package repository

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceRecheckRepository struct{ db *gorm.DB }

func NewComplianceRecheckRepository(db *gorm.DB) *ComplianceRecheckRepository {
	return &ComplianceRecheckRepository{db: db}
}

func (r *ComplianceRecheckRepository) Create(ctx context.Context, item *domain.RecheckTask) error {
	if item == nil {
		return errors.New("recheck task is required")
	}
	if r == nil || r.db == nil {
		return errors.New("recheck repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}
