package dto

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/utils/enums"
)

type TwilioTemplateApprovalDto struct {
	Id       uuid.UUID `json:"guid" param:"guid" validate:"required,uuid4"`
	Name     string    `json:"name" validate:"required"`
	Category string    `json:"category" validate:"required"`
}

type TwilioWebhookDto struct {
	AccountSid           string `json:"AccountSid" form:"AccountSid"`
	AppMonitorTriggerSid string `json:"AppMonitorTriggerSid" form:"AppMonitorTriggerSid"`
	CurrentValue         int    `json:"CurrentValue" form:"CurrentValue"`
	DateFired            string `json:"DateFired" form:"DateFired"`
	Description          string `json:"Description" form:"Description"`
	ErrorCode            int    `json:"ErrorCode" form:"ErrorCode"`
	IdempotencyToken     string `json:"IdempotencyToken" form:"IdempotencyToken"`
	Log                  string `json:"Log" form:"Log"`
	TimePeriod           string `json:"TimePeriod" form:"TimePeriod"`
	TriggerValue         string `json:"TriggerValue" form:"TriggerValue"`
}

type TwilioApprovalRequestDto struct {
	AllowCategoryChange bool         `json:"allow_category_change"`
	Category            string       `json:"category"`
	ContentType         string       `json:"content_type"`
	Flows               interface{}  `json:"flows"`
	Name                string       `json:"name"`
	RejectionReason     string       `json:"rejection_reason"`
	Status              enums.Status `json:"status"`
}
