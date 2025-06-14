package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email      string    `gorm:"type:text;not null;unique" json:"email"`
	Name       string    `gorm:"type:text" json:"name"`
	Avatar     string    `gorm:"type:text" json:"avatar"`
	Provider   string    `gorm:"type:text;not null" json:"provider"`    // e.g. "google"
	ProviderID string    `gorm:"type:text;not null" json:"provider_id"` // unique user ID from provider
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableUsers() string {
	return "users"
}

type Session struct {
	ID                    string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID                string    `gorm:"type:uuid;not null" json:"user_id"`
	RefreshToken          string    `gorm:"type:text;not null" json:"refresh_token"`
	AccessToken           string    `gorm:"type:text;not null" json:"access_Token"`
	UserAgent             string    `gorm:"type:text" json:"user_agent"`
	IPAddress             string    `gorm:"type:text" json:"ip_address"`
	Status                string    `gorm:"default:active" json:"status"`
	RefreshTokenExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time `gorm:"type:timestamptz;not null" json:"access_token_expires_at"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Session) TableSession() string {
	return "session"
}
