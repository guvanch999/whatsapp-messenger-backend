package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/messaging/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/messaging/service"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	_ "github.com/medium-messenger/messenger-backend/utils/util"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService,
	}
}

// SendMessage godoc
//
//	@Summary	Send message to recipients
//	@Tags		Messaging
//	@Accept		json
//	@Produce	json
//	@Param		Messaging 		body		dto.SendMessageDto				true	"Messaging information"
//	@Success	200				{object}	util.ListMessageDataWrapperDto[[]dto.SendMessageResponse]		"Send message information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/messages [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *MessageHandler) SendMessage(c echo.Context) error {
	var sendMessageDto dto.SendMessageDto
	if err := c.Bind(&sendMessageDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&sendMessageDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	process, err := h.service.SendMessages(user, sendMessageDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"message": "Messages processed",
			"list":    process,
		},
	)
}

// SendMessageList godoc
//
//	@Summary	Send message to groups
//	@Tags		Messaging
//	@Accept		json
//	@Produce	json
//	@Param		Messaging 		body		dto.SendMessageToListDto				true	"Messaging information"
//	@Success	200				{object}	util.ListMessageDataWrapperDto[[]dto.SendMessageResponse]		"Send message information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/messages/to-list [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *MessageHandler) SendMessageList(c echo.Context) error {
	var sendMessageDto dto.SendMessageToListDto
	if err := c.Bind(&sendMessageDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(&sendMessageDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	user := c.Get("user").(auth.UserDetail)
	process, err := h.service.SendMessageToGroup(user, sendMessageDto)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(
		c, map[string]any{
			"message": "Messages processed",
			"list":    process,
		},
	)
}
