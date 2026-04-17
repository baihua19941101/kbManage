package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceScanExecutionListFilter struct {
	WorkspaceID   *uint64
	ProjectID     *uint64
	ProfileID     *uint64
	Status        domain.ComplianceScanStatus
	TriggerSource domain.ComplianceTriggerSource
	From          *time.Time
	To            *time.Time
}

type ComplianceScanExecutionRepository struct{ db *gorm.DB }

func NewComplianceScanExecutionRepository(db *gorm.DB) *ComplianceScanExecutionRepository {
	return &ComplianceScanExecutionRepository{db: db}
}

func (r *ComplianceScanExecutionRepository) Create(ctx context.Context, item *domain.ScanExecution) error {
	if item == nil {
		return errors.New("scan execution is required")
	}
	if r == nil || r.db == nil {
		return errors.New("scan execution repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ComplianceScanExecutionRepository) GetByID(ctx context.Context, id uint64) (*domain.ScanExecution, error) {
	if id == 0 {
		return nil, errors.New("scan id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("scan execution repository is not configured")
	}
	var item domain.ScanExecution
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ComplianceScanExecutionRepository) List(ctx context.Context, filter ComplianceScanExecutionListFilter) ([]domain.ScanExecution, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("scan execution repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.ScanExecution{})
	if filter.WorkspaceID != nil {
		query = query.Where("workspace_id = ?", *filter.WorkspaceID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.ProfileID != nil {
		query = query.Where("profile_id = ?", *filter.ProfileID)
	}
	if strings.TrimSpace(string(filter.Status)) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if strings.TrimSpace(string(filter.TriggerSource)) != "" {
		query = query.Where("trigger_source = ?", filter.TriggerSource)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}
	items := make([]domain.ScanExecution, 0)
	if err := query.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ComplianceScanExecutionRepository) ListPending(ctx context.Context, limit int) ([]domain.ScanExecution, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("scan execution repository is not configured")
	}
	if limit <= 0 {
		limit = 20
	}
	items := make([]domain.ScanExecution, 0, limit)
	if err := r.db.WithContext(ctx).Where("status IN ?", []domain.ComplianceScanStatus{domain.ComplianceScanStatusPending, domain.ComplianceScanStatusRunning}).Order("id ASC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ComplianceScanExecutionRepository) UpdateFields(ctx context.Context, id uint64, updates map[string]any) error {
	if id == 0 {
		return errors.New("scan id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("scan execution repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.ScanExecution{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
