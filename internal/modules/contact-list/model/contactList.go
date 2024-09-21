package model

import (
	"github.com/google/uuid"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contact-list/dto"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contacts/models"
	"time"
)

type ContactList struct {
	Id        uuid.UUID     `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID     `json:"user_id"`
	Name      string        `json:"name"`
	Contacts  []UserContact `gorm:"many2many:contact_list_contacts;"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (*ContactList) TableName() string {
	return "contact_list"
}

func (s *ContactList) ToResponseDto() *ResponseList {
	return &ResponseList{
		Id:        s.Id,
		Name:      s.Name,
		Contacts:  s.Contacts,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
