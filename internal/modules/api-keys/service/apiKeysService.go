package service

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/config"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/repo"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type ApiKeysService struct {
	cnf        *config.Schema
	repository *repo.ApiKeyRepository
}

func NewApiKeysService(schema *config.Schema, repository *repo.ApiKeyRepository) *ApiKeysService {
	return &ApiKeysService{
		cnf:        schema,
		repository: repository,
	}
}

func (s *ApiKeysService) AddApiKey(user auth.UserDetail, keyDto dto.ApiKeyDto) (*dto.ApiKeyResponse, error) {
	apiKey := models.ApiKey{
		UserID: user.ID,
		Name:   keyDto.Name,
	}
	hash, err := util.HashString(keyDto.ApiKey)
	if err != nil {
		return nil, err
	}
	encrypted, err := util.EncryptString(keyDto.ApiKey, s.cnf.SecretKeyForHash)
	if err != nil {
		return nil, err
	}
	apiKey.Hash = hash
	apiKey.Encoded = encrypted
	key, err := s.repository.AddApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	return key.ToResponseDto(), nil
}

func (s *ApiKeysService) GetUserApiKeys(user auth.UserDetail) ([]dto.ApiKeyResponse, error) {
	list, err := s.repository.GetUserApiKeys(user.ID)
	if err != nil {
		return nil, err
	}
	return util.Map(
		list, func(key models.ApiKey) dto.ApiKeyResponse {
			return *key.ToResponseDto()
		},
	), nil
}

func (s *ApiKeysService) DeleteApiKey(user auth.UserDetail, id uuid.UUID) error {
	_, err := s.checkAccess(user, id)
	if err != nil {
		return err
	}
	return s.repository.DeleteApiKey(id)
}

func (s *ApiKeysService) GetDetail(user auth.UserDetail, id uuid.UUID) (*dto.ApiKeyResponse, error) {
	apiKey, err := s.checkAccess(user, id)
	if err != nil {
		return nil, err
	}
	return apiKey.ToResponseDto(), nil
}

func (s *ApiKeysService) GetValueOfApiKey(user auth.UserDetail, id uuid.UUID) (*dto.ApiKeyDetailResponse, error) {
	apiKey, err := s.checkAccess(user, id)
	if err != nil {
		return nil, err
	}
	value, err := util.DecryptString(apiKey.Encoded, s.cnf.SecretKeyForHash)
	if err != nil {
		return nil, err
	}
	return &dto.ApiKeyDetailResponse{
		ApiKey: value,
	}, nil
}

func (s *ApiKeysService) checkAccess(user auth.UserDetail, id uuid.UUID) (*models.ApiKey, error) {
	apiKey, err := s.repository.GetDetail(id)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && apiKey.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return apiKey, nil
}
