package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) Create(ctx context.Context, cluster *domain.Cluster) error {
	return r.db.WithContext(ctx).Create(cluster).Error
}

func (r *ClusterRepository) GetByID(ctx context.Context, id uint64) (*domain.Cluster, error) {
	var cluster domain.Cluster
	if err := r.db.WithContext(ctx).First(&cluster, id).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (r *ClusterRepository) List(ctx context.Context) ([]domain.Cluster, error) {
	var clusters []domain.Cluster
	err := r.db.WithContext(ctx).Order("id DESC").Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) UpdateStatus(ctx context.Context, id uint64, status domain.ClusterStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Cluster{}).Where("id = ?", id).Update("status", status).Error
}
