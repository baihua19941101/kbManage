package domain

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID           uint64     `gorm:"primaryKey"`
	Username     string     `gorm:"size:128;uniqueIndex;not null"`
	DisplayName  string     `gorm:"size:128"`
	Email        string     `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string     `gorm:"size:255;not null"`
	Status       UserStatus `gorm:"size:32;not null;default:active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID           uint64    `gorm:"primaryKey"`
	UserID       uint64    `gorm:"index;not null"`
	RefreshToken string    `gorm:"size:255;not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	Revoked      bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
