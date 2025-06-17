package usecase

import (
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	google "github.com/johnquangdev/oauth2/service/google"
	uImlp "github.com/johnquangdev/oauth2/usecase/imlp"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ImlpUseCase struct {
	db            *gorm.DB
	oauth2Service *google.ServiceOauthGoogle
	repo          rInterfaces.Repo
	redis         *redis.Client
}

func (u *ImlpUseCase) Auth() (interfaces.Auth, error) {
	return uImlp.NewOAuthUsecase(u.db, u.redis)
}

func NewUseCase(repo rInterfaces.Repo, db *gorm.DB, oauth2Service *google.ServiceOauthGoogle, redis *redis.Client) (interfaces.UseCase, error) {
	return &ImlpUseCase{
		db:            db,
		oauth2Service: oauth2Service,
		repo:          repo,
		redis:         redis,
	}, nil
}
