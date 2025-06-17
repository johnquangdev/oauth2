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
	CreateSesion(*models.Session) error
	IsUserExists(string) (bool, error)
	BlockedUserByUserID(context.Context, uuid.UUID) error
}

type Redis interface {
	AddBackList(userID string, token string, duration time.Duration) error
	GetRefreshToken(userID string) (string, error)
	DeleteRefreshToken(userID string) error
	IsTokenBlacklisted(tokenID uuid.UUID) (bool, error)
}

type Repo interface {
	Auth() Auth
	Redis() Redis
}
