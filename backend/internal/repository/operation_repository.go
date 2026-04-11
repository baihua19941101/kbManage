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
	fromStatuses := []domain.OperationStatus{}
	switch status {
	case domain.OperationStatusRunning:
		fromStatuses = []domain.OperationStatus{domain.OperationStatusPending}
	case domain.OperationStatusSucceeded, domain.OperationStatusFailed:
		fromStatuses = []domain.OperationStatus{domain.OperationStatusRunning}
	default:
		fromStatuses = []domain.OperationStatus{status}
	}
	failureReason := ""
	if status == domain.OperationStatusFailed {
		failureReason = resultMessage
	}
	_, _, err := r.TransitionStatus(ctx, id, fromStatuses, status, strings.TrimSpace(resultMessage), resultMessage, failureReason)
	return err
}

func (r *OperationRepository) UpdateProgress(ctx context.Context, id uint64, progressMessage string) error {
	progressMessage = strings.TrimSpace(progressMessage)
	if progressMessage == "" {
		return nil
	}

	if r.db != nil {
		res := r.db.WithContext(ctx).Model(&domain.OperationRequest{}).Where("id = ?", id).Updates(map[string]any{
			"progress_message": progressMessage,
			"updated_at":       time.Now(),
		})
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
	item.ProgressMessage = progressMessage
	item.UpdatedAt = time.Now()
	r.byID[id] = item
	return nil
}

func (r *OperationRepository) TransitionStatus(
	ctx context.Context,
	id uint64,
	fromStatuses []domain.OperationStatus,
	nextStatus domain.OperationStatus,
	progressMessage string,
	resultMessage string,
	failureReason string,
) (*domain.OperationRequest, bool, error) {
	progressMessage = strings.TrimSpace(progressMessage)
	resultMessage = strings.TrimSpace(resultMessage)
	failureReason = strings.TrimSpace(failureReason)

	if r.db != nil {
		var item domain.OperationRequest
		if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
			return nil, false, err
		}

		if !statusAllowed(item.Status, fromStatuses) || !item.Status.CanTransitTo(nextStatus) {
			return &item, false, nil
		}
		if item.Status == nextStatus {
			return &item, false, nil
		}

		now := time.Now()
		updates := map[string]any{
			"status":     nextStatus,
			"updated_at": now,
		}
		if progressMessage != "" {
			updates["progress_message"] = progressMessage
		}
		if resultMessage != "" {
			updates["result_message"] = resultMessage
		}
		if nextStatus == domain.OperationStatusFailed {
			updates["failure_reason"] = failureReason
		}
		if nextStatus == domain.OperationStatusSucceeded {
			updates["failure_reason"] = ""
		}
		if nextStatus.IsTerminal() {
			updates["completed_at"] = now
		}

		res := r.db.WithContext(ctx).
			Model(&domain.OperationRequest{}).
			Where("id = ? AND status = ?", id, item.Status).
			Updates(updates)
		if res.Error != nil {
			return nil, false, res.Error
		}
		if res.RowsAffected == 0 {
			latest, err := r.GetByID(ctx, id)
			if err != nil {
				return nil, false, err
			}
			return latest, false, nil
		}

		item.Status = nextStatus
		if progressMessage != "" {
			item.ProgressMessage = progressMessage
		}
		if resultMessage != "" {
			item.ResultMessage = resultMessage
		}
		if nextStatus == domain.OperationStatusFailed {
			item.FailureReason = failureReason
		}
		if nextStatus == domain.OperationStatusSucceeded {
			item.FailureReason = ""
		}
		item.UpdatedAt = now
		if nextStatus.IsTerminal() {
			completedAt := now
			item.CompletedAt = &completedAt
		}
		return &item, true, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[id]
	if !ok {
		return nil, false, gorm.ErrRecordNotFound
	}
	if !statusAllowed(item.Status, fromStatuses) || !item.Status.CanTransitTo(nextStatus) {
		copyItem := item
		return &copyItem, false, nil
	}
	if item.Status == nextStatus {
		copyItem := item
		return &copyItem, false, nil
	}

	now := time.Now()
	item.Status = nextStatus
	if progressMessage != "" {
		item.ProgressMessage = progressMessage
	}
	if resultMessage != "" {
		item.ResultMessage = resultMessage
	}
	if nextStatus == domain.OperationStatusFailed {
		item.FailureReason = failureReason
	}
	if nextStatus == domain.OperationStatusSucceeded {
		item.FailureReason = ""
	}
	item.UpdatedAt = now
	if nextStatus.IsTerminal() {
		completedAt := now
		item.CompletedAt = &completedAt
	}
	r.byID[id] = item

	copyItem := item
	return &copyItem, true, nil
}

func statusAllowed(current domain.OperationStatus, allowed []domain.OperationStatus) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, candidate := range allowed {
		if current == candidate {
			return true
		}
	}
	return false
}
