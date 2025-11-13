package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProviderGoogle = "google"
	ProviderGitHub = "github"
	//ProviderFacebook = "facebook"
)

const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusBlocked = "blocked"
	StatusBanned  = "banned"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Avatar     string    `json:"avatar"`
	Provider   string    `json:"provider"`
	ProviderId string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ExchangeTokenRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type TokenJwt struct {
	AccessToken           string        `json:"access_token,omitempty"`
	RefreshToken          string        `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  time.Duration `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt time.Duration `json:"refresh_token_expires_at,omitempty"`
}
