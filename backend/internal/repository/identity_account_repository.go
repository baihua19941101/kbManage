package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type IdentityAccountRepository struct {
	db *gorm.DB
}

func NewIdentityAccountRepository(db *gorm.DB) *IdentityAccountRepository {
	return &IdentityAccountRepository{db: db}
}

func (r *IdentityAccountRepository) Create(ctx context.Context, item *domain.IdentityAccount) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *IdentityAccountRepository) Update(ctx context.Context, item *domain.IdentityAccount) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *IdentityAccountRepository) FindBySourceExternalRef(ctx context.Context, sourceID uint64, externalRef string) (*domain.IdentityAccount, error) {
	var item domain.IdentityAccount
	if err := r.db.WithContext(ctx).
		Where("identity_source_id = ? AND external_ref = ?", sourceID, externalRef).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IdentityAccountRepository) ListBySourceID(ctx context.Context, sourceID uint64) ([]domain.IdentityAccount, error) {
	var items []domain.IdentityAccount
	err := r.db.WithContext(ctx).
		Where("identity_source_id = ?", sourceID).
		Order("id ASC").
		Find(&items).Error
	return items, err
}

func (r *IdentityAccountRepository) ListByUserID(ctx context.Context, userID uint64) ([]domain.IdentityAccount, error) {
	var items []domain.IdentityAccount
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&items).Error
	return items, err
}
