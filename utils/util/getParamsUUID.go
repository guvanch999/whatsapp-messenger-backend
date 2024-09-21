package util

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
)

func GetParamsUUID(c echo.Context, name string) (uuid.UUID, error) {
	userGuid, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}

	return userGuid, nil
}
