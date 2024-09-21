package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type OrganizationHandler struct {
	service *service.OrganizationService
}

func NewOrganizationHandler(organizationService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService,
	}
}

// GetUserOrganizations godoc
//
//	@Summary	Get organizations
//	@Tags		Organizations
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ResponseOrganizationDto]   "Organizations"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/organizations [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *OrganizationHandler) GetUserOrganizations(c echo.Context) error {
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.GetUserOrganizations(user.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}
func (h *OrganizationHandler) GetAllOrganizations(c echo.Context) error {
	data, err := h.service.GetAllOrganizations()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// CreateOrganization godoc
//
//	@Summary	Create organization
//	@Tags		Organizations
//	@Accept		json
//	@Produce	json
//	@Param		Organization detail 	body		dto.OrganizationDto				true	"Organization detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseOrganizationDto]   "Organization"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/organizations/{guid} [put]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	var organizationDto dto.OrganizationDto
	if err := c.Bind(&organizationDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&organizationDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.CreateOrganization(user.ID, organizationDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// UpdateOrganization godoc
//
//	@Summary	Update organization
//	@Tags		Organizations
//	@Accept		json
//	@Produce	json
//	@Param		Organization detail 	body		dto.UpdateOrganization				true	"Organization detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseOrganizationDto]   "Update organization"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/organizations/{guid} [put]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *OrganizationHandler) UpdateOrganization(c echo.Context) error {
	var updateOrganization dto.UpdateOrganization
	if err := c.Bind(&updateOrganization); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&updateOrganization); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.UpdateOrganization(user, updateOrganization)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// DeleteOrganization godoc
//
//	@Summary	Delete organization
//	@Tags		Organizations
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.MessageWrapperDto   "Delete organization"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/organizations/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *OrganizationHandler) DeleteOrganization(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err = h.service.DeleteOrganization(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]string{
			"message": "Organization is removed",
		},
	)
}

// GetDetail godoc
//
//	@Summary	Organization detail
//	@Tags		Organizations
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseOrganizationDto ]   "Organization detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/organizations/{guid} [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *OrganizationHandler) GetDetail(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	detail, err := h.service.GetOrganizationDetail(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, detail)
}
