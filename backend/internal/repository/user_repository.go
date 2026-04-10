package repository

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB

	mu         sync.RWMutex
	nextID     uint64
	byID       map[uint64]domain.User
	byUsername map[string]uint64
	byEmail    map[string]uint64
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:         db,
		nextID:     1,
		byID:       make(map[uint64]domain.User),
		byUsername: make(map[string]uint64),
		byEmail:    make(map[string]uint64),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return gorm.ErrInvalidData
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(user).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byUsername[user.Username]; ok {
		return gorm.ErrDuplicatedKey
	}
	if _, ok := r.byEmail[user.Email]; ok {
		return gorm.ErrDuplicatedKey
	}

	copyUser := *user
	if copyUser.ID == 0 {
		copyUser.ID = r.nextID
		r.nextID++
	}
	now := time.Now()
	if copyUser.CreatedAt.IsZero() {
		copyUser.CreatedAt = now
	}
	copyUser.UpdatedAt = now
	if copyUser.Status == "" {
		copyUser.Status = domain.UserStatusActive
	}

	r.byID[copyUser.ID] = copyUser
	r.byUsername[copyUser.Username] = copyUser.ID
	r.byEmail[copyUser.Email] = copyUser.ID
	*user = copyUser
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	if r.db != nil {
		var user domain.User
		if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyUser := user
	return &copyUser, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, errors.New("username is required")
	}
	if r.db != nil {
		var user domain.User
		if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byUsername[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	user := r.byID[id]
	copyUser := user
	return &copyUser, nil
}
