package dto

import "github.com/medium-messenger/messenger-backend/utils/enums"

type ChangeUserInfoDto struct {
	Role enums.Role `json:"role" validate:"required,oneof=user admin"`
}
