package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/repository/models"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) interfaces.Auth {
	return &repository{
		db: db,
	}
}

func (r repository) GetUserByUserId(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where(&models.User{Id: id}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) GetUserByProviderAndProviderId(ctx context.Context, provider string, providerId string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where(&models.User{Provider: provider, ProviderId: providerId}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) CreateUser(user *models.User) error {
	return r.db.Create(&user).Error
}

func (r repository) CreateSession(session *models.Session) error {
	return r.db.Create(&session).Error
}

func (r repository) BlockedUserByUserID(ctx context.Context, userID uuid.UUID) error {
	var s *models.User
	result := r.db.WithContext(ctx).Model(&s).
		Where("id=?", userID).
		Updates(map[string]interface{}{
			"status": "blocked",
		})
	if result.Error != nil {
		return fmt.Errorf("failed to blocked user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user does not exist")
	}
	return nil
}

func (r repository) UserExists(email string) (bool, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}
