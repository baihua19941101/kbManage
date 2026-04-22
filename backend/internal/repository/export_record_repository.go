package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ExportRecordRepository struct{ db *gorm.DB }

func NewExportRecordRepository(db *gorm.DB) *ExportRecordRepository {
	return &ExportRecordRepository{db: db}
}

func (r *ExportRecordRepository) Create(ctx context.Context, item *domain.ExportRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ExportRecordRepository) List(ctx context.Context) ([]domain.ExportRecord, error) {
	var items []domain.ExportRecord
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
