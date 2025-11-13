package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/repository/models"
	"github.com/johnquangdev/oauth2/service/github"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	uModels "github.com/johnquangdev/oauth2/usecase/models"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GithubOAuth2Impl struct {
	repo  rInterfaces.Repo
	cfg   utils.Config
	git   *github.Oauth2GithubService
	redis *redis.Client
}

func NewOAuth2Github(cfg utils.Config, redis *redis.Client, r rInterfaces.Repo) interfaces.GithubOauth2 {
	github, err := github.NewGithubOauth2Service(r, cfg)
	if err != nil {
		return nil
	}
	return &GithubOAuth2Impl{
		repo:  r,
		cfg:   cfg,
		git:   github,
		redis: redis,
	}
}

func (g *GithubOAuth2Impl) GetAuthURL() (string, error) {
	url, err := g.git.GenerateAuthURL()
	if err != nil {
		return "", fmt.Errorf("cannot generate github auth url: %w", err)
	}
	return url.Url, nil
}

func (g *GithubOAuth2Impl) Login(ctx context.Context, code string) (*uModels.TokenJwt, *uModels.User, error) {
	// Exchange code for access token
	githubAccessToken, err := g.git.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	//get user info from github
	userInfoGithub, err := g.git.GetUserInfoGithub(ctx, githubAccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info from github: %w", err)
	}

	//get user email from github
	emailInfoGithub, err := g.git.GetUserEmail(ctx, githubAccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user email from github: %w", err)
	}
	//check user exist in db
	userExist, err := g.repo.Auth().GetUserByProviderAndProviderId(ctx, uModels.ProviderGitHub, userInfoGithub.Provider)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userExist = &models.User{
				Id:         uuid.New(),
				Email:      emailInfoGithub,
				Name:       userInfoGithub.Name,
				Avatar:     userInfoGithub.AvatarURL,
				Status:     uModels.StatusPending,
				Provider:   uModels.ProviderGitHub,
				ProviderId: userInfoGithub.Provider,
			}
			if err := g.repo.Auth().CreateUser(userExist); err != nil {
				return nil, nil, fmt.Errorf("create user error: %w", err)
			}
		} else {
			return nil, nil, err
		}
	}
	//Create JWT token
	accessTokenTimeLife := time.Duration(g.cfg.AccessTokenTimeLife) * time.Minute
	refreshTokenTimeLife := time.Duration(g.cfg.RefreshTokenTimeLife) * time.Hour
	accessToken, claimsAccess := utils.GenerateToken(userExist.Id, userExist.Email, userExist.Name, accessTokenTimeLife, g.cfg.SecretKey)
	refreshToken, claimsRefresh := utils.GenerateToken(userExist.Id, userExist.Email, userExist.Name, refreshTokenTimeLife, g.cfg.SecretKey)
	// Save session to redis
	session := &models.Session{
		Id:                    uuid.New(),
		UserId:                userExist.Id,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: claimsRefresh.ExpiresAt.Time,
	}
	if err := g.repo.Auth().CreateSession(session); err != nil {
		return nil, nil, fmt.Errorf("create session error: %w", err)
	}
	// save accessToken for redis
	err = g.repo.Redis().CreateRecord(userExist.Id, accessToken, time.Until(claimsAccess.ExpiresAt.Time))
	if err != nil {
		return nil, nil, fmt.Errorf("create redis record error: %w", err)
	}

	return &uModels.TokenJwt{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  time.Until(claimsAccess.ExpiresAt.Time),
			RefreshTokenExpiresAt: time.Until(claimsRefresh.ExpiresAt.Time),
		}, &uModels.User{
			Id:         userExist.Id,
			Email:      userExist.Email,
			Name:       userExist.Name,
			Status:     userExist.Status,
			Avatar:     userExist.Avatar,
			Provider:   userExist.Provider,
			ProviderId: userExist.ProviderId,
			CreatedAt:  userExist.CreatedAt,
			UpdatedAt:  userExist.UpdatedAt,
		}, nil
}
