package usecase

import (
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/usecase/impl"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/redis/go-redis/v9"
)

type UseCase struct {
	redis *redis.Client
	cfg   utils.Config
	repo  rInterfaces.Repo
}

func (u UseCase) Auth() interfaces.AuthImpl {
	google := impl.NewOAuth2Google(u.cfg, u.repo)
	github := impl.NewOAuth2Github(u.cfg, u.redis, u.repo)
	auth := impl.NewSystemAuth(u.repo)
	return interfaces.AuthImpl{
		GoogleOauth2: google,
		GithubOauth2: github,
		SystemAuth:   auth,
	}
}

func NewUseCase(cfg utils.Config, repo rInterfaces.Repo, redis *redis.Client) (interfaces.UseCaseImpl, error) {
	return &UseCase{
		redis: redis,
		repo:  repo,
		cfg:   cfg,
	}, nil
}
