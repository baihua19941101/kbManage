package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type TemplateVersionRepository struct{ db *gorm.DB }

func NewTemplateVersionRepository(db *gorm.DB) *TemplateVersionRepository {
	return &TemplateVersionRepository{db: db}
}

func (r *TemplateVersionRepository) Create(ctx context.Context, item *domain.TemplateVersion) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *TemplateVersionRepository) Update(ctx context.Context, item *domain.TemplateVersion) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *TemplateVersionRepository) FindByTemplateVersion(ctx context.Context, templateID uint64, version string) (*domain.TemplateVersion, error) {
	var item domain.TemplateVersion
	if err := r.db.WithContext(ctx).Where("template_id = ? AND version = ?", templateID, version).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TemplateVersionRepository) FindByTemplateAndVersion(ctx context.Context, templateID uint64, version string) (*domain.TemplateVersion, error) {
	return r.FindByTemplateVersion(ctx, templateID, version)
}

func (r *TemplateVersionRepository) GetByID(ctx context.Context, id uint64) (*domain.TemplateVersion, error) {
	var item domain.TemplateVersion
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TemplateVersionRepository) ListByTemplateID(ctx context.Context, templateID uint64) ([]domain.TemplateVersion, error) {
	var items []domain.TemplateVersion
	err := r.db.WithContext(ctx).Where("template_id = ?", templateID).Order("id ASC").Find(&items).Error
	return items, err
}
