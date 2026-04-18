package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type NodePoolRepository struct {
	db *gorm.DB
}

func NewNodePoolRepository(db *gorm.DB) *NodePoolRepository {
	return &NodePoolRepository{db: db}
}

func (r *NodePoolRepository) Create(ctx context.Context, item *domain.NodePoolProfile) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *NodePoolRepository) Update(ctx context.Context, item *domain.NodePoolProfile) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *NodePoolRepository) GetByID(ctx context.Context, id uint64) (*domain.NodePoolProfile, error) {
	var item domain.NodePoolProfile
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *NodePoolRepository) ListByClusterID(ctx context.Context, clusterID uint64) ([]domain.NodePoolProfile, error) {
	var items []domain.NodePoolProfile
	err := r.db.WithContext(ctx).Where("cluster_id = ?", clusterID).Order("id ASC").Find(&items).Error
	return items, err
}
