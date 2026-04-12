package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ObservabilityDataSourceRepository struct {
	db *gorm.DB
}

func NewObservabilityDataSourceRepository(db *gorm.DB) *ObservabilityDataSourceRepository {
	return &ObservabilityDataSourceRepository{db: db}
}

func (r *ObservabilityDataSourceRepository) Create(ctx context.Context, item *domain.ObservabilityDataSource) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ObservabilityDataSourceRepository) GetByID(ctx context.Context, id uint64) (*domain.ObservabilityDataSource, error) {
	var item domain.ObservabilityDataSource
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ObservabilityDataSourceRepository) ListByCluster(ctx context.Context, clusterID uint64) ([]domain.ObservabilityDataSource, error) {
	var items []domain.ObservabilityDataSource
	err := r.db.WithContext(ctx).
		Where("cluster_id = ? OR cluster_id IS NULL", clusterID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *ObservabilityDataSourceRepository) Update(ctx context.Context, item *domain.ObservabilityDataSource) error {
	return r.db.WithContext(ctx).Save(item).Error
}

