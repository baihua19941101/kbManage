package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DRDrillReportRepository struct {
	db *gorm.DB
}

func NewDRDrillReportRepository(db *gorm.DB) *DRDrillReportRepository {
	return &DRDrillReportRepository{db: db}
}

func (r *DRDrillReportRepository) Create(ctx context.Context, item *domain.DRDrillReport) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DRDrillReportRepository) Update(ctx context.Context, item *domain.DRDrillReport) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *DRDrillReportRepository) GetByRecordID(ctx context.Context, recordID uint64) (*domain.DRDrillReport, error) {
	var item domain.DRDrillReport
	if err := r.db.WithContext(ctx).Where("drill_record_id = ?", recordID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
