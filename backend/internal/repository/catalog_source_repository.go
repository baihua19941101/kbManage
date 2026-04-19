package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type CatalogSourceListFilter struct {
	SourceType string
	Status     string
	Keyword    string
}

type CatalogSourceRepository struct{ db *gorm.DB }

func NewCatalogSourceRepository(db *gorm.DB) *CatalogSourceRepository {
	return &CatalogSourceRepository{db: db}
}

func (r *CatalogSourceRepository) Create(ctx context.Context, item *domain.CatalogSource) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CatalogSourceRepository) Update(ctx context.Context, item *domain.CatalogSource) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CatalogSourceRepository) FindByName(ctx context.Context, name string) (*domain.CatalogSource, error) {
	var item domain.CatalogSource
	if err := r.db.WithContext(ctx).Where("name = ?", strings.TrimSpace(name)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CatalogSourceRepository) GetByID(ctx context.Context, id uint64) (*domain.CatalogSource, error) {
	var item domain.CatalogSource
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CatalogSourceRepository) List(ctx context.Context, filter CatalogSourceListFilter) ([]domain.CatalogSource, error) {
	query := r.db.WithContext(ctx).Model(&domain.CatalogSource{})
	if v := strings.TrimSpace(filter.SourceType); v != "" {
		query = query.Where("source_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.Keyword); v != "" {
		like := "%" + v + "%"
		query = query.Where("name LIKE ? OR endpoint_ref LIKE ?", like, like)
	}
	var items []domain.CatalogSource
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}
