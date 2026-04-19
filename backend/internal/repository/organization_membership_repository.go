package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type OrganizationMembershipRepository struct {
	db *gorm.DB
}

func NewOrganizationMembershipRepository(db *gorm.DB) *OrganizationMembershipRepository {
	return &OrganizationMembershipRepository{db: db}
}

func (r *OrganizationMembershipRepository) Create(ctx context.Context, item *domain.OrganizationMembership) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *OrganizationMembershipRepository) Update(ctx context.Context, item *domain.OrganizationMembership) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *OrganizationMembershipRepository) ListByUnitID(ctx context.Context, unitID uint64) ([]domain.OrganizationMembership, error) {
	var items []domain.OrganizationMembership
	err := r.db.WithContext(ctx).
		Where("unit_id = ?", unitID).
		Order("id ASC").
		Find(&items).Error
	return items, err
}

func (r *OrganizationMembershipRepository) ListByMemberRef(ctx context.Context, memberType, memberRef string) ([]domain.OrganizationMembership, error) {
	var items []domain.OrganizationMembership
	err := r.db.WithContext(ctx).
		Where("member_type = ? AND member_ref = ?", memberType, memberRef).
		Order("id ASC").
		Find(&items).Error
	return items, err
}
