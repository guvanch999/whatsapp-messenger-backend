package models

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/dto"
	"time"
)

type UserContact struct {
	Id          uuid.UUID              `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	UserID      uuid.UUID              `json:"user_id"`
	Name        string                 `json:"name"`
	PhoneNumber string                 `json:"phone_number"`
	Email       string                 `json:"email"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"serializer:json"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func (*UserContact) TableName() string {
	return "user_contacts"
}

func (s *UserContact) FromDto(contactDto dto.UserContactDto) {
	s.Name = contactDto.Name
	s.PhoneNumber = contactDto.PhoneNumber
	s.Email = contactDto.Email
	s.Metadata = contactDto.Metadata
}
func (s *UserContact) FromUpdateDto(contactDto dto.UpdateContactDto) {
	s.Id = contactDto.Id
	s.Name = contactDto.Name
	s.PhoneNumber = contactDto.PhoneNumber
	s.Email = contactDto.Email
	s.Metadata = contactDto.Metadata
}

func (s *UserContact) ToResponseDto() *dto.ContactResponse {
	return &dto.ContactResponse{
		Id:          s.Id,
		Name:        s.Name,
		PhoneNumber: s.PhoneNumber,
		Email:       s.Email,
		Metadata:    s.Metadata,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
