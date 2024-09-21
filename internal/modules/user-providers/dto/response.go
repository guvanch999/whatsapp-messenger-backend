package dto

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"time"
)

type ResponseProviderDto struct {
	Id              uuid.UUID      `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Name            string         `json:"name"`
	FromPhoneNumber string         `json:"from_phone_number"`
	Status          enums.Status   `json:"status"` // inreview | approved |rejected | paused | disabled | unsubmitted
	Type            enums.Provider `json:"type"`   // twilio | plivo
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	WebhookUrl      string         `json:"webhook_url"`
}
