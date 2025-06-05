package utils

import (
	"context"
	"errors"
	"time"

	configGoogle "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// struct implement interfaces 1
type ServiceOauthGoogle struct {
	config *configGoogle.Config
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
}

// interfaces 1
type Oauth2 interface {
	ChangeCodeToToken(context.Context, string) (Token, error)
	GenerateAuthURL(state string) string
}

var googleOauthConfig = &configGoogle.Config{
	RedirectURL:  "http://localhost:8080/v1/oauth/callback",
	ClientID:     "487291772648-olmn5125kgmujjerkru6ihjh3nrbefa4.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func (g *ServiceOauthGoogle) ChangeCodeToToken(ctx context.Context, code string) (Token, error) {
	accessToken, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return Token{}, err
	}
	return Token{
		AccessToken: accessToken.AccessToken,
		Expiry:      accessToken.Expiry,
		TokenType:   accessToken.TokenType,
	}, nil
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

func NewGoogleOAuthService() *ServiceOauthGoogle {
	config := &configGoogle.Config{
		RedirectURL:  "http://localhost:8080/v1/auth/callback",
		ClientID:     "487291772648-olmn5125kgmujjerkru6ihjh3nrbefa4.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &ServiceOauthGoogle{
		config: config,
	}
}

func (s *ServiceOauthGoogle) GenerateAuthURL(state string) string {
	return s.config.AuthCodeURL(state, configGoogle.AccessTypeOffline)
}

func (s *ServiceOauthGoogle) ExchangeCodeForToken(ctx context.Context, code string) (OAuthResponse, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return OAuthResponse{}, err
	}

	response := OAuthResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry.String(),
	}

	return response, nil
}

func (s *ServiceOauthGoogle) ValidateConfig() error {
	if s.config.ClientID == "" {
		return errors.New("GOOGLE_CLIENT_ID environment variable is required")
	}
	if s.config.ClientSecret == "" {
		return errors.New("GOOGLE_CLIENT_SECRET environment variable is required")
	}
	return nil
}
