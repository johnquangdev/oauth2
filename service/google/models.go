package utils

import "time"

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
}

type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type AuthURLResponse struct {
	AuthURL     string `json:"auth_url"`
	State       string `json:"state"`
	RedirectURL string `json:"redirect_url"`
	ClientID    string `json:"client_id"`
}

type UserResponse struct {
	Id      string
	Email   string
	Name    string
	Picture string
}
