package dto

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"time"
)

type CreateTemplateDto struct {
	Name         string         `json:"name" validate:"required"`
	Content      interface{}    `json:"content" validate:"required"`
	ProviderId   uuid.UUID      `json:"provider_id" validate:"required,uuid4"`
	Platform     enums.Platform `json:"platform" validate:"required,oneof=WhatsApp sms email"`
	ProviderType enums.Provider `json:"provider_type" validate:"required,oneof=twilio plivo"`
}

type UpdateTemplateDto struct {
	Id uuid.UUID `json:"guid" param:"guid" validate:"required,uuid4"`
	CreateTemplateDto
}

type ResponseTemplateDto struct {
	Id             uuid.UUID      `json:"id" `
	Name           string         `json:"name"`
	Content        interface{}    `json:"content"`  // json
	Status         enums.Status   `json:"status"`   // pending | accepted | rejected
	Platform       enums.Platform `json:"platform"` // WhatsApp | Sms |Email
	ProviderType   enums.Provider `json:"provider"` // twilio | plivo
	ProviderId     uuid.UUID      `json:"provider_id"`
	ExternalId     string         `json:"external_id"`
	ExternalStatus enums.Status   `json:"external_status"` //enum
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
