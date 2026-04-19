package repository

import (
	"context"
	"strconv"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type CompatibilityStatementRepository struct{ db *gorm.DB }

func NewCompatibilityStatementRepository(db *gorm.DB) *CompatibilityStatementRepository {
	return &CompatibilityStatementRepository{db: db}
}

func (r *CompatibilityStatementRepository) Create(ctx context.Context, item *domain.CompatibilityStatement) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CompatibilityStatementRepository) Update(ctx context.Context, item *domain.CompatibilityStatement) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CompatibilityStatementRepository) FindLatestByExtensionID(ctx context.Context, extensionID uint64) (*domain.CompatibilityStatement, error) {
	var item domain.CompatibilityStatement
	if err := r.db.WithContext(ctx).Where("owner_type = ? AND owner_ref LIKE ?", domain.CompatibilityOwnerExtension, strconv.FormatUint(extensionID, 10)+":%").Order("id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CompatibilityStatementRepository) ListByOwner(ctx context.Context, ownerType domain.CompatibilityOwnerType, ownerRef string) ([]domain.CompatibilityStatement, error) {
	var items []domain.CompatibilityStatement
	err := r.db.WithContext(ctx).Where("owner_type = ? AND owner_ref = ?", ownerType, ownerRef).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *CompatibilityStatementRepository) ReplaceForOwner(ctx context.Context, ownerType domain.CompatibilityOwnerType, ownerRef string, items []domain.CompatibilityStatement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("owner_type = ? AND owner_ref = ?", ownerType, ownerRef).Delete(&domain.CompatibilityStatement{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OwnerType = ownerType
			items[i].OwnerRef = ownerRef
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
