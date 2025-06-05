package impl

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/usecase/models"
	"github.com/johnquangdev/oauth2/utils"
)

// New code
type OAuthUsecase struct {
	oauthService *utils.ServiceOauthGoogle
}

type GetAuthURLRequest struct {
	CustomState string `json:"custom_state,omitempty"`
}

func NewOAuthUsecase() interfaces.Auth {
	return &OAuthUsecase{
		oauthService: utils.NewGoogleOAuthService(),
	}
}

func (uc *OAuthUsecase) GetGoogleAuthURL(state string) (string, error) {
	// Generate auth URL
	response := uc.oauthService.GenerateAuthURL(state)
	return response, nil
}

func (uc *OAuthUsecase) ExchangeCodeForToken(ctx context.Context, req models.ExchangeTokenRequest) (models.OAuthResponse, error) {
	// Validate input
	if req.Code == "" {
		return models.OAuthResponse{}, errors.New("authorization code is required")
	}

	// Validate OAuth config
	if err := uc.oauthService.ValidateConfig(); err != nil {
		return models.OAuthResponse{}, err
	}

	// Exchange code for token
	tokenResponse, err := uc.oauthService.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return models.OAuthResponse{}, errors.New("failed to exchange code for token: " + err.Error())
	}

	return models.OAuthResponse{
		AccessToken:  tokenResponse.AccessToken,
		TokenType:    tokenResponse.TokenType,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}, nil
}

func (uc *OAuthUsecase) GenerateRandomState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
