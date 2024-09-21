package service

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUserInfo(userGuid string) (*dto.UserInfoDto, error) {
	return s.repository.GetUserDetail(userGuid)
}

func (s *UserService) GetAllUsers() ([]dto.UserInfoDto, error) {
	return s.repository.GetAllUsers()
}

func (s *UserService) ChangeUserDetail(userGuid uuid.UUID, userInfo dto.ChangeUserInfoDto) error {
	userDetail := models.UserInfo{
		Role:     userInfo.Role,
		UserGuid: userGuid,
	}
	return s.repository.ChangeOrCreateUserRole(&userDetail)
}
