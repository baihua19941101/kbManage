package testutil

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"

	"gorm.io/gorm"
)

type SeedUserInput struct {
	ID          uint64
	Username    string
	Password    string
	DisplayName string
	Email       string
	Status      domain.UserStatus
}

type SeededUser struct {
	User     domain.User
	Password string
}

func SeedUser(t *testing.T, db *gorm.DB, in SeedUserInput) SeededUser {
	t.Helper()

	if db == nil {
		t.Fatal("seed user requires non-nil db")
	}

	username := strings.TrimSpace(in.Username)
	if username == "" {
		t.Fatal("seed user requires username")
	}

	password := in.Password
	if strings.TrimSpace(password) == "" {
		t.Fatal("seed user requires password")
	}

	displayName := strings.TrimSpace(in.DisplayName)
	if displayName == "" {
		displayName = username
	}

	email := strings.TrimSpace(in.Email)
	if email == "" {
		email = fmt.Sprintf("%s@example.test", username)
	}

	passwordSvc := authSvc.NewPasswordService(0)
	passwordHash, err := passwordSvc.Hash(password)
	if err != nil {
		t.Fatalf("hash user password failed: %v", err)
	}

	status := in.Status
	if status == "" {
		status = domain.UserStatusActive
	}

	user := domain.User{
		ID:           in.ID,
		Username:     username,
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: passwordHash,
		Status:       status,
	}

	if err := db.WithContext(context.Background()).Create(&user).Error; err != nil {
		t.Fatalf("seed user failed: %v", err)
	}

	return SeededUser{
		User:     user,
		Password: password,
	}
}

func IssueAccessToken(t *testing.T, cfg repository.Config, userID uint64) string {
	t.Helper()

	tokenSvc := authSvc.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}
