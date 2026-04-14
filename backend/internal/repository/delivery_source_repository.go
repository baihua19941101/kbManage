package repository

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliverySourceStore interface {
	Create(ctx context.Context, item *domain.DeliverySource) error
	GetByID(ctx context.Context, id uint64) (*domain.DeliverySource, error)
	ListByScope(
		ctx context.Context,
		workspaceID *uint64,
		projectID *uint64,
		sourceType domain.DeliverySourceType,
		status domain.DeliverySourceStatus,
	) ([]domain.DeliverySource, error)
	Update(ctx context.Context, item *domain.DeliverySource) error
}

type DeliverySourceRepository struct {
	db *gorm.DB
}

func NewDeliverySourceRepository(db *gorm.DB) *DeliverySourceRepository {
	return &DeliverySourceRepository{db: db}
}

func (r *DeliverySourceRepository) Create(ctx context.Context, item *domain.DeliverySource) error {
	if item == nil {
		return errors.New("delivery source is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery source repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliverySourceRepository) GetByID(ctx context.Context, id uint64) (*domain.DeliverySource, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery source repository is not configured")
	}
	var item domain.DeliverySource
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DeliverySourceRepository) ListByScope(
	ctx context.Context,
	workspaceID *uint64,
	projectID *uint64,
	sourceType domain.DeliverySourceType,
	status domain.DeliverySourceStatus,
) ([]domain.DeliverySource, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("delivery source repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.DeliverySource{})
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	}
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}
	if strings.TrimSpace(string(sourceType)) != "" {
		query = query.Where("source_type = ?", sourceType)
	}
	if strings.TrimSpace(string(status)) != "" {
		query = query.Where("status = ?", status)
	}
	items := make([]domain.DeliverySource, 0)
	if err := query.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DeliverySourceRepository) Update(ctx context.Context, item *domain.DeliverySource) error {
	if item == nil || item.ID == 0 {
		return errors.New("delivery source id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("delivery source repository is not configured")
	}
	return r.db.WithContext(ctx).Model(&domain.DeliverySource{}).Where("id = ?", item.ID).Updates(item).Error
}
