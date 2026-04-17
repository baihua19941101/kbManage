package repository

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceTrendRepository struct{ db *gorm.DB }

func NewComplianceTrendRepository(db *gorm.DB) *ComplianceTrendRepository {
	return &ComplianceTrendRepository{db: db}
}

func (r *ComplianceTrendRepository) Create(ctx context.Context, item *domain.ComplianceTrendSnapshot) error {
	if item == nil {
		return errors.New("trend snapshot is required")
	}
	if r == nil || r.db == nil {
		return errors.New("trend repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}
