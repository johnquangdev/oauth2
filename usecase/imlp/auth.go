package impl

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/repository"
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/repository/models"
	google "github.com/johnquangdev/oauth2/service/google"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	UModels "github.com/johnquangdev/oauth2/usecase/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	ProviderGoogle   = "google"
	ProviderGitHub   = "github"
	ProviderFacebook = "facebook"
)

type OAuthUsecase struct {
	oauthService *google.ServiceOauthGoogle
	oauth2       *google.Soauth2
	repo         rInterfaces.Repo
}

type GetAuthURLRequest struct {
	CustomState string `json:"custom_state,omitempty"`
}

func NewOAuthUsecase(db *gorm.DB, redis *redis.Client) (interfaces.Auth, error) {
	u, err := google.Newoauth2()
	if err != nil {
		return nil, err
	}
	return &OAuthUsecase{
		oauth2:       u,
		repo:         repository.NewRepository(db, redis),
		oauthService: google.NewGoogleOAuthService(),
	}, err
}

func (uc *OAuthUsecase) GetGoogleAuthURL(state string) (string, error) {
	// Generate auth URL
	response := uc.oauthService.GenerateAuthURL(state)
	return response, nil
}

func (uc *OAuthUsecase) ExchangeCodeForToken(ctx context.Context, req UModels.ExchangeTokenRequest) (UModels.OAuthResponse, error) {
	// Validate input
	if req.Code == "" {
		return UModels.OAuthResponse{}, errors.New("authorization code is required")
	}

	// Exchange code for token
	tokenResponse, err := uc.oauthService.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return UModels.OAuthResponse{}, errors.New("failed to exchange code for token: " + err.Error())
	}

	return UModels.OAuthResponse{
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

func (uc *OAuthUsecase) CreateUser(ctx context.Context, user models.User) error {
	err := uc.repo.Auth().CreateUser(&models.User{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (uc *OAuthUsecase) CreateSesion(ctx context.Context, session models.Session) error {
	err := uc.repo.Auth().CreateSesion(&models.Session{
		ID:                    session.ID,
		UserID:                session.UserID,
		RefreshToken:          session.RefreshToken,
		AccessToken:           session.AccessToken,
		UserAgent:             session.UserAgent,
		IPAddress:             session.IPAddress,
		Status:                session.Status,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (uc *OAuthUsecase) GetUserByEmail(gmail string) (*UModels.User, error) {
	user, err := uc.repo.Auth().GetUserByEmail(gmail)
	if err != nil {
		return nil, err
	}
	return &UModels.User{
		Email:      user.Email,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
	}, nil
}
func (uc *OAuthUsecase) SaveRefreshToken(useID string, token string, expires time.Duration) error {
	uc.repo.Redis().SaveRefreshToken(useID, token, expires)
	return nil
}
func (uc *OAuthUsecase) GetUserInfoGoogle(accessToken string) (UModels.User, error) {
	userInfo, err := uc.oauth2.GetUserInfoGoogle(accessToken)
	if err != nil {
		return UModels.User{}, err
	}
	return UModels.User{
		Id:         uuid.New(),
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		Avatar:     userInfo.Picture,
		Provider:   ProviderGoogle,
		ProviderID: userInfo.Id,
	}, nil
}

func (uc *OAuthUsecase) IsUserExists(gmail string) error {
	exist, err := uc.repo.Auth().IsUserExists(gmail)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("user exists")
	}
	return nil
}

func (uc *OAuthUsecase) Logout(ctx context.Context, userID uuid.UUID) error {
	// call repo to block user
	err := uc.repo.Auth().BlockedUserByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't block user by err: %v", err)
	}
	return nil
}
