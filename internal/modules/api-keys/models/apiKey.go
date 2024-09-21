package models

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/dto"
	"time"
)

type ApiKey struct {
	Id        uuid.UUID `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Hash      string    `json:"hash"`
	Encoded   string    `json:"encoded"`
	CreatedAt time.Time `json:"created_at"`
}

func (k *ApiKey) ToResponseDto() *dto.ApiKeyResponse {
	return &dto.ApiKeyResponse{
		Id:        k.Id,
		Name:      k.Name,
		CreatedAt: k.CreatedAt,
	}
}
