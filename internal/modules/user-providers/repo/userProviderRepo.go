package repo

import (
	"errors"
	"github.com/google/uuid"
	. "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/model"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"gorm.io/gorm"
)

type UserProviderRepository struct {
	db *gorm.DB
}

func NewUserProviderRepository(db *gorm.DB) *UserProviderRepository {
	return &UserProviderRepository{
		db,
	}
}

func (r *UserProviderRepository) GetUserProviders(userId uuid.UUID) ([]UserProvider, error) {
	var list []UserProvider
	if err := r.db.Model(&UserProvider{}).Select("*").Where(
		"user_id = ?",
		userId,
	).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *UserProviderRepository) GetAllProviders() ([]UserProvider, error) {
	var list []UserProvider
	if err := r.db.Model(&UserProvider{}).Select("*").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *UserProviderRepository) GetDetail(providerId uuid.UUID) (*UserProvider, error) {
	var providerDetail UserProvider
	if err := r.db.Model(&UserProvider{}).Select("*").Where(
		"id = ?",
		providerId,
	).First(&providerDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &providerDetail, nil
}

func (r *UserProviderRepository) AddProvider(provider UserProvider) (*UserProvider, error) {
	if err := r.db.Create(&provider).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *UserProviderRepository) UpdateProvider(provider UserProvider) (*UserProvider, error) {
	if err := r.db.Model(&UserProvider{}).Where("id = ?", provider.Id).Updates(provider).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}
func (r *UserProviderRepository) DeleteProvider(providerId uuid.UUID) error {
	if err := r.db.Delete(&UserProvider{}, providerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}
