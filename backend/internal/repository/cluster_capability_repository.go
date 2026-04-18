package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterCapabilityRepository struct {
	db *gorm.DB
}

func NewClusterCapabilityRepository(db *gorm.DB) *ClusterCapabilityRepository {
	return &ClusterCapabilityRepository{db: db}
}

func (r *ClusterCapabilityRepository) ReplaceForOwner(ctx context.Context, ownerType domain.CapabilityOwnerType, ownerRef string, items []domain.CapabilityMatrixEntry) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Where("owner_type = ? AND owner_ref = ?", ownerType, ownerRef).Delete(&domain.CapabilityMatrixEntry{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for i := range items {
		items[i].OwnerType = ownerType
		items[i].OwnerRef = ownerRef
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *ClusterCapabilityRepository) ListByOwner(ctx context.Context, ownerType domain.CapabilityOwnerType, ownerRef string) ([]domain.CapabilityMatrixEntry, error) {
	var items []domain.CapabilityMatrixEntry
	err := r.db.WithContext(ctx).
		Where("owner_type = ? AND owner_ref = ?", ownerType, ownerRef).
		Order("capability_domain ASC").
		Find(&items).Error
	return items, err
}
