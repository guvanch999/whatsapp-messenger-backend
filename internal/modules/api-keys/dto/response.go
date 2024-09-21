package dto

import (
	"github.com/google/uuid"
	"time"
)

type ApiKeyResponse struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ApiKeyDetailResponse struct {
	ApiKey string `json:"api_key"`
}
