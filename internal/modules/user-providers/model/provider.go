package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/config"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/dto"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"time"
)

type UserProvider struct {
	Id                  uuid.UUID      `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID              uuid.UUID      `json:"user_id"`
	Name                string         `json:"name"`
	ProviderCredentials string         `json:"provider_credentials"`
	FromPhoneNumber     string         `json:"from_phone_number"`
	Status              enums.Status   `json:"status"` // inreview | approved |rejected | paused | disabled | unsubmitted
	Type                enums.Provider `json:"type"`   // twilio | plivo
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (*UserProvider) TableName() string {
	return "user_providers"
}

func (p *UserProvider) ToResponseDto(cnf *config.Schema) *dto.ResponseProviderDto {
	return &dto.ResponseProviderDto{
		Id:              p.Id,
		Name:            p.Name,
		FromPhoneNumber: p.FromPhoneNumber,
		Status:          p.Status,
		Type:            p.Type,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		WebhookUrl:      fmt.Sprintf("%s/v1/templates/webhook/%s", cnf.AppUrl, p.Id),
	}
}

type TwilioCred struct {
	TwilioAccountSid          string `json:"twilio_account_sid"`
	TwilioAuthToken           string `json:"twilio_auth_token"`
	TwilioMessagingServiceSid string `json:"twilio_messaging_service_sid"`
}
