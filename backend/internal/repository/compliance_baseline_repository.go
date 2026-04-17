package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceBaselineListFilter struct {
	StandardType domain.ComplianceStandardType
	Status       domain.ComplianceBaselineStatus
}

type ComplianceBaselineRepository struct{ db *gorm.DB }

func NewComplianceBaselineRepository(db *gorm.DB) *ComplianceBaselineRepository {
	return &ComplianceBaselineRepository{db: db}
}

func (r *ComplianceBaselineRepository) Create(ctx context.Context, item *domain.ComplianceBaseline) error {
	if item == nil {
		return errors.New("compliance baseline is required")
	}
	if r == nil || r.db == nil {
		return errors.New("compliance baseline repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ComplianceBaselineRepository) GetByID(ctx context.Context, id uint64) (*domain.ComplianceBaseline, error) {
	if id == 0 {
		return nil, errors.New("baseline id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("compliance baseline repository is not configured")
	}
	var item domain.ComplianceBaseline
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ComplianceBaselineRepository) List(ctx context.Context, filter ComplianceBaselineListFilter) ([]domain.ComplianceBaseline, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("compliance baseline repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.ComplianceBaseline{})
	if strings.TrimSpace(string(filter.StandardType)) != "" {
		query = query.Where("standard_type = ?", filter.StandardType)
	}
	if strings.TrimSpace(string(filter.Status)) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	items := make([]domain.ComplianceBaseline, 0)
	if err := query.Order("updated_at DESC, id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ComplianceBaselineRepository) UpdateFields(ctx context.Context, id uint64, updates map[string]any) error {
	if id == 0 {
		return errors.New("baseline id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("compliance baseline repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.ComplianceBaseline{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
