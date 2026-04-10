package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

// OperationRepository persists and queries operation requests.
type OperationRepository struct {
	db *gorm.DB

	mu          sync.RWMutex
	nextID      uint64
	byID        map[uint64]domain.OperationRequest
	byRequestID map[string]uint64
}

func NewOperationRepository(db *gorm.DB) *OperationRepository {
	return &OperationRepository{
		db:          db,
		nextID:      1,
		byID:        make(map[uint64]domain.OperationRequest),
		byRequestID: make(map[string]uint64),
	}
}

func (r *OperationRepository) Create(ctx context.Context, item *domain.OperationRequest) error {
	if item == nil {
		return errors.New("operation request is required")
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(item).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byRequestID[item.RequestID]; ok {
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
	r.byRequestID[copyItem.RequestID] = copyItem.ID
	*item = copyItem
	return nil
}

func (r *OperationRepository) GetByID(ctx context.Context, id uint64) (*domain.OperationRequest, error) {
	_ = ctx
	if r.db != nil {
		var item domain.OperationRequest
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

func (r *OperationRepository) GetByRequestID(ctx context.Context, requestID string) (*domain.OperationRequest, error) {
	_ = ctx
	if r.db != nil {
		var item domain.OperationRequest
		if err := r.db.WithContext(ctx).Where("request_id = ?", requestID).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byRequestID[requestID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	item := r.byID[id]
	copyItem := item
	return &copyItem, nil
}

func (r *OperationRepository) UpdateStatus(ctx context.Context, id uint64, status domain.OperationStatus, resultMessage string) error {
	if r.db != nil {
		return r.db.WithContext(ctx).Model(&domain.OperationRequest{}).Where("id = ?", id).Updates(map[string]any{
			"status":         status,
			"result_message": resultMessage,
			"updated_at":     time.Now(),
		}).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	item.Status = status
	item.ResultMessage = resultMessage
	item.UpdatedAt = time.Now()
	r.byID[id] = item
	return nil
}
