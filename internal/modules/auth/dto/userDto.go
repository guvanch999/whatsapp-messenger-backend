package dto

type UserLogin struct {
	Email    string `json:"email" form:"email" validate:"required,email" example:"test@gmail.com"`
	Password string `json:"password" form:"password" validate:"required,gte=6" example:"123456"`
}
