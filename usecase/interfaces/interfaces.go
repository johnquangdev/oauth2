package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	uModels "github.com/johnquangdev/oauth2/usecase/models"
)

type GoogleOauth2 interface {
	// NewUser(context.Context, rModels.User) error
	// NewSession(context.Context, rModels.Session) error
	Login(ctx context.Context, code string) (*uModels.TokenJwt, *uModels.User, error)
	Logout(context.Context, uuid.UUID) error
	GetAuthURL() (string, error)
	// AddBackList(uuid.UUID, string, time.Duration) error
}
type GithubOauth2 interface {
	Login(ctx context.Context, code string) (*uModels.TokenJwt, *uModels.User, error)
	GetAuthURL() (string, error)
}
type SystemAuth interface {
	AddBackList(uuid.UUID, string, time.Duration) error
	GetUserByProviderAndProviderId(context.Context, string, string) (*uModels.User, error)
	GetUserById(context.Context, uuid.UUID) (*uModels.User, error)
}
type AuthImpl struct {
	GoogleOauth2 GoogleOauth2
	GithubOauth2 GithubOauth2
	SystemAuth   SystemAuth
	//FacebookOauth2() FacebookOauth2
}

type UseCaseImpl interface {
	Auth() AuthImpl
}
