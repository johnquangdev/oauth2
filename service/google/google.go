package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	configGoogle "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// struct implement interfaces 1
type ServiceOauthGoogle struct {
	config *configGoogle.Config
}

type Soauth2 struct {
	googleOauth *oauth2.UserinfoService
}

// interfaces 1
type Oauth2 interface {
	ChangeCodeToToken(context.Context, string) (Token, error)
	GenerateAuthURL(state string) string
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

func Newoauth2() (*Soauth2, error) {
	service, err := oauth2.NewService(context.Background(), option.WithCredentialsFile("keyfile.json"))
	if err != nil {
		return &Soauth2{}, err
	}
	userv2 := oauth2.NewUserinfoService(service)
	return &Soauth2{
		googleOauth: userv2,
	}, nil
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

func (s Soauth2) GetUserInfoGoogle(accessToken string) (UserResponse, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return UserResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserResponse{}, fmt.Errorf("google API error: %s", resp.Status)
	}

	var userInfo UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return UserResponse{}, err
	}

	return userInfo, nil
}

//test sau nhé------
// func (s Soauth2) LGetUserInfo(accessToken string) (UserResponse, error) {
// 	userinfo, err := s.googleOauth.Get().Do(googleapi.QueryParameter("Authorization", "Bearer "+accessToken))
// 	if err != nil {
// 		return UserResponse{}, err
// 	}
// 	return UserResponse{
// 		Email:   userinfo.Email,
// 		Name:    userinfo.GivenName,
// 		Picture: userinfo.Picture,
// 	}, nil
// }
