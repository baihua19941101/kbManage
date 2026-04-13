package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type WorkloadActionRepository struct {
	db *gorm.DB

	mu          sync.RWMutex
	next        uint64
	byID        map[uint64]domain.WorkloadActionRequest
	byRequestID map[string]uint64
}

func NewWorkloadActionRepository(db *gorm.DB) *WorkloadActionRepository {
	return &WorkloadActionRepository{
		db:          db,
		next:        1,
		byID:        make(map[uint64]domain.WorkloadActionRequest),
		byRequestID: make(map[string]uint64),
	}
}

func (r *WorkloadActionRepository) Create(ctx context.Context, item *domain.WorkloadActionRequest) error {
	if item == nil {
		return errors.New("workload action request is required")
	}
	if r.db != nil {
		return r.db.WithContext(ctx).Create(item).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	copyItem := *item
	if copyItem.ID == 0 {
		copyItem.ID = r.next
		r.next++
	}
	now := time.Now()
	if copyItem.CreatedAt.IsZero() {
		copyItem.CreatedAt = now
	}
	copyItem.UpdatedAt = now
	r.byID[copyItem.ID] = copyItem
	if copyItem.RequestID != "" {
		r.byRequestID[copyItem.RequestID] = copyItem.ID
	}
	*item = copyItem
	return nil
}

func (r *WorkloadActionRepository) GetByID(ctx context.Context, id uint64) (*domain.WorkloadActionRequest, error) {
	if r.db != nil {
		var item domain.WorkloadActionRequest
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

func (r *WorkloadActionRepository) GetByRequestID(ctx context.Context, requestID string) (*domain.WorkloadActionRequest, error) {
	if r.db != nil {
		var item domain.WorkloadActionRequest
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

func (r *WorkloadActionRepository) UpdateExecutionResult(
	ctx context.Context,
	id uint64,
	status domain.OperationStatus,
	progressMessage string,
	resultMessage string,
	failureReason string,
) error {
	now := time.Now()
	if r.db != nil {
		updates := map[string]any{
			"status":           status,
			"progress_message": progressMessage,
			"result_message":   resultMessage,
			"failure_reason":   failureReason,
			"updated_at":       now,
		}
		if status == domain.OperationStatusRunning {
			updates["started_at"] = now
		}
		if status == domain.OperationStatusSucceeded || status == domain.OperationStatusFailed {
			updates["completed_at"] = now
		}
		res := r.db.WithContext(ctx).Model(&domain.WorkloadActionRequest{}).Where("id = ?", id).Updates(updates)
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
	item.ProgressMessage = progressMessage
	item.ResultMessage = resultMessage
	item.FailureReason = failureReason
	if status == domain.OperationStatusRunning {
		item.StartedAt = &now
	}
	if status == domain.OperationStatusSucceeded || status == domain.OperationStatusFailed {
		item.CompletedAt = &now
	}
	item.UpdatedAt = now
	r.byID[id] = item
	return nil
}
