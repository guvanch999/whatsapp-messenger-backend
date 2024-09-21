package dto

import (
	"github.com/google/uuid"
	"time"
)

type ResponseOrganizationDto struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
