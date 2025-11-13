package impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	rInterfaces "github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/usecase/models"
)

type AuthImpl struct {
	repo rInterfaces.Repo
}

func NewSystemAuth(repo rInterfaces.Repo) interfaces.SystemAuth {
	return &AuthImpl{
		repo: repo,
	}
}

func (u AuthImpl) GetUserByProviderAndProviderId(ctx context.Context, provider string, providerId string) (*models.User, error) {
	// Call to repository layer to get user by provider ID
	user, err := u.repo.Auth().GetUserByProviderAndProviderId(ctx, provider, providerId)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id:       user.Id,
		Email:    user.Email,
		Provider: user.Provider,
		Status:   user.Status,
	}, nil
}
func (u AuthImpl) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.repo.Auth().GetUserByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Avatar:     user.Avatar,
		Provider:   user.Provider,
		Status:     user.Status,
		ProviderId: user.ProviderId,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}
func (u AuthImpl) AddBackList(uuid.UUID, string, time.Duration) error {
	return nil
}
