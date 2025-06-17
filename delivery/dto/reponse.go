package dto

import "time"

type UserResponse struct {
	Gmail   string
	Name    string
	Picture string
}

type JwtResponse struct {
	AccessToken           string
	RefreshToken          string
	ExpiresAtAccessToken  time.Time
	ExpiresAtRefreshToken time.Time
}
