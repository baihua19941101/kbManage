package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RestorePointListFilter struct {
	WorkspaceIDs []uint64
	PolicyID     uint64
	Result       string
}

type RestorePointRepository struct {
	db *gorm.DB
}

func NewRestorePointRepository(db *gorm.DB) *RestorePointRepository {
	return &RestorePointRepository{db: db}
}

func (r *RestorePointRepository) Create(ctx context.Context, item *domain.RestorePoint) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RestorePointRepository) Update(ctx context.Context, item *domain.RestorePoint) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *RestorePointRepository) GetByID(ctx context.Context, id uint64) (*domain.RestorePoint, error) {
	var item domain.RestorePoint
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RestorePointRepository) List(ctx context.Context, filter RestorePointListFilter) ([]domain.RestorePoint, error) {
	query := r.db.WithContext(ctx).Model(&domain.RestorePoint{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if filter.PolicyID != 0 {
		query = query.Where("policy_id = ?", filter.PolicyID)
	}
	if v := strings.TrimSpace(filter.Result); v != "" {
		query = query.Where("result = ?", v)
	}
	var items []domain.RestorePoint
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
