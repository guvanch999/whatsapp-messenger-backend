package models

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/dto"
	"time"
)

type Organization struct {
	Id        uuid.UUID `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	OwnerId   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (o *Organization) TableName() string {
	return "organizations"
}

func (o *Organization) ToResponseDto() *dto.ResponseOrganizationDto {
	return &dto.ResponseOrganizationDto{
		Id:        o.Id,
		Name:      o.Name,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
