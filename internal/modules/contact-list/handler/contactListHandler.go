package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type ContactListHandler struct {
	service *service.ContactListService
}

func NewContactListHandler(listService *service.ContactListService) *ContactListHandler {
	return &ContactListHandler{
		service: listService,
	}
}

// GetUserContactLists godoc
//
//	@Summary	Get contact group
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ResponseList]   "User Contacts groups"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/contact-list [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) GetUserContactLists(c echo.Context) error {
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.GetUserLists(user.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}
func (h *ContactListHandler) GetAllContactList(c echo.Context) error {
	data, err := h.service.GetAllLists()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// CreateContactList godoc
//
//	@Summary	Add contact group
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Param		Group detail 	body		dto.ContactListDto				true	"Contact group"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseList]		"Contact group information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/contact-list [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) CreateContactList(c echo.Context) error {
	var contactListDto dto.ContactListDto
	if err := c.Bind(&contactListDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)

	}
	if err := c.Validate(&contactListDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.AddContactList(user, contactListDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// AddContactToList godoc
//
//	@Summary	Add contact to list
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Param		Contacts 		body		dto.ContactListUpdateDto				true	"Contact list"
//	@Success	200				{object}	util.MessageWrapperDto			"Status message"
//	@Failure	400				{object}	exceptions.BadRequestError		"Bad request"
//	@Failure	500				{object}	string							"Internal server error"
//	@Router		/contact-list/add-contact/{guid} [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) AddContactToList(c echo.Context) error {
	var updateDto dto.ContactListUpdateDto
	if err := c.Bind(&updateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&updateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err := h.service.AddContactToList(user, updateDto)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(
		c, map[string]any{
			"message": "Contacts list updated",
		},
	)
}

// UpdateContactListName godoc
//
//	@Summary	Update contact group name
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Param		Contact group detail 		body		dto.UpdateContactListNameDto				true	"Contact group detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseList]			"Contact group detail"
//	@Failure	400				{object}	exceptions.BadRequestError		"Bad request"
//	@Failure	500				{object}	string							"Internal server error"
//	@Router		/contact-list/name/{guid} [put]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) UpdateContactListName(c echo.Context) error {
	var updateListDto dto.UpdateContactListNameDto
	if err := c.Bind(&updateListDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&updateListDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	data, err := h.service.UpdateContactListName(user, updateListDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, data)
}

// DeleteContactList godoc
//
//	@Summary	Delete contact group
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.MessageWrapperDto			"Status message"
//	@Failure	400				{object}	exceptions.BadRequestError		"Bad request"
//	@Failure	500				{object}	string							"Internal server error"
//	@Router		/contact-list/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) DeleteContactList(c echo.Context) error {
	guid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err = h.service.DeleteList(user, guid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]string{
			"message": "Contact list is removed",
		},
	)
}

// GetDetail godoc
//
//	@Summary	Get group detail
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.DataWrapperDto[dto.ResponseList]			"Group detail"
//	@Failure	400				{object}	exceptions.BadRequestError		"Bad request"
//	@Failure	500				{object}	string							"Internal server error"
//	@Router		/contact-list/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) GetDetail(c echo.Context) error {
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

// DeleteContactFromList godoc
//
//	@Summary	Remove contact to group
//	@Tags		Contact group
//	@Accept		json
//	@Produce	json
//	@Param		Contacts 		body		dto.ContactListUpdateDto		true	"Contact group"
//	@Success	200				{object}	util.MessageWrapperDto			"Status message"
//	@Failure	400				{object}	exceptions.BadRequestError		"Bad request"
//	@Failure	500				{object}	string							"Internal server error"
//	@Router		/contact-list/remove-contact/{guid} [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *ContactListHandler) DeleteContactFromList(c echo.Context) error {
	var updateDto dto.ContactListUpdateDto
	if err := c.Bind(&updateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&updateDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	err := h.service.RemoveContactFromList(user, updateDto)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(
		c, map[string]any{
			"message": "Contacts list updated",
		},
	)
}
