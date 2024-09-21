package service

import (
	"context"
	"github.com/medium-messenger/messenger-backend/internal/modules/auth/dto"
	supa "github.com/nedpals/supabase-go"
	"log"
)

// AuthService TODO Search method that allows us to handle different error types
// from Supabase, such as 'Not Found' errors, 'Client' errors, and others.
type AuthService struct {
	sp *supa.Client
}

func NewAuthService(supabase *supa.Client) *AuthService {

	return &AuthService{
		sp: supabase,
	}
}

func (s *AuthService) Auth(loginDto *dto.UserLogin) (*supa.User, error) {
	ctx := context.Background()
	user, err := s.sp.Auth.SignUp(
		ctx, supa.UserCredentials{
			Email:    loginDto.Email,
			Password: loginDto.Password,
			Data: map[string]string{
				"role": "user",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(loginDto *dto.UserLogin) (*supa.AuthenticatedDetails, error) {
	ctx := context.Background()
	user, err := s.sp.Auth.SignIn(
		ctx, supa.UserCredentials{
			Email:    loginDto.Email,
			Password: loginDto.Password,
		},
	)
	if err != nil {
		log.Printf("Error on sign in: %s\n", err.Error())
		return nil, err
	}

	return user, nil
}

func (s *AuthService) CheckToken(token string) (*supa.User, error) {
	cnt := context.Background()
	user, err := s.sp.Auth.User(cnt, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) RefreshToken(refTknDto *dto.RefreshTokenDto) (*supa.AuthenticatedDetails, error) {
	ctx := context.Background()
	user, err := s.sp.Auth.RefreshUser(ctx, refTknDto.AccessToken, refTknDto.RefreshToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ChangePassword(token string, newPassword string) (*supa.User, error) {
	updateData := map[string]interface{}{
		"password": newPassword,
	}
	user, err := s.sp.Auth.UpdateUser(context.Background(), token, updateData)
	if err != nil {
		return nil, err
	}
	return user, nil
}
