package service

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/model"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type ContactListService struct {
	repository *repository.ContactListRepository
}

func NewContactListService(listRepository *repository.ContactListRepository) *ContactListService {
	return &ContactListService{
		repository: listRepository,
	}
}

func (s *ContactListService) GetAllLists() ([]dto.ResponseList, error) {
	list, err := s.repository.GetAllContactLists()
	if err != nil {
		return nil, err
	}
	return util.Map(
		list, func(l model.ContactList) dto.ResponseList {
			return *l.ToResponseDto()
		},
	), err
}

func (s *ContactListService) GetUserLists(userId uuid.UUID) ([]dto.ResponseList, error) {
	list, err := s.repository.GetUserContactLists(userId)
	if err != nil {
		return nil, err
	}
	return util.Map(
		list, func(l model.ContactList) dto.ResponseList {
			return *l.ToResponseDto()
		},
	), err
}

func (s *ContactListService) AddContactList(user auth.UserDetail, listDto dto.ContactListDto) (
	*dto.ResponseList,
	error,
) {
	listModel := model.ContactList{
		UserID: user.ID,
		Name:   listDto.Name,
		Contacts: util.Map(
			listDto.ContactList, func(contactId uuid.UUID) models.UserContact {
				return models.UserContact{
					Id: contactId,
				}
			},
		),
	}
	contactList, err := s.repository.AddContactList(listModel)
	if err != nil {
		return nil, err
	}
	return contactList.ToResponseDto(), nil
}

func (s *ContactListService) UpdateContactListName(
	user auth.UserDetail,
	updateDto dto.UpdateContactListNameDto,
) (*dto.ResponseList, error) {
	listModel, err := s.checkAccess(user, updateDto.Id)
	if err != nil {
		return nil, err
	}
	listModel.Name = updateDto.Name
	result, err := s.repository.UpdateContactList(*listModel)
	if err != nil {
		return nil, err
	}
	return result.ToResponseDto(), nil
}
func (s *ContactListService) GetDetail(user auth.UserDetail, id uuid.UUID) (*dto.ResponseList, error) {
	listModel, err := s.checkAccess(user, id)
	if err != nil {
		return nil, err
	}
	return listModel.ToResponseDto(), nil
}

func (s *ContactListService) DeleteList(user auth.UserDetail, id uuid.UUID) error {
	_, err := s.checkAccess(user, id)
	if err != nil {
		return err
	}
	return s.repository.DeleteContactList(id)
}

func (s *ContactListService) RemoveContactFromList(user auth.UserDetail, updateDto dto.ContactListUpdateDto) error {
	_, err := s.checkAccess(user, updateDto.Id)
	if err != nil {
		return err
	}
	return s.repository.RemoveContactFromList(updateDto.Id, updateDto.ContactList)
}

func (s *ContactListService) AddContactToList(user auth.UserDetail, updateDto dto.ContactListUpdateDto) error {
	_, err := s.checkAccess(user, updateDto.Id)
	if err != nil {
		return err
	}
	return s.repository.AddContactToList(updateDto.Id, updateDto.ContactList)
}

func (s *ContactListService) checkAccess(user auth.UserDetail, id uuid.UUID) (*model.ContactList, error) {
	contactModel, err := s.repository.GetDetail(id)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && contactModel.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return contactModel, nil
}
