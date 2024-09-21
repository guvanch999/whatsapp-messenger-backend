package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/service"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return response.Error(c, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) ChangeUserInfo(c echo.Context) error {
	userGuid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	userInfo := dto.ChangeUserInfoDto{}
	if err = c.Bind(&userInfo); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err = c.Validate(userInfo); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	err = h.service.ChangeUserDetail(userGuid, userInfo)

	if err != nil {
		return response.Error(c, err)
	}
	return c.String(http.StatusOK, "Details is changed")
}
