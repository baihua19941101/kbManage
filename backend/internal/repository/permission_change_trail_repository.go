package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PermissionChangeTrailRepository struct{ db *gorm.DB }

func NewPermissionChangeTrailRepository(db *gorm.DB) *PermissionChangeTrailRepository {
	return &PermissionChangeTrailRepository{db: db}
}

func (r *PermissionChangeTrailRepository) Create(ctx context.Context, item *domain.PermissionChangeTrail) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PermissionChangeTrailRepository) List(ctx context.Context) ([]domain.PermissionChangeTrail, error) {
	var items []domain.PermissionChangeTrail
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
