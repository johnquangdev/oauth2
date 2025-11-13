package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/repository/models"
	"github.com/johnquangdev/oauth2/service/google"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	uModels "github.com/johnquangdev/oauth2/usecase/models"
	"github.com/johnquangdev/oauth2/utils"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type GoogleOAuth2Impl struct {
	oauthService *google.ServiceOauthGoogle
	repo         rInterfaces.Repo
	cfg          utils.Config
}

type GetAuthURLRequest struct {
	CustomState string `json:"custom_state,omitempty"`
}

func NewOAuth2Google(cfg utils.Config, r rInterfaces.Repo) interfaces.GoogleOauth2 {
	g, err := google.NewGoogleOAuthService(cfg)
	if err != nil {
		return nil
	}
	return &GoogleOAuth2Impl{
		repo:         r,
		oauthService: g,
		cfg:          cfg,
	}
}
func (u *GoogleOAuth2Impl) GetAuthURL() (string, error) {
	url, err := u.oauthService.GenerateAuthURL()
	if err != nil {
		return "", fmt.Errorf("cannot generate google auth url: %w", err)
	}
	return url, nil
}

func (u *GoogleOAuth2Impl) Login(ctx context.Context, code string) (*uModels.TokenJwt, *uModels.User, error) {
	// validate code
	if code == "" {
		return nil, nil, fmt.Errorf("code is required")
	}

	//Exchange code for googleAccessToken
	googleAccessToken, idToken, err := u.oauthService.ChangeCodeToToken(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	//Verify ID Token (xác thực danh tính)
	// idtoken.Validate tự động verify với Google's public keys và kiểm tra audience
	payload, err := idtoken.Validate(ctx, idToken, u.cfg.ClientId_Google)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid id_token: %v", err)
	}

	// Verify issuer
	if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
		return nil, nil, fmt.Errorf("invalid issuer: %s", payload.Issuer)
	}

	// Verify email_verified
	emailVerified, ok := payload.Claims["email_verified"].(bool)
	if !ok || !emailVerified {
		return nil, nil, fmt.Errorf("email not verified")
	}

	// Verify email exists
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, nil, fmt.Errorf("email claim missing")
	}

	//Get userInfoGoogle
	userInfoGoogle, err := u.oauthService.GetUserInfoGoogle(googleAccessToken)
	if err != nil {
		return nil, nil, err
	}

	//Check userExist
	userExist, err := u.repo.Auth().GetUserByProviderAndProviderId(ctx, uModels.ProviderGoogle, userInfoGoogle.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userExist = &models.User{
				Id:         uuid.New(),
				Email:      email,
				Name:       userInfoGoogle.Name,
				Avatar:     userInfoGoogle.Picture,
				Status:     uModels.StatusPending,
				Provider:   uModels.ProviderGoogle,
				ProviderId: payload.Claims["sub"].(string),
			}
			if err := u.repo.Auth().CreateUser(userExist); err != nil {
				return nil, nil, fmt.Errorf("create user error: %w", err)
			}
		} else {
			return nil, nil, err
		}
	}

	// create JWT (access + refresh token)
	accessTokenTimeLife := time.Duration(u.cfg.AccessTokenTimeLife) * time.Minute
	refreshTokenTimeLife := time.Duration(u.cfg.RefreshTokenTimeLife) * time.Hour
	accessToken, claimsAccess := utils.GenerateToken(userExist.Id, userExist.Email, userExist.Name, accessTokenTimeLife, u.cfg.SecretKey)
	refreshToken, claimsRefresh := utils.GenerateToken(userExist.Id, userExist.Email, userExist.Name, refreshTokenTimeLife, u.cfg.SecretKey)

	// create session
	session := &models.Session{
		Id:                    uuid.New(),
		UserId:                userExist.Id,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: claimsRefresh.ExpiresAt.Time,
	}
	if err := u.repo.Auth().CreateSession(session); err != nil {
		return nil, nil, fmt.Errorf("create session error: %w", err)
	}

	// save accessToken for redis
	err = u.repo.Redis().CreateRecord(userExist.Id, accessToken, time.Until(claimsAccess.ExpiresAt.Time))
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
			ProviderId: userExist.ProviderId,
			Provider:   userExist.Provider,
			CreatedAt:  userExist.CreatedAt,
			UpdatedAt:  userExist.UpdatedAt,
		}, nil
}

func (u *GoogleOAuth2Impl) Logout(ctx context.Context, userID uuid.UUID) error {
	// call repo to block user
	err := u.repo.Auth().BlockedUserByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("can't block user by err: %v", err)
	}
	return nil
}
