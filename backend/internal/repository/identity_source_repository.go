package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type IdentitySourceListFilter struct {
	SourceType string
	Status     string
}

type IdentitySourceRepository struct {
	db *gorm.DB
}

func NewIdentitySourceRepository(db *gorm.DB) *IdentitySourceRepository {
	return &IdentitySourceRepository{db: db}
}

func (r *IdentitySourceRepository) Create(ctx context.Context, item *domain.IdentitySource) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *IdentitySourceRepository) Update(ctx context.Context, item *domain.IdentitySource) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *IdentitySourceRepository) GetByID(ctx context.Context, id uint64) (*domain.IdentitySource, error) {
	var item domain.IdentitySource
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IdentitySourceRepository) FindByName(ctx context.Context, name string) (*domain.IdentitySource, error) {
	var item domain.IdentitySource
	err := r.db.WithContext(ctx).
		Where("lower(name) = ?", strings.ToLower(strings.TrimSpace(name))).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IdentitySourceRepository) FindLocal(ctx context.Context) (*domain.IdentitySource, error) {
	var item domain.IdentitySource
	err := r.db.WithContext(ctx).
		Where("source_type = ?", domain.IdentitySourceTypeLocal).
		Order("id ASC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IdentitySourceRepository) List(ctx context.Context, filter IdentitySourceListFilter) ([]domain.IdentitySource, error) {
	query := r.db.WithContext(ctx).Model(&domain.IdentitySource{})
	if v := strings.TrimSpace(filter.SourceType); v != "" {
		query = query.Where("source_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.IdentitySource
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}
