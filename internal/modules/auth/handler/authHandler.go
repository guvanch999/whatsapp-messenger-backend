package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/config"
	"github.com/medium-messenger/messenger-backend/internal/modules/auth/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/auth/service"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/response"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
	conf    *config.Schema
}

func NewAuthHandler(
	service *service.AuthService,
	conf *config.Schema,
) *AuthHandler {
	return &AuthHandler{
		service: service,
		conf:    conf,
	}
}

// RegisterHandler godoc
//
//	@Summary	Register
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		credentials 	body		dto.UserLogin				true	"Register credentials"
//	@Success	200				{object}	models.User					"Registered user information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/auth/register [post]
func (h *AuthHandler) RegisterHandler(c echo.Context) error {
	u := new(dto.UserLogin)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(u); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	data, err := h.service.Auth(u)
	if err != nil {
		return response.Error(c, err)
	}

	userDetail := models.User{}
	userDetail.FromSupabaseUser(
		data,
	)
	return c.JSON(http.StatusOK, userDetail)
}

// LoginHandler godoc
//
//	@Summary	Login
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		credentials 	body		dto.UserLogin				true	"Login credentials"
//	@Success	200				{object}	dto.AuthenticatedDetails						"Login user information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/auth/login [post]
func (h *AuthHandler) LoginHandler(c echo.Context) error {
	u := new(dto.UserLogin)
	if err := c.Bind(u); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(u); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	data, err := h.service.Login(u)
	if err != nil {
		return response.Error(c, err)
	}
	authDetail := dto.AuthenticatedDetails{}
	authDetail.FromSupabaseAuthDetails(data)

	return c.JSON(http.StatusOK, authDetail)
}

// GetMe godoc
//
//	@Summary	Get my profile detail
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}	models.User	                "Login user information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/auth/me [get]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *AuthHandler) GetMe(c echo.Context) error {
	user := c.Get("user").(models.UserDetail)
	return c.JSON(
		http.StatusOK, user.User,
	)
}

// RefreshHandler godoc
//
//	@Summary	Refresh token
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		token credentials 	body		dto.RefreshTokenDto				true	"Token credentials"
//	@Success	200				{object}	dto.AuthenticatedDetails						"Token information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/auth/refresh [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *AuthHandler) RefreshHandler(c echo.Context) error {
	refreshDto := new(dto.RefreshTokenDto)
	if err := c.Bind(refreshDto); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(refreshDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}

	user, err := h.service.RefreshToken(refreshDto)
	if err != nil {
		return err
	}

	authDetail := dto.AuthenticatedDetails{}
	authDetail.FromSupabaseAuthDetails(user)
	return c.JSON(http.StatusOK, authDetail)
}

// ChangePassword godoc
//
//	@Summary	Change password
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		credentials 	body		dto.ChangePasswordDto				true	"Change password credentials"
//	@Success	200				{object}	models.User						"Change password information"
//	@Failure	400				{object}	exceptions.BadRequestError	"Bad request"
//	@Failure	500				{object}	string						"Internal server error"
//	@Router		/auth/change-password [post]
//	@Security	Bearer
//	@Security	X-API-KEY
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	changePasswordDto := new(dto.ChangePasswordDto)
	if err := c.Bind(changePasswordDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}
	if err := c.Validate(changePasswordDto); err != nil {
		return response.Error(
			c, &exceptions.BadRequestError{
				Message: err.Error(),
			},
		)
	}

	token, ok := c.Get("token").(string)
	if !ok {
		return response.Error(
			c, &exceptions.Forbidden{
				Message: "This function is not available with api key authorization",
			},
		)
	}
	user, err := h.service.ChangePassword(token, changePasswordDto.Password)
	if err != nil {
		return response.Error(c, err)
	}

	userDetail := models.User{}
	userDetail.FromSupabaseUser(
		user,
	)
	return c.JSON(http.StatusOK, userDetail)
}
