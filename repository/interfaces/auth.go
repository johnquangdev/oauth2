package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/repository/models"
)

type Auth interface {
	GetUserByUserId(context.Context, uuid.UUID) (*models.User, error)
	CreateUser(*models.User) error
	CreateSession(*models.Session) error
	UserExists(string) (bool, error)
	BlockedUserByUserID(context.Context, uuid.UUID) error
	GetUserByProviderAndProviderId(context.Context, string, string) (*models.User, error)
}

type Redis interface {
	AddBackList(userID string, token string, duration time.Duration) error
	IsTokenBlacklisted(tokenID uuid.UUID) (bool, error)
	CreateRecord(userId uuid.UUID, accessToken string, accessTokenTimeLife time.Duration) error
}

type Repo interface {
	Auth() Auth
	Redis() Redis
}
