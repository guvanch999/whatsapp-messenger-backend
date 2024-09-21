package repo

import (
	"errors"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"gorm.io/gorm"
)

type UserContactsRepository struct {
	db *gorm.DB
}

func NewUserContactRepository(db *gorm.DB) *UserContactsRepository {
	return &UserContactsRepository{
		db,
	}
}

func (r *UserContactsRepository) GetUserContacts(userId uuid.UUID) ([]models.UserContact, error) {
	var userContacts []models.UserContact
	if err := r.db.Model(&models.UserContact{}).Where(
		"user_id = ?",
		userId,
	).Select("*").Scan(&userContacts).Error; err != nil {
		return nil, err
	}
	return userContacts, nil
}

func (r *UserContactsRepository) GetUserContactsList(userId uuid.UUID, numbers []string) ([]models.UserContact, error) {
	var userContacts []models.UserContact
	if err := r.db.Model(&models.UserContact{}).Where(
		"user_id = ? and phone_number in (?)",
		userId,
		numbers,
	).Select("*").Scan(&userContacts).Error; err != nil {
		return nil, err
	}
	return userContacts, nil
}

func (r *UserContactsRepository) GetAllContacts() ([]models.UserContact, error) {
	var userContacts []models.UserContact
	if err := r.db.Model(&models.UserContact{}).Select("*").Scan(&userContacts).Error; err != nil {
		return nil, err
	}
	return userContacts, nil
}
func (r *UserContactsRepository) GetContactDetail(id uuid.UUID) (*models.UserContact, error) {
	var userContact models.UserContact
	if err := r.db.Model(&models.UserContact{}).Where("id = ?", id).Find(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &userContact, nil
}

func (r *UserContactsRepository) AddContact(contactModel models.UserContact) (*models.UserContact, error) {
	if err := r.db.Save(&contactModel).Error; err != nil {
		return nil, err
	}
	return &contactModel, nil
}

func (r *UserContactsRepository) AddMany(contacts []models.UserContact) ([]models.UserContact, error) {
	if err := r.db.Save(contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *UserContactsRepository) UpdateContact(contactModel models.UserContact) (*models.UserContact, error) {
	if err := r.db.Model(&models.UserContact{}).Where(
		"id = ?",
		contactModel.Id,
	).Updates(contactModel).Error; err != nil {
		return nil, err
	}
	return &contactModel, nil
}

func (r *UserContactsRepository) DeleteUserContact(id uuid.UUID) error {
	if err := r.db.Delete(&models.UserContact{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}
