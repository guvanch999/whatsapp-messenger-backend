package dto

import (
	"github.com/google/uuid"
	"time"
)

type ContactResponse struct {
	Id          uuid.UUID              `json:"id,omitempty"`
	Name        string                 `json:"name"`
	PhoneNumber string                 `json:"phone_number"`
	Email       string                 `json:"email"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type NumberValidateResponse struct {
	Number       string `json:"number"`
	IsValid      bool   `json:"is_valid"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type CreatedOrExistResponse struct {
	Created []ContactResponse `json:"created"`
	Exist   []ContactResponse `json:"exist"`
}
