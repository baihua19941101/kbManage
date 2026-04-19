package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ExtensionLifecycleRepository struct{ db *gorm.DB }

func NewExtensionLifecycleRepository(db *gorm.DB) *ExtensionLifecycleRepository {
	return &ExtensionLifecycleRepository{db: db}
}

func (r *ExtensionLifecycleRepository) Create(ctx context.Context, item *domain.ExtensionLifecycleRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ExtensionLifecycleRepository) ListByExtensionID(ctx context.Context, extensionID uint64) ([]domain.ExtensionLifecycleRecord, error) {
	var items []domain.ExtensionLifecycleRecord
	err := r.db.WithContext(ctx).Where("extension_package_id = ?", extensionID).Order("id DESC").Find(&items).Error
	return items, err
}
