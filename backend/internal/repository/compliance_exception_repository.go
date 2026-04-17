package repository

import (
	"context"
	"errors"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceExceptionRepository struct{ db *gorm.DB }

func NewComplianceExceptionRepository(db *gorm.DB) *ComplianceExceptionRepository {
	return &ComplianceExceptionRepository{db: db}
}

func (r *ComplianceExceptionRepository) Create(ctx context.Context, item *domain.ComplianceExceptionRequest) error {
	if item == nil {
		return errors.New("exception request is required")
	}
	if r == nil || r.db == nil {
		return errors.New("exception repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ComplianceExceptionRepository) ExpireActiveBefore(ctx context.Context, now time.Time) ([]uint64, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("exception repository is not configured")
	}
	items := make([]domain.ComplianceExceptionRequest, 0)
	if err := r.db.WithContext(ctx).Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?", domain.ComplianceExceptionStatusActive, now).Find(&items).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	if len(ids) == 0 {
		return ids, nil
	}
	if err := r.db.WithContext(ctx).Model(&domain.ComplianceExceptionRequest{}).Where("id IN ?", ids).Update("status", domain.ComplianceExceptionStatusExpired).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
