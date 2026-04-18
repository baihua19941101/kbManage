package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterDriverRepository struct {
	db *gorm.DB
}

func NewClusterDriverRepository(db *gorm.DB) *ClusterDriverRepository {
	return &ClusterDriverRepository{db: db}
}

func (r *ClusterDriverRepository) Create(ctx context.Context, item *domain.ClusterDriverVersion) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ClusterDriverRepository) Update(ctx context.Context, item *domain.ClusterDriverVersion) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ClusterDriverRepository) GetByID(ctx context.Context, id uint64) (*domain.ClusterDriverVersion, error) {
	var item domain.ClusterDriverVersion
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterDriverRepository) List(ctx context.Context, providerType string) ([]domain.ClusterDriverVersion, error) {
	query := r.db.WithContext(ctx).Model(&domain.ClusterDriverVersion{})
	if v := strings.TrimSpace(providerType); v != "" {
		query = query.Where("provider_type = ?", v)
	}
	var items []domain.ClusterDriverVersion
	err := query.Order("driver_key ASC, id DESC").Find(&items).Error
	return items, err
}

func (r *ClusterDriverRepository) FindByKeyVersion(ctx context.Context, key, version string) (*domain.ClusterDriverVersion, error) {
	var item domain.ClusterDriverVersion
	err := r.db.WithContext(ctx).Where("driver_key = ? AND version = ?", strings.TrimSpace(key), strings.TrimSpace(version)).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
