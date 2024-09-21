package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
)

func CheckAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(models.UserDetail)
		if user.Role != enums.Admin {
			return response.Error(c, &exceptions.AccessDenied{})
		}
		return next(c)
	}
}
