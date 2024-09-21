package response

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"net/http"
)

func GetError(err error) any {
	switch true {
	case errors.Is(err, &exceptions.NotFoundError{}):
		return err.Error()
	case errors.Is(err, &exceptions.BadRequestError{}):
		return map[string]interface{}{
			"message": err.(*exceptions.BadRequestError).Message,
			"error":   err.Error(),
		}
	case errors.Is(err, &exceptions.AuthFailed{}):
		return map[string]interface{}{
			"message": err.(*exceptions.AuthFailed).Message,
			"error":   err.Error(),
		}
	case errors.Is(err, &exceptions.Forbidden{}):
		return map[string]interface{}{
			"message": err.(*exceptions.Forbidden).Message,
			"error":   err.Error(),
		}
	case errors.Is(err, &exceptions.AccessDenied{}):
		return err.Error()
	case errors.Is(err, &exceptions.ResponseError{}):
		return map[string]interface{}{
			"url":     err.(*exceptions.ResponseError).Url,
			"message": err.(*exceptions.ResponseError).Message,
			"error":   err.Error(),
		}
	default:
		return err.Error()
	}
}

func Error(c echo.Context, err error) error {
	switch true {
	case errors.Is(err, &exceptions.NotFoundError{}):
		return echo.ErrNotFound
	case errors.Is(err, &exceptions.BadRequestError{}):
		return c.JSON(
			http.StatusBadRequest, map[string]interface{}{
				"message": err.(*exceptions.BadRequestError).Message,
				"error":   err.Error(),
			},
		)
	case errors.Is(err, &exceptions.AuthFailed{}):
		return c.JSON(
			http.StatusUnauthorized, map[string]interface{}{
				"message": err.(*exceptions.AuthFailed).Message,
				"error":   err.Error(),
			},
		)
	case errors.Is(err, &exceptions.Forbidden{}):
		return c.JSON(
			http.StatusForbidden, map[string]interface{}{
				"message": err.(*exceptions.Forbidden).Message,
				"error":   err.Error(),
			},
		)
	case errors.Is(err, &exceptions.AccessDenied{}):
		return c.JSON(http.StatusNotAcceptable, err.Error())
	case errors.Is(err, &exceptions.ResponseError{}):
		return c.JSON(
			http.StatusBadRequest, map[string]interface{}{
				"url":     err.(*exceptions.ResponseError).Url,
				"message": err.(*exceptions.ResponseError).Message,
				"error":   err.Error(),
			},
		)
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}
