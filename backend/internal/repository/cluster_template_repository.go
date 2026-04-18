package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterTemplateRepository struct {
	db *gorm.DB
}

func NewClusterTemplateRepository(db *gorm.DB) *ClusterTemplateRepository {
	return &ClusterTemplateRepository{db: db}
}

func (r *ClusterTemplateRepository) Create(ctx context.Context, item *domain.ClusterTemplate) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ClusterTemplateRepository) Update(ctx context.Context, item *domain.ClusterTemplate) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ClusterTemplateRepository) GetByID(ctx context.Context, id uint64) (*domain.ClusterTemplate, error) {
	var item domain.ClusterTemplate
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterTemplateRepository) List(ctx context.Context, driverKey, infrastructureType string) ([]domain.ClusterTemplate, error) {
	query := r.db.WithContext(ctx).Model(&domain.ClusterTemplate{})
	if v := strings.TrimSpace(driverKey); v != "" {
		query = query.Where("driver_key = ?", v)
	}
	if v := strings.TrimSpace(infrastructureType); v != "" {
		query = query.Where("infrastructure_type = ?", v)
	}
	var items []domain.ClusterTemplate
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
