package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type UserProviderHandler struct {
	service *service.UserProviderService
}

func NewUserProviderHandler(providerService *service.UserProviderService) *UserProviderHandler {
	return &UserProviderHandler{
		providerService,
	}
}

// GetUserProviders godoc
//
//	@Summary	Get user providers
//	@Tags		User providers
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ResponseProviderDto]   "User providers"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-providers [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserProviderHandler) GetUserProviders(c echo.Context) error {
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.GetUserProviders(user.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}
func (h *UserProviderHandler) GetAllProviders(c echo.Context) error {
	data, err := h.service.GetAllProviders()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// CreateProvider godoc
//
//	@Summary	Add provider
//	@Tags		User providers
//	@Accept		json
//	@Produce	json
//	@Param		Add provider 	body		dto.UserProviderDto				true	"Provider detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseProviderDto]   "Provider detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-providers [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserProviderHandler) CreateProvider(c echo.Context) error {
	var providerDto dto.UserProviderDto
	if err := c.Bind(&providerDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	cred, err := h.service.ValidateProviderDto(c, &providerDto)
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.CreateProvider(user.ID, providerDto, cred)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// DeleteProvider godoc
//
//	@Summary	Delete provider
//	@Tags		User providers
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.MessageWrapperDto   "Delete provider"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-providers/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserProviderHandler) DeleteProvider(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err = h.service.DeleteProvider(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]string{
			"message": "Provider is removed",
		},
	)
}

// GetDetail godoc
//
//	@Summary	Get provider detail
//	@Tags		User providers
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseProviderDto]   "User provider detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-providers/{guid} [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserProviderHandler) GetDetail(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	detail, err := h.service.GetProviderDetail(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, detail)
}
