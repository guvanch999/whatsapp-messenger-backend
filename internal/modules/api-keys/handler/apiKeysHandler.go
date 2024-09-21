package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type ApiKeysHandler struct {
	service *service.ApiKeysService
}

func NewApiKeysHandler(keysService *service.ApiKeysService) *ApiKeysHandler {
	return &ApiKeysHandler{
		keysService,
	}
}

// GetUserApiKeys godoc
//
//	@Summary	Get user api keys
//	@Tags		Api keys
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ApiKeyResponse]   "User api keys"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/api-keys 		[get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ApiKeysHandler) GetUserApiKeys(c echo.Context) error {
	user := c.Get("user").(auth.UserDetail)
	list, err := h.service.GetUserApiKeys(user)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": list,
		},
	)
}

// AddApiKey godoc
//
//	@Summary	Create api key
//	@Tags		Api keys
//	@Accept		json
//	@Produce	json
//	@Param		Api key detail 	body		dto.ApiKeyDto				true	"Api key detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ApiKeyResponse]   "Api key detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/api-keys [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ApiKeysHandler) AddApiKey(c echo.Context) error {
	var apiKeyDto dto.ApiKeyDto
	if err := c.Bind(&apiKeyDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&apiKeyDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.AddApiKey(user, apiKeyDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// DeleteApiKey godoc
//
//	@Summary	Delete api key
//	@Tags		Api keys
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.MessageWrapperDto   "Delete api key"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/api-keys/{guid} 		[delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ApiKeysHandler) DeleteApiKey(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(auth.UserDetail)

	err = h.service.DeleteApiKey(user, guid)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(
		c, map[string]string{
			"message": "Api key is removed",
		},
	)
}

// GetDetail godoc
//
//	@Summary	Get user api key detail
//	@Tags		Api keys
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[dto.ApiKeyResponse ]   "Api key detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/api-keys/{guid} 		[get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ApiKeysHandler) GetDetail(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(auth.UserDetail)

	data, err := h.service.GetDetail(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// GetApiKeyValue godoc
//
//	@Summary	Api key value
//	@Tags		Api keys
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[dto.ApiKeyDetailResponse]   "Api key detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/api-keys/value/{guid} 		[get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ApiKeysHandler) GetApiKeyValue(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.GetValueOfApiKey(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}
