package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type TerminalSessionRepository struct {
	db *gorm.DB

	mu          sync.RWMutex
	nextID      uint64
	byID        map[uint64]domain.TerminalSession
	bySessionID map[string]uint64
}

func NewTerminalSessionRepository(db *gorm.DB) *TerminalSessionRepository {
	return &TerminalSessionRepository{
		db:          db,
		nextID:      1,
		byID:        make(map[uint64]domain.TerminalSession),
		bySessionID: make(map[string]uint64),
	}
}

func (r *TerminalSessionRepository) Create(ctx context.Context, item *domain.TerminalSession) error {
	if item == nil {
		return errors.New("terminal session is required")
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(item).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bySessionID[item.SessionKey]; exists {
		return gorm.ErrDuplicatedKey
	}

	copyItem := *item
	if copyItem.ID == 0 {
		copyItem.ID = r.nextID
		r.nextID++
	}
	now := time.Now()
	if copyItem.CreatedAt.IsZero() {
		copyItem.CreatedAt = now
	}
	copyItem.UpdatedAt = now
	r.byID[copyItem.ID] = copyItem
	r.bySessionID[copyItem.SessionKey] = copyItem.ID
	*item = copyItem
	return nil
}

func (r *TerminalSessionRepository) GetByID(ctx context.Context, id uint64) (*domain.TerminalSession, error) {
	if r.db != nil {
		var item domain.TerminalSession
		if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copyItem := item
	return &copyItem, nil
}

func (r *TerminalSessionRepository) UpdateStatus(ctx context.Context, id uint64, status domain.TerminalSessionStatus, reason string) error {
	if r.db != nil {
		updates := map[string]any{
			"status":       status,
			"close_reason": reason,
			"updated_at":   time.Now(),
		}
		if status == domain.TerminalSessionStatusClosed || status == domain.TerminalSessionStatusExpired || status == domain.TerminalSessionStatusFailed {
			now := time.Now()
			updates["ended_at"] = now
		}
		res := r.db.WithContext(ctx).Model(&domain.TerminalSession{}).Where("id = ?", id).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.byID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	item.Status = status
	item.CloseReason = reason
	item.UpdatedAt = time.Now()
	if status == domain.TerminalSessionStatusClosed || status == domain.TerminalSessionStatusExpired || status == domain.TerminalSessionStatusFailed {
		now := time.Now()
		item.EndedAt = &now
	}
	r.byID[id] = item
	return nil
}
