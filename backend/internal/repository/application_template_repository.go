package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ApplicationTemplateListFilter struct {
	CatalogSourceID uint64
	Category        string
	Status          string
	Keyword         string
}

type ApplicationTemplateRepository struct{ db *gorm.DB }

func NewApplicationTemplateRepository(db *gorm.DB) *ApplicationTemplateRepository {
	return &ApplicationTemplateRepository{db: db}
}

func (r *ApplicationTemplateRepository) Create(ctx context.Context, item *domain.ApplicationTemplate) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ApplicationTemplateRepository) Update(ctx context.Context, item *domain.ApplicationTemplate) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ApplicationTemplateRepository) FindBySlug(ctx context.Context, slug string) (*domain.ApplicationTemplate, error) {
	var item domain.ApplicationTemplate
	if err := r.db.WithContext(ctx).Where("slug = ?", strings.TrimSpace(slug)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ApplicationTemplateRepository) FindBySourceAndSlug(ctx context.Context, sourceID uint64, slug string) (*domain.ApplicationTemplate, error) {
	var item domain.ApplicationTemplate
	if err := r.db.WithContext(ctx).Where("catalog_source_id = ? AND slug = ?", sourceID, strings.TrimSpace(slug)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ApplicationTemplateRepository) GetByID(ctx context.Context, id uint64) (*domain.ApplicationTemplate, error) {
	var item domain.ApplicationTemplate
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ApplicationTemplateRepository) List(ctx context.Context, filter ApplicationTemplateListFilter) ([]domain.ApplicationTemplate, error) {
	query := r.db.WithContext(ctx).Model(&domain.ApplicationTemplate{})
	if v := strings.TrimSpace(filter.Category); v != "" {
		query = query.Where("category = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("publish_status = ?", v)
	}
	if filter.CatalogSourceID != 0 {
		query = query.Where("catalog_source_id = ?", filter.CatalogSourceID)
	}
	if v := strings.TrimSpace(filter.Keyword); v != "" {
		query = query.Where("name LIKE ? OR slug LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	var items []domain.ApplicationTemplate
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}
