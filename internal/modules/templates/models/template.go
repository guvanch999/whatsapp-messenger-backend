package models

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/model"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	openapi "github.com/twilio/twilio-go/rest/content/v1"
	"time"
)

type Template struct {
	Id           uuid.UUID           `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID       uuid.UUID           `json:"user_id"`
	Name         string              `json:"name"`
	Content      interface{}         `json:"content" gorm:"serializer:json"` // json
	Status       enums.Status        `json:"status"`                         // inreview | approved |rejected | paused | disabled | unsubmitted
	Platform     enums.Platform      `json:"platform"`                       // WhatsApp | Sms |Email
	ProviderType enums.Provider      `json:"provider_type"`                  // twilio | plivo
	ProviderId   uuid.UUID           `json:"provider_id"`
	Provider     *model.UserProvider `json:"provider,omitempty" gorm:"foreignKey:provider_id;references:id;constraint:OnDelete:set null;"`
	ExternalId   string              `json:"external_id"`
	NextCheck    time.Time           `json:"next_check" gorm:"default:null"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func (t *Template) GetContent() (*openapi.ContentCreateRequest, error) {
	jsonBytes, err := json.Marshal(t.Content)
	if err != nil {
		return nil, fmt.Errorf("content is not valid json: %s", err.Error())
	}
	content := new(openapi.ContentCreateRequest)
	if err = json.Unmarshal(jsonBytes, &content); err != nil {
		return nil, fmt.Errorf("content is not valid: %s", err.Error())
	}

	return content, nil
}

func (t *Template) TableName() string {
	return "templates"
}

func (t *Template) FromDto(templateDto *dto.CreateTemplateDto) {
	t.Content = templateDto.Content
	t.Name = templateDto.Name
	t.Status = enums.Unsubmitted
	t.Platform = templateDto.Platform
	t.ProviderType = templateDto.ProviderType
	t.ProviderId = templateDto.ProviderId
}

func (t *Template) FromUpdateDto(templateDto *dto.UpdateTemplateDto) {
	t.Id = templateDto.Id
	t.Content = templateDto.Content
	t.Name = templateDto.Name
	t.Status = enums.Unsubmitted
	t.Platform = templateDto.Platform
	t.ProviderType = templateDto.ProviderType
	t.ProviderId = templateDto.ProviderId
}

func (t *Template) ToResponseDto() *dto.ResponseTemplateDto {
	return &dto.ResponseTemplateDto{
		Id:           t.Id,
		Name:         t.Name,
		Content:      t.Content,
		Status:       t.Status,
		Platform:     t.Platform,
		ProviderType: t.ProviderType,
		ProviderId:   t.ProviderId,
		ExternalId:   t.ExternalId,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
