package dto

type ApiKeyDto struct {
	Name   string `json:"name" validate:"required,gt=0"`
	ApiKey string `json:"api_key" validate:"required,gt=10,printascii"`
}
