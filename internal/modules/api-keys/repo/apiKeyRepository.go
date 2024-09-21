package repo

import (
	"errors"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"gorm.io/gorm"
)

type ApiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{
		db,
	}
}

func (r *ApiKeyRepository) AddApiKey(apiKey models.ApiKey) (*models.ApiKey, error) {
	if err := r.db.Save(&apiKey).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *ApiKeyRepository) GetUserApiKeys(userId uuid.UUID) ([]models.ApiKey, error) {
	var list []models.ApiKey
	if err := r.db.Model(&models.ApiKey{}).Select("*").Where("user_id = ?", userId).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ApiKeyRepository) DeleteApiKey(id uuid.UUID) error {
	if err := r.db.Delete(&models.ApiKey{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ApiKeyRepository) GetDetail(id uuid.UUID) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	if err := r.db.Model(&models.ApiKey{}).Where("id = ?", id).First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &apiKey, nil
}
