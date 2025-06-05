package models

type User struct {
	Id      string
	Name    string
	Email   string
	Picture string
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
