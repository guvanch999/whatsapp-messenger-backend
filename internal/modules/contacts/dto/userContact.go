package dto

import "github.com/google/uuid"

type UserContactDto struct {
	Name        string                 `json:"name" validate:"required,gte=1"`
	PhoneNumber string                 `json:"phone_number" validate:"required,phone_number"`
	Email       string                 `json:"email" validate:"required,email"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateContactDto struct {
	Id          uuid.UUID              `json:"guid" param:"guid" validate:"required,uuid4"`
	Name        string                 `json:"name" validate:"required,gte=1"`
	PhoneNumber string                 `json:"phone_number" validate:"required,phone_number"`
	Email       string                 `json:"email" validate:"required,email"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UserContactsListDto struct {
	UserContactsList []UserContactDto `json:"user_contacts_list" validate:"required,gt=0,dive"`
}

type ValidateNumbersDto struct {
	Numbers []string `json:"numbers" validate:"required,gt=0"`
}
