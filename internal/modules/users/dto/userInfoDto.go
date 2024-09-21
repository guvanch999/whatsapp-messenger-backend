package dto

import (
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
)

type UserInfoDto struct {
	User models.User `json:"user"`
	Role enums.Role  `json:"role"`
}
