package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type TenantScopeMappingRepository struct {
	db *gorm.DB
}

func NewTenantScopeMappingRepository(db *gorm.DB) *TenantScopeMappingRepository {
	return &TenantScopeMappingRepository{db: db}
}

func (r *TenantScopeMappingRepository) Create(ctx context.Context, item *domain.TenantScopeMapping) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *TenantScopeMappingRepository) ListByUnitID(ctx context.Context, unitID uint64) ([]domain.TenantScopeMapping, error) {
	var items []domain.TenantScopeMapping
	err := r.db.WithContext(ctx).
		Where("unit_id = ?", unitID).
		Order("id ASC").
		Find(&items).Error
	return items, err
}

func (r *TenantScopeMappingRepository) FindByUnitScope(ctx context.Context, unitID uint64, scopeType, scopeRef string) (*domain.TenantScopeMapping, error) {
	var item domain.TenantScopeMapping
	if err := r.db.WithContext(ctx).
		Where("unit_id = ? AND scope_type = ? AND scope_ref = ?", unitID, scopeType, scopeRef).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
