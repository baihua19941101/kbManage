package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliveryOperationStore interface {
	Create(ctx context.Context, item *domain.DeliveryOperation) error
	GetByID(ctx context.Context, id uint64) (*domain.DeliveryOperation, error)
	GetByRequestID(ctx context.Context, requestID string) (*domain.DeliveryOperation, error)
	UpdateStatus(ctx context.Context, id uint64, status domain.DeliveryOperationStatus, progress int, summary, reason string) error
	UpdatePayload(ctx context.Context, id uint64, payloadJSON string) error
}

type DeliveryOperationRepository struct {
	db *gorm.DB
}

func NewDeliveryOperationRepository(db *gorm.DB) *DeliveryOperationRepository {
	return &DeliveryOperationRepository{db: db}
}

func (r *DeliveryOperationRepository) Create(ctx context.Context, item *domain.DeliveryOperation) error {
	if item == nil {
		return errors.New("delivery operation is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery operation repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliveryOperationRepository) GetByID(ctx context.Context, id uint64) (*domain.DeliveryOperation, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery operation repository is not configured")
	}
	var item domain.DeliveryOperation
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DeliveryOperationRepository) GetByRequestID(ctx context.Context, requestID string) (*domain.DeliveryOperation, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery operation repository is not configured")
	}
	var item domain.DeliveryOperation
	if err := r.db.WithContext(ctx).Where("request_id = ?", strings.TrimSpace(requestID)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DeliveryOperationRepository) UpdateStatus(
	ctx context.Context,
	id uint64,
	status domain.DeliveryOperationStatus,
	progress int,
	summary string,
	reason string,
) error {
	if id == 0 {
		return errors.New("delivery operation id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery operation repository is not configured")
	}
	now := time.Now()
	updates := map[string]any{
		"status":           status,
		"progress_percent": progress,
		"result_summary":   strings.TrimSpace(summary),
		"failure_reason":   strings.TrimSpace(reason),
		"updated_at":       now,
	}
	if status == domain.DeliveryOperationStatusRunning {
		updates["started_at"] = now
	}
	if status == domain.DeliveryOperationStatusSucceeded ||
		status == domain.DeliveryOperationStatusPartiallySucceeded ||
		status == domain.DeliveryOperationStatusFailed ||
		status == domain.DeliveryOperationStatusCanceled {
		updates["completed_at"] = now
	}
	res := r.db.WithContext(ctx).Model(&domain.DeliveryOperation{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DeliveryOperationRepository) UpdatePayload(ctx context.Context, id uint64, payloadJSON string) error {
	if id == 0 {
		return errors.New("delivery operation id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery operation repository is not configured")
	}
	res := r.db.WithContext(ctx).
		Model(&domain.DeliveryOperation{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"payload_json": strings.TrimSpace(payloadJSON),
			"updated_at":   time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
