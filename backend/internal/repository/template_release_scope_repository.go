package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type TemplateReleaseScopeRepository struct{ db *gorm.DB }

type TemplateReleaseScopeListFilter struct {
	TemplateID uint64
	ScopeType  string
	ScopeRef   string
	Status     string
}

func NewTemplateReleaseScopeRepository(db *gorm.DB) *TemplateReleaseScopeRepository {
	return &TemplateReleaseScopeRepository{db: db}
}

func (r *TemplateReleaseScopeRepository) Create(ctx context.Context, item *domain.TemplateReleaseScope) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *TemplateReleaseScopeRepository) Update(ctx context.Context, item *domain.TemplateReleaseScope) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *TemplateReleaseScopeRepository) ListByTemplateID(ctx context.Context, templateID uint64) ([]domain.TemplateReleaseScope, error) {
	var items []domain.TemplateReleaseScope
	err := r.db.WithContext(ctx).Where("template_id = ?", templateID).Order("id DESC").Find(&items).Error
	return items, err
}

func (r *TemplateReleaseScopeRepository) List(ctx context.Context, filter TemplateReleaseScopeListFilter) ([]domain.TemplateReleaseScope, error) {
	query := r.db.WithContext(ctx).Model(&domain.TemplateReleaseScope{})
	if filter.TemplateID != 0 {
		query = query.Where("template_id = ?", filter.TemplateID)
	}
	if v := strings.TrimSpace(filter.ScopeType); v != "" {
		query = query.Where("scope_type = ?", v)
	}
	if v := strings.TrimSpace(filter.ScopeRef); v != "" {
		query = query.Where("scope_ref = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.TemplateReleaseScope
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *TemplateReleaseScopeRepository) FindByTemplateTarget(ctx context.Context, templateID uint64, targetType, targetRef string) (*domain.TemplateReleaseScope, error) {
	var item domain.TemplateReleaseScope
	if err := r.db.WithContext(ctx).Where("template_id = ? AND scope_type = ? AND scope_ref = ?", templateID, strings.TrimSpace(targetType), strings.TrimSpace(targetRef)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TemplateReleaseScopeRepository) FindByTemplateScope(ctx context.Context, templateID uint64, scopeType, scopeRef string) (*domain.TemplateReleaseScope, error) {
	return r.FindByTemplateTarget(ctx, templateID, scopeType, scopeRef)
}
