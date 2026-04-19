package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type SessionRecordListFilter struct {
	Status    string
	RiskLevel string
	UserID    uint64
}

type SessionRecordRepository struct {
	db *gorm.DB
}

func NewSessionRecordRepository(db *gorm.DB) *SessionRecordRepository {
	return &SessionRecordRepository{db: db}
}

func (r *SessionRecordRepository) Create(ctx context.Context, item *domain.SessionRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *SessionRecordRepository) Update(ctx context.Context, item *domain.SessionRecord) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *SessionRecordRepository) GetByID(ctx context.Context, id uint64) (*domain.SessionRecord, error) {
	var item domain.SessionRecord
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SessionRecordRepository) List(ctx context.Context, filter SessionRecordListFilter) ([]domain.SessionRecord, error) {
	query := r.db.WithContext(ctx).Model(&domain.SessionRecord{})
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.RiskLevel); v != "" {
		query = query.Where("risk_level = ?", v)
	}
	var items []domain.SessionRecord
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *SessionRecordRepository) FindByUserSource(ctx context.Context, userID, sourceID uint64) (*domain.SessionRecord, error) {
	var item domain.SessionRecord
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND identity_source_id = ?", userID, sourceID).
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
