package interfaces

import (
	"context"

	"github.com/johnquangdev/oauth2/usecase/models"
)

type Auth interface {
	GetGoogleAuthURL(string) (string, error)
	GenerateRandomState() (string, error)
	ExchangeCodeForToken(context.Context, models.ExchangeTokenRequest) (models.OAuthResponse, error)
}

type UseCase interface {
	Auth() Auth
}
