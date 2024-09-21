package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserDetail(userGuid string) (*dto.UserInfoDto, error) {
	var user models.User
	if err := r.db.Table("auth.users").Select("*").Where("id = ?", userGuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}

	detail := new(models.UserInfo)
	if err := r.db.Model(models.UserInfo{}).Where("user_guid = ?", userGuid).First(detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			detail = &models.UserInfo{
				Id:       uuid.Nil,
				Role:     enums.User,
				UserGuid: user.ID,
			}
		}
		return nil, err
	}

	return &dto.UserInfoDto{
		User: user,
		Role: detail.Role,
	}, nil
}

func (r *UserRepository) GetAllUsers() ([]dto.UserInfoDto, error) {
	var users []models.User
	if err := r.db.Table("auth.users").Select("*").Scan(&users).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var userInfo []models.UserInfo

	if err := r.db.Model(&models.UserInfo{}).Select("*").Scan(&userInfo).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return util.Map(
		users, func(user models.User) dto.UserInfoDto {
			userDto := dto.UserInfoDto{
				Role: enums.User,
				User: user,
			}
			for _, info := range userInfo {
				if info.UserGuid == user.ID {
					userDto.Role = info.Role
					break
				}
			}
			return userDto
		},
	), nil
}

func (r *UserRepository) ChangeOrCreateUserRole(userInfo *models.UserInfo) error {
	detail := new(models.UserInfo)
	if err := r.db.Model(&models.UserInfo{}).Select("*").Where(
		"user_guid = ?",
		userInfo.UserGuid,
	).First(detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.CheckUserExists(userInfo.UserGuid)
			if err != nil {
				return err
			}
			err = r.CreateDetail(userInfo)
			if err != nil {
				return err
			}
		}
		return err
	}
	detail.Role = userInfo.Role
	return r.UpdateDetail(detail)
}

func (r *UserRepository) UpdateDetail(userInfo *models.UserInfo) error {
	if err := r.db.Model(&models.UserInfo{}).Where("id = ?", userInfo.Id).Updates(userInfo).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateDetail(userInfo *models.UserInfo) error {
	if err := r.db.Model(&models.UserInfo{}).Create(userInfo).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CheckUserExists(userGuid uuid.UUID) error {
	var usersShort []models.User
	if err := r.db.Table("auth.users").Select("*").Where("id = ?", userGuid).First(&usersShort).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}
