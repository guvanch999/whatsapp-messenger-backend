package dto

type ChangePasswordDto struct {
	Password string `json:"password" validate:"required,gte=6"`
}
