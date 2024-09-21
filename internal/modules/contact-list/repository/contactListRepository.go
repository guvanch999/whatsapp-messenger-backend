package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/model"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"gorm.io/gorm"
)

type ContactListRepository struct {
	db *gorm.DB
}

type ContactListContacts struct {
	ContactListId uuid.UUID `json:"contact_list_id"`
	UserContactId uuid.UUID `json:"user_contact_id"`
}

func NewContactListRepository(db *gorm.DB) *ContactListRepository {
	return &ContactListRepository{
		db,
	}
}

func (r *ContactListRepository) GetUserContactLists(userId uuid.UUID) ([]model.ContactList, error) {
	var list []model.ContactList
	if err := r.db.Model(&model.ContactList{}).Select("*").Preload("Contacts").Where(
		"user_id = ?",
		userId,
	).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ContactListRepository) GetContactsWithIds(contactIds []uuid.UUID) ([]models.UserContact, error) {
	var list []models.UserContact
	if err := r.db.Model(&models.UserContact{}).Select("*").Where(
		"id in (?)",
		contactIds,
	).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ContactListRepository) GetAllContactLists() ([]model.ContactList, error) {
	var list []model.ContactList
	if err := r.db.Model(&model.ContactList{}).Select("*").Preload("Contacts").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ContactListRepository) GetDetail(listId uuid.UUID) (*model.ContactList, error) {
	var listDetail model.ContactList
	if err := r.db.Model(&model.ContactList{}).Select("*").Preload("Contacts").Where(
		"id = ?",
		listId,
	).First(&listDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &listDetail, nil
}

func (r *ContactListRepository) AddContactList(contactList model.ContactList) (*model.ContactList, error) {
	if err := r.db.Model(&model.ContactList{}).Create(&contactList).Error; err != nil {
		return nil, err
	}
	return &contactList, nil
}

func (r *ContactListRepository) UpdateContactList(contactList model.ContactList) (*model.ContactList, error) {
	if err := r.db.Model(&model.ContactList{}).Where("id = ?", contactList.Id).Updates(contactList).Error; err != nil {
		return nil, err
	}
	return &contactList, nil
}
func (r *ContactListRepository) DeleteContactList(listId uuid.UUID) error {
	if err := r.db.Delete(&model.ContactList{}, listId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}

func (r *ContactListRepository) RemoveContactFromList(listId uuid.UUID, contactIds []uuid.UUID) error {
	if err := r.db.Table("contact_list_contacts").Where(
		" contact_list_id = ?  and user_contact_id in (?)",
		listId,
		contactIds,
	).Delete(&ContactListContacts{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}

func (r *ContactListRepository) AddContactToList(listId uuid.UUID, contactIds []uuid.UUID) error {
	contacts := util.Map(
		contactIds, func(contactId uuid.UUID) ContactListContacts {
			return ContactListContacts{
				ContactListId: listId,
				UserContactId: contactId,
			}
		},
	)
	if err := r.db.Table("contact_list_contacts").Create(contacts).Error; err != nil {
		return err
	}
	return nil
}
