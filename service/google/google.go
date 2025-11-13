package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/johnquangdev/oauth2/utils"
	configGoogle "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// struct implement interfaces 1
type ServiceOauthGoogle struct {
	config *configGoogle.Config
	cfg    utils.Config
}

// // interfaces
type Oauth2Custom interface {
	ChangeCodeToToken(context.Context, string) (string, error)
	GenerateAuthURL(state string) string
	GetUserInfoGoogle(accessToken string) (UserInfoResp, error)
}

func NewGoogleOAuthService(cfg utils.Config) (*ServiceOauthGoogle, error) {
	// Parse scopes from config
	scopes := []string{}
	if cfg.Scopes_Google != "" {
		for _, s := range strings.Split(cfg.Scopes_Google, ",") {
			scopes = append(scopes, strings.TrimSpace(s))
		}
	}

	config := &configGoogle.Config{
		RedirectURL:  cfg.RedirectUrl_Google,
		ClientID:     cfg.ClientId_Google,
		ClientSecret: cfg.ClientSecret_Google,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
	return &ServiceOauthGoogle{
		config: config,
		cfg:    cfg,
	}, nil
}

func (s *ServiceOauthGoogle) GenerateAuthURL() (string, error) {
	state, err := GenerateRandomState()
	if err != nil {
		return "", err
	}
	return s.config.AuthCodeURL(state, configGoogle.AccessTypeOffline), nil
}

func (s *ServiceOauthGoogle) GetUserInfoGoogle(accessToken string) (*UserInfoResp, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var UserInfo UserInfoResp
	if err := json.NewDecoder(resp.Body).Decode(&UserInfo); err != nil {
		return nil, fmt.Errorf("decode user info: %w", err)
	}
	UserInfo.Provider = "google"
	return &UserInfo, nil
}

func (o ServiceOauthGoogle) ChangeCodeToToken(ctx context.Context, code string) (string, string, error) {
	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("exchange code failed by err: %w", err)
	}
	if !token.Valid() {
		return "", "", fmt.Errorf("invalid token received from Google")
	}
	idToken := token.Extra("id_token").(string)
	return token.AccessToken, idToken, nil
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
