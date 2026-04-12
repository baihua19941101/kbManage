package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type AlertIncidentRepository struct {
	db *gorm.DB
}

func NewAlertIncidentRepository(db *gorm.DB) *AlertIncidentRepository {
	return &AlertIncidentRepository{db: db}
}

func (r *AlertIncidentRepository) UpsertBySourceKey(ctx context.Context, item *domain.AlertIncidentSnapshot) error {
	tx := r.db.WithContext(ctx)
	var existing domain.AlertIncidentSnapshot
	err := tx.Where("source_incident_key = ?", item.SourceIncidentKey).First(&existing).Error
	if err == nil {
		item.ID = existing.ID
		return tx.Save(item).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return tx.Create(item).Error
}

func (r *AlertIncidentRepository) GetByID(ctx context.Context, id uint64) (*domain.AlertIncidentSnapshot, error) {
	var item domain.AlertIncidentSnapshot
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AlertIncidentRepository) Update(ctx context.Context, item *domain.AlertIncidentSnapshot) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *AlertIncidentRepository) List(ctx context.Context, status domain.AlertIncidentStatus, limit int) ([]domain.AlertIncidentSnapshot, error) {
	if limit <= 0 {
		limit = 100
	}
	tx := r.db.WithContext(ctx).Model(&domain.AlertIncidentSnapshot{})
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	var items []domain.AlertIncidentSnapshot
	err := tx.Order("id DESC").Limit(limit).Find(&items).Error
	return items, err
}

func (r *AlertIncidentRepository) CreateHandlingRecord(ctx context.Context, item *domain.AlertHandlingRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AlertIncidentRepository) ListHandlingRecords(ctx context.Context, incidentID uint64) ([]domain.AlertHandlingRecord, error) {
	var items []domain.AlertHandlingRecord
	err := r.db.WithContext(ctx).
		Where("incident_id = ?", incidentID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}
