package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, map[string]any{
		"data": data,
	})
}
