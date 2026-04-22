package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliveryArtifactRepository struct{ db *gorm.DB }

func NewDeliveryArtifactRepository(db *gorm.DB) *DeliveryArtifactRepository {
	return &DeliveryArtifactRepository{db: db}
}

func (r *DeliveryArtifactRepository) Create(ctx context.Context, item *domain.DeliveryArtifact) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliveryArtifactRepository) List(ctx context.Context) ([]domain.DeliveryArtifact, error) {
	var items []domain.DeliveryArtifact
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
