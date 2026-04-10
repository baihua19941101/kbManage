package auth

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty   = errors.New("password must not be empty")
	ErrPasswordTooLong = errors.New("password exceeds bcrypt limit (72 bytes)")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters")
)

type PasswordService struct {
	cost int
}

func NewPasswordService(cost int) *PasswordService {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &PasswordService{cost: cost}
}

func (s *PasswordService) Hash(password string) (string, error) {
	if err := validatePassword(password); err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *PasswordService) Verify(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func validatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}
	if len([]byte(password)) > 72 {
		return ErrPasswordTooLong
	}
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordTooWeak
	}
	return nil
}
