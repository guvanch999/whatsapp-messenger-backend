package service

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/repo"
	models2 "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"github.com/nyaruka/phonenumbers"
	"log"
)

type UserContactsService struct {
	repository *repo.UserContactsRepository
}

func NewUserContactsService(userRepository *repo.UserContactsRepository) *UserContactsService {
	return &UserContactsService{
		repository: userRepository,
	}
}

func (s *UserContactsService) GetUserContacts(userId uuid.UUID) ([]dto.ContactResponse, error) {
	contacts, err := s.repository.GetUserContacts(userId)
	if err != nil {
		return nil, err
	}
	return util.Map(
		contacts, func(contact models.UserContact) dto.ContactResponse {
			return *contact.ToResponseDto()
		},
	), nil
}

func (s *UserContactsService) GetAllUserContacts() ([]dto.ContactResponse, error) {
	contacts, err := s.repository.GetAllContacts()
	if err != nil {
		return nil, err
	}
	return util.Map(
		contacts, func(contact models.UserContact) dto.ContactResponse {
			return *contact.ToResponseDto()
		},
	), nil
}

func (s *UserContactsService) AddUserContact(
	user models2.UserDetail,
	userContactDto dto.UserContactDto,
) (*dto.CreatedOrExistResponse, error) {
	responseDto := new(dto.CreatedOrExistResponse)
	existList, err := s.getExistingNumbers(user.ID, userContactDto.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if len(existList) > 0 {
		responseDto.Exist = append(responseDto.Exist, *existList[0].ToResponseDto())
		return responseDto, nil
	}

	contactModel := models.UserContact{}
	contactModel.FromDto(userContactDto)
	contactModel.UserID = user.ID
	contact, err := s.repository.AddContact(contactModel)
	if err != nil {
		return nil, err
	}
	responseDto.Created = append(responseDto.Created, *contact.ToResponseDto())
	return responseDto, nil
}

func (s *UserContactsService) AddContactsList(
	user models2.UserDetail,
	listDto dto.UserContactsListDto,
) (*dto.CreatedOrExistResponse, error) {
	responseDto := new(dto.CreatedOrExistResponse)
	numbers := util.Map(
		listDto.UserContactsList, func(uDto dto.UserContactDto) string {
			return uDto.PhoneNumber
		},
	)
	existList, err := s.getExistingNumbers(user.ID, numbers...)
	if err != nil {
		return nil, err
	}
	var contacts []models.UserContact
	for _, uDto := range listDto.UserContactsList {
		exist := false
		for _, cModel := range existList {
			if cModel.PhoneNumber == uDto.PhoneNumber {
				responseDto.Exist = append(responseDto.Exist, *cModel.ToResponseDto())
				exist = true
				break
			}
		}
		if !exist {
			contactModel := models.UserContact{}
			contactModel.FromDto(uDto)
			contactModel.UserID = user.ID
			contacts = append(contacts, contactModel)
		}
	}

	if len(contacts) == 0 {
		return responseDto, nil
	}
	log.Println("must create")
	data, err := s.repository.AddMany(contacts)
	if err != nil {
		return nil, err
	}
	for _, crtd := range data {
		responseDto.Created = append(responseDto.Created, *crtd.ToResponseDto())
	}

	return responseDto, nil
}

func (s *UserContactsService) UpdateContact(
	user models2.UserDetail,
	contactDto dto.UpdateContactDto,
) (*dto.ContactResponse, error) {
	contactModel, err := s.checkAccess(user, contactDto.Id)
	if err != nil {
		return nil, err
	}
	contactModel.FromUpdateDto(contactDto)
	contact, err := s.repository.UpdateContact(*contactModel)
	if err != nil {
		return nil, err
	}
	return contact.ToResponseDto(), nil
}

func (s *UserContactsService) DeleteContact(user models2.UserDetail, contactId uuid.UUID) error {
	_, err := s.checkAccess(user, contactId)
	if err != nil {
		return err
	}
	return s.repository.DeleteUserContact(contactId)
}
func (s *UserContactsService) checkAccess(user models2.UserDetail, contactId uuid.UUID) (*models.UserContact, error) {
	contact, err := s.repository.GetContactDetail(contactId)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && contact.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return contact, nil
}

func (s *UserContactsService) GetContactDetail(user models2.UserDetail, contactId uuid.UUID) (
	*dto.ContactResponse,
	error,
) {
	contact, err := s.checkAccess(user, contactId)
	if err != nil {
		return nil, err
	}
	return contact.ToResponseDto(), nil
}

func (s *UserContactsService) ValidateNumbersSync(numbers []string) []dto.NumberValidateResponse {
	processedResult := make([]dto.NumberValidateResponse, len(numbers))
	for index, number := range numbers {
		_, err := phonenumbers.Parse(number, "")
		if err != nil {
			processedResult[index] = dto.NumberValidateResponse{
				Number:       number,
				IsValid:      false,
				ErrorMessage: err.Error(),
			}
		} else {
			processedResult[index] = dto.NumberValidateResponse{
				Number:  number,
				IsValid: true,
			}
		}
	}
	return processedResult
}

func (s *UserContactsService) getExistingNumbers(userId uuid.UUID, numbers ...string) ([]models.UserContact, error) {
	if len(numbers) == 0 {
		return []models.UserContact{}, nil
	}
	numberExistingList, err := s.repository.GetUserContactsList(userId, numbers)
	if err != nil {
		return nil, err
	}
	return numberExistingList, nil
}

func (s *UserContactsService) ValidateNumbers(numbers []string) []dto.NumberValidateResponse {
	lenNumbers := len(numbers)
	jobs := make(chan string, lenNumbers)
	results := make(chan dto.NumberValidateResponse, lenNumbers)

	for w := 0; w < 10; w++ {
		go validateWorker(jobs, results)
	}

	for j := 0; j < lenNumbers; j++ {
		jobs <- numbers[j]
	}
	close(jobs)
	processedResult := make([]dto.NumberValidateResponse, lenNumbers)
	for a := 0; a < lenNumbers; a++ {
		processedResult[a] = <-results
	}
	return processedResult
}

func validateWorker(numbers <-chan string, results chan<- dto.NumberValidateResponse) {
	for number := range numbers {
		_, err := phonenumbers.Parse(number, "")
		if err != nil {
			results <- dto.NumberValidateResponse{
				Number:       number,
				IsValid:      false,
				ErrorMessage: err.Error(),
			}
		} else {
			results <- dto.NumberValidateResponse{
				Number:  number,
				IsValid: true,
			}
		}

	}
}
