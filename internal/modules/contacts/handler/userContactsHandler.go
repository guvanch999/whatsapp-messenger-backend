package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/service"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type UserContactsHandler struct {
	service *service.UserContactsService
}

func NewUserContactsHandler(contactsService *service.UserContactsService) *UserContactsHandler {
	return &UserContactsHandler{
		service: contactsService,
	}
}

// GetMyContacts godoc
//
//	@Summary	Get user Contacts
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ContactResponse]   "User Contact list"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) GetMyContacts(c echo.Context) error {
	user := c.Get("user").(models.UserDetail)
	data, err := h.service.GetUserContacts(user.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// GetAllContacts godoc
//
//	@Summary	Get my Contacts
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.ContactResponse]   "All Contact list"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts/all [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) GetAllContacts(c echo.Context) error {
	data, err := h.service.GetAllUserContacts()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}

// GetContactDetail godoc
//
//		@Summary	Get Contact detail
//		@Tags		Contacts
//		@Accept		json
//		@Produce	json
//	 	@Param		guid			path		string							true	"Contact ID"
//		@Success	200				{object}	util.DataWrapperDto[dto.ContactResponse]   "Contact detail"
//		@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//		@Failure	500				{object}	string						"Internal server error"
//		@Router		/user-contacts/{guid} [get]
//		@Security	Bearer
//		@Security	X-API-KEY
func (h *UserContactsHandler) GetContactDetail(c echo.Context) error {
	userGuid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(models.UserDetail)
	data, err := h.service.GetContactDetail(user, userGuid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, data,
	)
}

// AddContact godoc
//
//	@Summary	Add contact
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Param		Contact detail 	body		dto.UserContactDto				true	"Add contact detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.CreatedOrExistResponse]		"Contact information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) AddContact(c echo.Context) error {
	var createDto dto.UserContactDto
	if err := c.Bind(&createDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&createDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(models.UserDetail)
	data, err := h.service.AddUserContact(user, createDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c,
		data,
	)
}

// AddListOfContacts godoc
//
//	@Summary	Add contact list
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Param		Contact detail 	body		dto.UserContactsListDto				true	"Add contact list"
//	@Success	200				{object}	util.DataWrapperDto[dto.CreatedOrExistResponse]		"Contact information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts/list [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) AddListOfContacts(c echo.Context) error {
	var userContactsList dto.UserContactsListDto
	if err := c.Bind(&userContactsList); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&userContactsList); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}

	user := c.Get("user").(models.UserDetail)
	data, err := h.service.AddContactsList(user, userContactsList)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, data,
	)
}

// UpdateContactDetail godoc
//
//	@Summary	Update contact
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Param		guid			path		string							true	"Contact ID"
//	@Param		Contact detail 	body		dto.UpdateContactDto				true	"Update contact detail"
//	@Success	200				{object}	util.DataWrapperDto[dto.ContactResponse]		"Updated contact information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts/{guid} [put]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) UpdateContactDetail(c echo.Context) error {
	var updateDto dto.UpdateContactDto
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
	user := c.Get("user").(models.UserDetail)
	data, err := h.service.UpdateContact(user, updateDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, data,
	)
}

// DeleteContactDetail godoc
//
//	@Summary	Delete contact
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Param		guid			path		string							true	"Contact ID"
//	@Success	200				{object}	util.MessageWrapperDto		"Delete information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts/{guid} [delete]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *UserContactsHandler) DeleteContactDetail(c echo.Context) error {
	contactUuid, err := util.GetParamsUUID(c, "guid")
	if err != nil {
		return response.Error(c, err)
	}
	user := c.Get("user").(models.UserDetail)
	err = h.service.DeleteContact(user, contactUuid)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"message": "Contact is removed",
		},
	)
}

// ValidateNumber godoc
//
//	@Summary	Validate phone number
//	@Tags		Contacts
//	@Accept		json
//	@Produce	json
//	@Param		Numbers  	 	body		dto.ValidateNumbersDto				true	"List of numbers"
//	@Success	200				{object}	util.ListDataWrapperDto[[]dto.NumberValidateResponse]		"Validated result"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/user-contacts/validate [post]
//	@Security	Bearer
func (h *UserContactsHandler) ValidateNumber(c echo.Context) error {
	var validateNumberDto dto.ValidateNumbersDto
	if err := c.Bind(&validateNumberDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&validateNumberDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	data := h.service.ValidateNumbers(validateNumberDto.Numbers)

	return response.Success(
		c, map[string]any{
			"list": data,
		},
	)
}
