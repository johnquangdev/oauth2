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

func (r repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r repository) CreateUser(user *models.User) error {
	return r.db.Create(&user).Error
}

func (r repository) CreateSesion(session *models.Session) error {
	return r.db.Create(&session).Error
}

func (r repository) BlockedUserByUserID(ctx context.Context, userID uuid.UUID) error {
	var s *models.Session
	result := r.db.WithContext(ctx).Model(&s).
		Where("user_id=?", userID).
		Updates(map[string]interface{}{
			"status":        "blocked",
			"refresh_token": "",
		})
	if result.Error != nil {
		return fmt.Errorf("failed to block session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no active session found with provided token")
	}
	return nil
}

func (r repository) IsUserExists(email string) (bool, error) {
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
