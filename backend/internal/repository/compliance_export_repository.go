package repository

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceExportRepository struct{ db *gorm.DB }

func NewComplianceExportRepository(db *gorm.DB) *ComplianceExportRepository {
	return &ComplianceExportRepository{db: db}
}

func (r *ComplianceExportRepository) Create(ctx context.Context, item *domain.ArchiveExportTask) error {
	if item == nil {
		return errors.New("archive export task is required")
	}
	if r == nil || r.db == nil {
		return errors.New("archive export repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}
