package dto

import (
	"github.com/google/uuid"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"time"
)

type ResponseList struct {
	Id        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Contacts  []UserContact `json:"contacts,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
