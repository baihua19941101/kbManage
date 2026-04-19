package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type InstallationRecordListFilter struct {
	TemplateID uint64
	ScopeType  string
	ScopeRef   string
	Status     string
}

type InstallationRecordRepository struct{ db *gorm.DB }

func NewInstallationRecordRepository(db *gorm.DB) *InstallationRecordRepository {
	return &InstallationRecordRepository{db: db}
}

func (r *InstallationRecordRepository) Create(ctx context.Context, item *domain.InstallationRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *InstallationRecordRepository) Update(ctx context.Context, item *domain.InstallationRecord) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *InstallationRecordRepository) List(ctx context.Context, filter InstallationRecordListFilter) ([]domain.InstallationRecord, error) {
	query := r.db.WithContext(ctx).Model(&domain.InstallationRecord{})
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
		query = query.Where("lifecycle_status = ?", v)
	}
	var items []domain.InstallationRecord
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *InstallationRecordRepository) FindByTemplateScope(ctx context.Context, templateID uint64, scopeType, scopeRef string) (*domain.InstallationRecord, error) {
	var item domain.InstallationRecord
	if err := r.db.WithContext(ctx).
		Where("template_id = ? AND scope_type = ? AND scope_ref = ?", templateID, strings.TrimSpace(scopeType), strings.TrimSpace(scopeRef)).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
