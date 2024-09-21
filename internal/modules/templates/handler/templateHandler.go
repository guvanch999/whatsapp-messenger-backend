package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"log"
)

type TemplateHandler struct {
	service *service.TemplateService
}

func NewTemplateHandler(templateService *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService,
	}
}

// GetMyTemplates godoc
//
//	@Summary	Get user templates
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ResponseTemplateDto]   "User templates"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates 		[get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) GetMyTemplates(c echo.Context) error {
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.GetUserTemplates(user.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

func (h *TemplateHandler) GetAllTemplates(c echo.Context) error {
	data, err := h.service.GetAllTemplates()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// CreateTemplate godoc
//
//	@Summary	Add template
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Param		Add template 	body		dto.CreateTemplateDto				true	"Template detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseProviderDto]   "Provider detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) CreateTemplate(c echo.Context) error {
	var templateDto dto.CreateTemplateDto
	if err := c.Bind(&templateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&templateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}

	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.AddTemplate(user, templateDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// UpdateTemplate godoc
//
//	@Summary	Update template
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Param		Update template 	body		dto.UpdateTemplateDto				true	"Template detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseProviderDto]   "Provider detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates/{guid} [put]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) UpdateTemplate(c echo.Context) error {
	var updateTemplateDto dto.UpdateTemplateDto
	if err := c.Bind(updateTemplateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(updateTemplateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.UpdateTemplate(user, updateTemplateDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// DeleteTemplate godoc
//
//	@Summary	Delete template
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.MessageWrapperDto   "Delete information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) DeleteTemplate(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err = h.service.DeleteTemplate(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]string{
			"message": "Template is removed",
		},
	)
}

// GetDetail godoc
//
//	@Summary	Template detail
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	dto.ResponseTemplateDto   "Template information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates/{guid} [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) GetDetail(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	detail, err := h.service.GetDetail(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, detail)
}

// ApproveTemplate godoc
//
//	@Summary	Approve template
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Param		Approve template  body		dto.TwilioTemplateApprovalDto				true	"Template detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseTemplateDto ]   "Template detail"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates/approve/{guid} [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *TemplateHandler) ApproveTemplate(c echo.Context) error {
	var approveDto dto.TwilioTemplateApprovalDto
	if err := c.Bind(&approveDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&approveDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.ApproveTemplate(user, approveDto)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, data)

}

// Sync godoc
//
//	@Summary	Sync templates
//	@Tags		Templates
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[any]   "Sync template"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/templates/sync [get]
func (h *TemplateHandler) Sync(c echo.Context) error {

	err := h.service.SyncTemplateStatuses()
	if err != nil {
		log.Printf("Webhook cannot work correctly: %s\n", err.Error())
		return response.Error(c, err)
	}

	return response.Success(
		c, map[string]string{
			"message": "Sync is successful",
		},
	)
}
