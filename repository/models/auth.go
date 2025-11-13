package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string
type User struct {
	Id         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email      string    `gorm:"type:text;not null;unique" json:"email"`
	Name       string    `gorm:"type:text" json:"name"`
	Avatar     string    `gorm:"type:text" json:"avatar"`
	Status     string    `gorm:"type:text; check:status IN ('active','blocked','banned')"`
	Provider   string    `gorm:"type:text;not null" json:"provider"`
	ProviderId string    `gorm:"type:text;not null" json:"provider_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableUsers() string {
	return "users"
}

type Session struct {
	Id                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserId                uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RefreshToken          string    `gorm:"type:text;not null" json:"refresh_token"`
	UserAgent             string    `gorm:"type:text" json:"user_agent"`
	IPAddress             string    `gorm:"type:text" json:"ip_address"`
	RefreshTokenExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"refresh_token_expires_at"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Session) TableSession() string {
	return "session"
}
