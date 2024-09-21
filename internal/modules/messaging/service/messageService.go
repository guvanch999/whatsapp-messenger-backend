package service

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/messaging/dto"
	template "github.com/medium-messenger/messenger-backend/internal/modules/templates/service"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/model"
	providers "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"github.com/nyaruka/phonenumbers"
	"github.com/twilio/twilio-go"
	openapi2 "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/gorm"
)

type MessageService struct {
	db                    *gorm.DB
	secretManagerClient   *secretmanager.Client
	templateService       *template.TemplateService
	contactListRepository *repository.ContactListRepository
}

func NewMessageService(
	db *gorm.DB,
	client *secretmanager.Client,
	service *template.TemplateService,
	listRepository *repository.ContactListRepository,
) *MessageService {
	return &MessageService{
		db,
		client,
		service,
		listRepository,
	}
}

func (s *MessageService) SendMessages(
	user auth.UserDetail,
	sendMessageDto dto.SendMessageDto,
) ([]dto.SendMessageResponse, error) {
	provider, cred, err := providers.GetProviderWithCred[model.TwilioCred](
		s.db,
		s.secretManagerClient,
		user,
		sendMessageDto.ProviderId,
	)
	if err != nil {
		return nil, err
	}

	contacts, err := s.contactListRepository.GetContactsWithIds(
		util.Map(
			sendMessageDto.Recipients, func(d dto.Recipient) uuid.UUID {
				return d.RecipientId
			},
		),
	)
	if err != nil {
		return nil, err
	}

	twilioClient := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: cred.TwilioAccountSid,
			Password: cred.TwilioAuthToken,
		},
	)

	teml, err := s.templateService.GetDetail(user, sendMessageDto.TemplateId)
	if err != nil {
		return nil, err
	}

	lenContacts := len(sendMessageDto.Recipients)
	jobs := make(chan dto.MessageDetailDto, lenContacts)
	results := make(chan dto.SendMessageResponse, lenContacts)

	for w := 0; w < 10; w++ {
		go sendMessageWorker(twilioClient, jobs, results)
	}

	for j := 0; j < lenContacts; j++ {
		phoneNumber := ""
		for _, cont := range contacts {
			if cont.Id == sendMessageDto.Recipients[j].RecipientId {
				phoneNumber = cont.PhoneNumber
			}
		}
		if len(phoneNumber) == 0 {
			continue
		}

		jobs <- dto.MessageDetailDto{
			PhoneNumber:       phoneNumber,
			FromPhoneNumber:   provider.FromPhoneNumber,
			ServiceId:         cred.TwilioMessagingServiceSid,
			TemplateId:        teml.ExternalId,
			TemplateVariables: sendMessageDto.Recipients[j].Variables,
		}
	}
	close(jobs)
	processedResult := make([]dto.SendMessageResponse, lenContacts)
	for a := 0; a < lenContacts; a++ {
		processedResult[a] = <-results
	}
	return processedResult, nil
}

func sendMessageWorker(
	twilioClient *twilio.RestClient,
	messages <-chan dto.MessageDetailDto,
	results chan<- dto.SendMessageResponse,
) {
	for message := range messages {
		params := &openapi2.CreateMessageParams{}
		parsedNumber, err := phonenumbers.Parse(message.PhoneNumber, "")
		if err != nil {
			results <- dto.SendMessageResponse{
				PhoneNumber:  message.PhoneNumber,
				Status:       enums.Fail,
				ErrorMessage: err.Error(),
			}
		}

		formattedNumber := phonenumbers.Format(parsedNumber, phonenumbers.E164)

		params.SetTo("whatsapp:" + formattedNumber)
		params.SetFrom("whatsapp:" + message.FromPhoneNumber)

		params.SetMessagingServiceSid(message.ServiceId)
		params.SetContentSid(message.TemplateId)

		contentByte, err := json.Marshal(message.TemplateVariables)
		if err != nil {
			results <- dto.SendMessageResponse{
				PhoneNumber:  message.PhoneNumber,
				Status:       enums.Fail,
				ErrorMessage: err.Error(),
			}
		}
		params.SetContentVariables(string(contentByte))

		_, err = twilioClient.Api.CreateMessage(params)
		if err != nil {
			results <- dto.SendMessageResponse{
				PhoneNumber:  message.PhoneNumber,
				Status:       enums.Fail,
				ErrorMessage: err.Error(),
			}
		}
		results <- dto.SendMessageResponse{
			PhoneNumber: message.PhoneNumber,
			Status:      enums.Success,
		}

	}
}

func (s *MessageService) SendMessageToGroup(
	user auth.UserDetail,
	sendMessageDto dto.SendMessageToListDto,
) ([]dto.SendMessageResponse, error) {
	provider, cred, err := providers.GetProviderWithCred[model.TwilioCred](
		s.db,
		s.secretManagerClient,
		user,
		sendMessageDto.ProviderId,
	)
	if err != nil {
		return nil, err
	}
	var contacts []models.UserContact

	detail, err := s.contactListRepository.GetDetail(*sendMessageDto.ContactListId)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && detail.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	contacts = append(contacts, detail.Contacts...)

	twilioClient := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: cred.TwilioAccountSid,
			Password: cred.TwilioAuthToken,
		},
	)

	teml, err := s.templateService.GetDetail(user, sendMessageDto.TemplateId)
	if err != nil {
		return nil, err
	}

	lenContacts := len(contacts)
	jobs := make(chan dto.MessageDetailDto, lenContacts)
	results := make(chan dto.SendMessageResponse, lenContacts)

	for w := 0; w < 10; w++ {
		go sendMessageWorker(twilioClient, jobs, results)
	}

	for j := 0; j < lenContacts; j++ {
		jobs <- dto.MessageDetailDto{
			PhoneNumber:       contacts[j].PhoneNumber,
			FromPhoneNumber:   provider.FromPhoneNumber,
			ServiceId:         cred.TwilioMessagingServiceSid,
			TemplateId:        teml.ExternalId,
			TemplateVariables: sendMessageDto.TemplateVariables,
		}
	}
	close(jobs)
	processedResult := make([]dto.SendMessageResponse, lenContacts)
	for a := 0; a < lenContacts; a++ {
		processedResult[a] = <-results
	}
	return processedResult, nil
}
