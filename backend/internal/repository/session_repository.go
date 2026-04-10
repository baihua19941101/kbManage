package repository

import (
	"context"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB

	mu        sync.RWMutex
	nextID    uint64
	byID      map[uint64]domain.Session
	byTokenID map[string]uint64
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db:        db,
		nextID:    1,
		byID:      make(map[uint64]domain.Session),
		byTokenID: make(map[string]uint64),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if session == nil {
		return gorm.ErrInvalidData
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(session).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byTokenID[session.RefreshToken]; ok {
		return gorm.ErrDuplicatedKey
	}

	copySession := *session
	if copySession.ID == 0 {
		copySession.ID = r.nextID
		r.nextID++
	}
	now := time.Now()
	if copySession.CreatedAt.IsZero() {
		copySession.CreatedAt = now
	}
	copySession.UpdatedAt = now

	r.byID[copySession.ID] = copySession
	r.byTokenID[copySession.RefreshToken] = copySession.ID
	*session = copySession
	return nil
}

func (r *SessionRepository) GetActiveByToken(ctx context.Context, token string) (*domain.Session, error) {
	if r.db != nil {
		var s domain.Session
		err := r.db.WithContext(ctx).
			Where("refresh_token = ? AND revoked = 0 AND expires_at > ?", token, time.Now()).
			First(&s).Error
		if err != nil {
			return nil, err
		}
		return &s, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byTokenID[token]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	s, ok := r.byID[id]
	if !ok || s.Revoked || !s.ExpiresAt.After(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}
	copySession := s
	return &copySession, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, id uint64) error {
	if r.db != nil {
		return r.db.WithContext(ctx).Model(&domain.Session{}).Where("id = ?", id).Update("revoked", true).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	item.Revoked = true
	item.UpdatedAt = time.Now()
	r.byID[id] = item
	return nil
}

func (r *SessionRepository) Rotate(ctx context.Context, revokeSessionID uint64, newSession *domain.Session) error {
	if newSession == nil {
		return gorm.ErrInvalidData
	}
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		old, ok := r.byID[revokeSessionID]
		if !ok || old.Revoked {
			return gorm.ErrRecordNotFound
		}
		if _, exists := r.byTokenID[newSession.RefreshToken]; exists {
			return gorm.ErrDuplicatedKey
		}

		old.Revoked = true
		old.UpdatedAt = time.Now()
		r.byID[old.ID] = old

		next := *newSession
		if next.ID == 0 {
			next.ID = r.nextID
			r.nextID++
		}
		now := time.Now()
		if next.CreatedAt.IsZero() {
			next.CreatedAt = now
		}
		next.UpdatedAt = now

		r.byID[next.ID] = next
		r.byTokenID[next.RefreshToken] = next.ID
		*newSession = next
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.Session{}).
			Where("id = ? AND revoked = 0", revokeSessionID).
			Update("revoked", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Create(newSession).Error
	})
}
