package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DelegationGrantRepository struct {
	db *gorm.DB
}

func NewDelegationGrantRepository(db *gorm.DB) *DelegationGrantRepository {
	return &DelegationGrantRepository{db: db}
}

func (r *DelegationGrantRepository) Create(ctx context.Context, item *domain.DelegationGrant) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DelegationGrantRepository) Update(ctx context.Context, item *domain.DelegationGrant) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *DelegationGrantRepository) GetByID(ctx context.Context, id uint64) (*domain.DelegationGrant, error) {
	var item domain.DelegationGrant
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DelegationGrantRepository) List(ctx context.Context, delegateRef string) ([]domain.DelegationGrant, error) {
	query := r.db.WithContext(ctx).Model(&domain.DelegationGrant{})
	if v := strings.TrimSpace(delegateRef); v != "" {
		query = query.Where("delegate_ref = ?", v)
	}
	var items []domain.DelegationGrant
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *DelegationGrantRepository) ListByDelegate(ctx context.Context, delegateRef string) ([]domain.DelegationGrant, error) {
	var items []domain.DelegationGrant
	err := r.db.WithContext(ctx).
		Where("delegate_ref = ?", delegateRef).
		Order("id ASC").
		Find(&items).Error
	return items, err
}
