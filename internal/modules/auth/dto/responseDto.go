package dto

import (
	"github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/nedpals/supabase-go"
)

type AuthenticatedDetails struct {
	AccessToken          string      `json:"access_token"`
	TokenType            string      `json:"token_type"`
	ExpiresIn            int         `json:"expires_in"`
	RefreshToken         string      `json:"refresh_token"`
	User                 models.User `json:"user"`
	ProviderToken        string      `json:"provider_token"`
	ProviderRefreshToken string      `json:"provider_refresh_token"`
}

func (ad *AuthenticatedDetails) FromSupabaseAuthDetails(auth *supabase.AuthenticatedDetails) {
	ad.AccessToken = auth.AccessToken
	ad.TokenType = auth.TokenType
	ad.ExpiresIn = auth.ExpiresIn
	ad.RefreshToken = auth.RefreshToken
	ad.ProviderToken = auth.ProviderToken
	ad.ProviderRefreshToken = auth.ProviderRefreshToken
	authUser := models.User{}
	authUser.FromSupabaseUser(&auth.User)
	ad.User = authUser
}
