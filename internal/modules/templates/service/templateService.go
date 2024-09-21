package service

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/repository"
	dto2 "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/dto"
	providers "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/content/v1"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

// TemplateService Todo make cron for status changes
type TemplateService struct {
	db                  *gorm.DB
	secretManagerClient *secretmanager.Client
	repository          *repository.TemplateRepository
}

func NewTemplateService(
	db *gorm.DB,
	client *secretmanager.Client,
	templateRepository *repository.TemplateRepository,
) *TemplateService {
	return &TemplateService{
		db:                  db,
		secretManagerClient: client,
		repository:          templateRepository,
	}
}

func (s *TemplateService) GetAllTemplates() ([]dto.ResponseTemplateDto, error) {
	list, err := s.repository.GetAllTemplates()
	if err != nil {
		return nil, err
	}
	return util.Map(
		list, func(l models.Template) dto.ResponseTemplateDto {
			return *l.ToResponseDto()
		},
	), err
}

func (s *TemplateService) GetUserTemplates(userId uuid.UUID) ([]dto.ResponseTemplateDto, error) {
	list, err := s.repository.GetUserTemplates(userId)
	if err != nil {
		return nil, err
	}
	return util.Map(
		list, func(l models.Template) dto.ResponseTemplateDto {
			return *l.ToResponseDto()
		},
	), err
}

func (s *TemplateService) AddTemplate(user auth.UserDetail, templateDto dto.CreateTemplateDto) (
	*dto.ResponseTemplateDto,
	error,
) {
	templateModel := models.Template{}
	templateModel.FromDto(&templateDto)
	templateModel.UserID = user.ID
	templateSid, err := s.CreateTemplateInTwilio(templateModel, user, templateDto.ProviderId)
	if err != nil {
		return nil, err
	}
	templateModel.ExternalId = templateSid
	template, err := s.repository.AddTemplate(templateModel)
	if err != nil {
		return nil, err
	}
	return template.ToResponseDto(), nil
}

func (s *TemplateService) UpdateTemplate(
	user auth.UserDetail,
	updateDto dto.UpdateTemplateDto,
) (*dto.ResponseTemplateDto, error) {
	template, err := s.checkAccess(user, updateDto.Id)
	if err != nil {
		return nil, err
	}
	template.FromUpdateDto(&updateDto)
	result, err := s.repository.UpdateTemplate(*template)
	if err != nil {
		return nil, err
	}
	return result.ToResponseDto(), nil
}
func (s *TemplateService) GetDetail(user auth.UserDetail, id uuid.UUID) (*dto.ResponseTemplateDto, error) {
	template, err := s.checkAccess(user, id)
	if err != nil {
		return nil, err
	}
	return template.ToResponseDto(), nil
}

func (s *TemplateService) DeleteTemplate(user auth.UserDetail, id uuid.UUID) error {
	template, err := s.checkAccess(user, id)
	if err != nil {
		return err
	}
	_, cred, err := providers.GetProviderWithCred[dto2.TwilioCredDto](
		s.db,
		s.secretManagerClient,
		user,
		template.ProviderId,
	)
	if err != nil {
		return err
	}
	twilioClient := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: cred.TwilioAccountSid,
			Password: cred.TwilioAuthToken,
		},
	)
	err = twilioClient.ContentV1.DeleteContent(template.ExternalId)
	if err != nil {
		return err
	}
	return s.repository.DeleteTemplate(id)
}

func (s *TemplateService) ApproveTemplate(
	user auth.UserDetail,
	approvalDto dto.TwilioTemplateApprovalDto,
) (*dto.ResponseTemplateDto, error) {
	template, err := s.checkAccess(user, approvalDto.Id)
	if err != nil {
		return nil, err
	}
	_, cred, err := providers.GetProviderWithCred[dto2.TwilioCredDto](
		s.db,
		s.secretManagerClient,
		user,
		template.ProviderId,
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
	data, err := twilioClient.ContentV1.CreateApprovalCreate(
		template.ExternalId, &openapi.CreateApprovalCreateParams{
			ContentApprovalRequest: &openapi.ContentApprovalRequest{
				Name:     approvalDto.Name,
				Category: approvalDto.Category,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	status, err := enums.StatusFromString(*data.Status)
	if err != nil {
		return nil, err
	}
	template.Status = status
	template.NextCheck = time.Now().Add(time.Minute * 5)

	tmp, err := s.repository.UpdateTemplate(*template)
	if err != nil {
		return nil, err
	}
	return tmp.ToResponseDto(), nil
}

func (s *TemplateService) checkAccess(user auth.UserDetail, id uuid.UUID) (*models.Template, error) {
	template, err := s.repository.GetDetail(id)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && template.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return template, nil
}

func (s *TemplateService) CreateTemplateInTwilio(
	templateModel models.Template,
	user auth.UserDetail,
	providerId uuid.UUID,
) (string, error) {
	content, err := templateModel.GetContent()
	if err != nil {
		return "", err
	}
	_, cred, err := providers.GetProviderWithCred[dto2.TwilioCredDto](s.db, s.secretManagerClient, user, providerId)
	if err != nil {
		return "", err
	}
	twilioClient := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: cred.TwilioAccountSid,
			Password: cred.TwilioAuthToken,
		},
	)

	data, err := twilioClient.ContentV1.CreateContent(
		&openapi.CreateContentParams{
			ContentCreateRequest: content,
		},
	)
	if err != nil {
		return "", err
	}
	return *data.Sid, nil
}

func (s *TemplateService) SyncTemplateStatuses() error {
	templates, err := s.repository.GetTemplatesBeforeTime(time.Now())
	var wg sync.WaitGroup
	wg.Add(len(templates))
	for _, template := range templates {
		go s.syncWorker(&wg, template)
	}
	wg.Wait()
	return err
}

func (s *TemplateService) syncWorker(wg *sync.WaitGroup, template models.Template) {
	defer wg.Done()
	_, cred, err := providers.GetProviderWithCredWithoutCheck[dto2.TwilioCredDto](
		s.db,
		s.secretManagerClient,
		template.ProviderId,
	)
	if err != nil {
		log.Printf("Cannot get provider detail of template: %s\n", template.Id)
		return
	}
	twilioClient := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: cred.TwilioAccountSid,
			Password: cred.TwilioAuthToken,
		},
	)
	providerTemplate, err := twilioClient.ContentV1.FetchApprovalFetch(template.ExternalId)
	if err != nil {
		log.Printf("Cannot get provider template detail of template: %s\n", template.Id)
		return
	}
	byteArr, err := json.Marshal(providerTemplate.Whatsapp)
	if err != nil {
		log.Printf("Cannot marshal to json: %s\n", template.Id)
		return
	}
	var approvalRequest dto.TwilioApprovalRequestDto
	err = json.Unmarshal(byteArr, &approvalRequest)
	if err != nil {
		log.Printf("Cannot unmarshal from json: %s\n", template.Id)
		return
	}
	updates := make(map[string]any)
	if approvalRequest.Status != template.Status {
		updates["status"] = approvalRequest.Status
	}
	if approvalRequest.Status == enums.InReview || approvalRequest.Status == "submitted" {
		updates["next_check"] = time.Now().Add(time.Minute * 5)
	} else {
		updates["next_check"] = nil
	}
	err = s.repository.UpdateTemplateWithUpdates(template.Id, updates)
	if err != nil {
		log.Printf("cannot upate template with id:%s \n", err.Error())
		return
	}
}
