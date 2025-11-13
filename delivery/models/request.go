package models

type LoginOauth2 struct {
	Code string `json:"code" validate:"required"`
}

type GetUserInfo struct {
	AccessToken string `json:"accesstoken" validate:"required"`
}

type User struct {
	Id         string `json:"id"`
	Gmail      string `json:"gmail" validate:"omitempty"`
	Name       string `json:"name" validate:"omitempty"`
	Phone      string `json:"phone"`
	Picture    string `json:"picture"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}

type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
