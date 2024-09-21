package dto

import "github.com/google/uuid"

type Recipient struct {
	RecipientId uuid.UUID   `json:"recipient_id" validate:"required,uuid4"`
	Variables   interface{} `json:"template_variables"`
}

type SendMessageDto struct {
	Recipients []Recipient `json:"recipients" validate:"required,gt=0,dive"`
	ProviderId uuid.UUID   `json:"provider_id" validate:"required,uuid4"`
	TemplateId uuid.UUID   `json:"template_id" validate:"required,uuid4"`
}

type SendMessageToListDto struct {
	ProviderId        uuid.UUID   `json:"provider_id" validate:"required,uuid4"`
	TemplateId        uuid.UUID   `json:"template_id" validate:"required,uuid4"`
	TemplateVariables interface{} `json:"template_variables"`
	ContactListId     *uuid.UUID  `json:"contact_list_id" validate:"required,uuid4"`
}
