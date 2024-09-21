package dto

import "github.com/medium-messenger/messenger-backend/utils/enums"

type SendMessageResponse struct {
	PhoneNumber  string                  `json:"phone_number"`
	Status       enums.MessageSendStatus `json:"status"`
	ErrorMessage string                  `json:"error_message,omitempty"`
}
