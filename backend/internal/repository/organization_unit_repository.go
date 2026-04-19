package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type OrganizationUnitListFilter struct {
	UnitType     string
	ParentUnitID uint64
}

type OrganizationUnitRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitRepository(db *gorm.DB) *OrganizationUnitRepository {
	return &OrganizationUnitRepository{db: db}
}

func (r *OrganizationUnitRepository) Create(ctx context.Context, item *domain.OrganizationUnit) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *OrganizationUnitRepository) Update(ctx context.Context, item *domain.OrganizationUnit) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *OrganizationUnitRepository) GetByID(ctx context.Context, id uint64) (*domain.OrganizationUnit, error) {
	var item domain.OrganizationUnit
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrganizationUnitRepository) FindByParentName(ctx context.Context, parentID *uint64, name string) (*domain.OrganizationUnit, error) {
	var item domain.OrganizationUnit
	query := r.db.WithContext(ctx).Where("lower(name) = ?", strings.ToLower(strings.TrimSpace(name)))
	if parentID == nil || *parentID == 0 {
		query = query.Where("parent_unit_id IS NULL")
	} else {
		query = query.Where("parent_unit_id = ?", *parentID)
	}
	if err := query.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrganizationUnitRepository) List(ctx context.Context, filter OrganizationUnitListFilter) ([]domain.OrganizationUnit, error) {
	query := r.db.WithContext(ctx).Model(&domain.OrganizationUnit{})
	if v := strings.TrimSpace(filter.UnitType); v != "" {
		query = query.Where("unit_type = ?", v)
	}
	if filter.ParentUnitID != 0 {
		query = query.Where("parent_unit_id = ?", filter.ParentUnitID)
	}
	var items []domain.OrganizationUnit
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}
