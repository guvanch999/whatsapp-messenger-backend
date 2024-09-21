package middleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	models2 "github.com/medium-messenger/messenger-backend/internal/modules/api-keys/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"github.com/nedpals/supabase-go"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func AuthMiddleware(supabaseClient *supabase.Client, db *gorm.DB) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenHeader := c.Request().Header.Get("Authorization")
			apiKey := c.Request().Header.Get("X-Api-Key")
			if tokenHeader == "" && apiKey == "" {
				return handleUnauthorized(c, "Invalid authorization")
			}
			if tokenHeader != "" {
				err := authWithToken(c, supabaseClient, db)
				if err != nil {
					return response.Error(c, err)
				}
			} else {
				err := authWithApiKey(c, db)
				if err != nil {
					return response.Error(c, err)
				}
			}
			return next(c)
		}
	}
}

func authWithToken(c echo.Context, supabaseClient *supabase.Client, db *gorm.DB) error {
	tokenHeader := c.Request().Header.Get("Authorization")
	var token string
	if token = strings.TrimPrefix(tokenHeader, "Bearer "); tokenHeader == "" || tokenHeader == token {
		return &exceptions.AuthFailed{
			Message: "Invalid authorization",
		}
	}
	user, err := supabaseClient.Auth.User(context.Background(), token)
	if err != nil {
		return &exceptions.AuthFailed{
			Message: "Invalid authorization",
		}
	}
	detail, err := getUserDetail(db, uuid.MustParse(user.ID))
	if err != nil {
		return err
	}
	userDetail := models.UserDetail{}
	userDetail.FromSupabaseUser(user, detail)
	c.Set("user", userDetail)
	c.Set("token", token)
	return nil
}

func authWithApiKey(c echo.Context, db *gorm.DB) error {
	apiKey := c.Request().Header.Get("X-Api-Key")
	user, err := checkApiKey(db, apiKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.AuthFailed{
				Message: "Invalid authorization",
			}
		}
		return err
	}
	detail, err := getUserDetail(db, user.ID)
	if err != nil {
		return err
	}
	userDetail := models.UserDetail{}
	userDetail.FromUser(user, detail)
	c.Set("user", userDetail)
	c.Set("api-key", apiKey)
	return nil
}

func handleUnauthorized(c echo.Context, message string) error {
	return c.JSON(
		http.StatusUnauthorized, &exceptions.AuthFailed{
			Message: message,
		},
	)
}

func getUserDetail(db *gorm.DB, userGuid uuid.UUID) (*models.UserInfo, error) {
	detail := new(models.UserInfo)
	if err := db.Table("user_info").Where("user_guid = ?", userGuid).First(detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.UserInfo{
				Id:       uuid.Nil,
				Role:     enums.User,
				UserGuid: userGuid,
			}, nil
		}
		return nil, err
	}
	return detail, nil
}

func checkApiKey(db *gorm.DB, apiKeyString string) (*models.User, error) {
	var list []models2.ApiKey
	if err := db.Model(&models2.ApiKey{}).Select("*").Scan(&list).Error; err != nil {
		return nil, err
	}
	apiKey := new(models2.ApiKey)
	for _, key := range list {
		err := util.CompareString(key.Hash, apiKeyString)
		if err == nil {
			apiKey = &key
			break
		}
	}
	if apiKey == nil {
		return nil, &exceptions.NotFoundError{}
	}

	var user models.User
	if err := db.Table("auth.users").Where("id = ?", apiKey.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}

	return &user, nil
}
