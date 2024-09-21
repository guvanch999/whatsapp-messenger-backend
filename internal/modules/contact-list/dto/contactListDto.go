package dto

import "github.com/google/uuid"

type ContactListDto struct {
	Name        string      `json:"name" validate:"required,gt=0"`
	ContactList []uuid.UUID `json:"contact_list" validate:"required,unique,dive,required,uuid4"`
}

type UpdateContactListNameDto struct {
	Id   uuid.UUID `json:"guid" param:"guid" validate:"required,uuid4"`
	Name string    `json:"name" validate:"required,gt=0"`
}

type ContactListUpdateDto struct {
	Id          uuid.UUID   `json:"guid" param:"guid" validate:"required,uuid4"`
	ContactList []uuid.UUID `json:"contact_list" validate:"required,unique,gt=0,dive,required,uuid4"`
}
