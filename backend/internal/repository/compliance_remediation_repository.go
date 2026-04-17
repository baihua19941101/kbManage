package repository

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceRemediationRepository struct{ db *gorm.DB }

func NewComplianceRemediationRepository(db *gorm.DB) *ComplianceRemediationRepository {
	return &ComplianceRemediationRepository{db: db}
}

func (r *ComplianceRemediationRepository) Create(ctx context.Context, item *domain.RemediationTask) error {
	if item == nil {
		return errors.New("remediation task is required")
	}
	if r == nil || r.db == nil {
		return errors.New("remediation repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}
