package models

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Avatar     string    `json:"avatar"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"provider_id"`
}

type ExchangeTokenRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}
