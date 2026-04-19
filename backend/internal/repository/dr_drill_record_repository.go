package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DRDrillRecordRepository struct {
	db *gorm.DB
}

func NewDRDrillRecordRepository(db *gorm.DB) *DRDrillRecordRepository {
	return &DRDrillRecordRepository{db: db}
}

func (r *DRDrillRecordRepository) Create(ctx context.Context, item *domain.DRDrillRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DRDrillRecordRepository) Update(ctx context.Context, item *domain.DRDrillRecord) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *DRDrillRecordRepository) GetByID(ctx context.Context, id uint64) (*domain.DRDrillRecord, error) {
	var item domain.DRDrillRecord
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
