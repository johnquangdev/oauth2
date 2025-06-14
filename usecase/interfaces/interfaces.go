package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	rModels "github.com/johnquangdev/oauth2/repository/models"
	"github.com/johnquangdev/oauth2/usecase/models"
)

type Auth interface {
	GetGoogleAuthURL(string) (string, error)
	GenerateRandomState() (string, error)
	ExchangeCodeForToken(context.Context, models.ExchangeTokenRequest) (models.OAuthResponse, error)
	CreateUser(context.Context, rModels.User) error
	CreateSesion(context.Context, rModels.Session) error
	GetUserByEmail(string) (*models.User, error)
	GetUserInfoGoogle(string) (models.User, error)
	IsUserExists(string) error
	Logout(context.Context, uuid.UUID) error
	SaveRefreshToken(string, string, time.Duration) error
}

type UseCase interface {
	Auth() (Auth, error)
}
