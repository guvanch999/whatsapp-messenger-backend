package dto

import "github.com/google/uuid"

type OrganizationDto struct {
	Name string `json:"name" validate:"required,gt=0"`
}

type UpdateOrganization struct {
	Id   uuid.UUID `json:"guid" param:"guid" validate:"required,uuid4"`
	Name string    `json:"name" validate:"required,gt=0"`
}
