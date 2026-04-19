package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type BackupAuditListFilter struct {
	WorkspaceIDs []uint64
	Action       string
	Outcome      string
	TargetType   string
}

type BackupAuditRepository struct {
	db *gorm.DB
}

func NewBackupAuditRepository(db *gorm.DB) *BackupAuditRepository {
	return &BackupAuditRepository{db: db}
}

func (r *BackupAuditRepository) Create(ctx context.Context, item *domain.BackupAuditEvent) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *BackupAuditRepository) List(ctx context.Context, filter BackupAuditListFilter) ([]domain.BackupAuditEvent, error) {
	query := r.db.WithContext(ctx).Model(&domain.BackupAuditEvent{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if v := strings.TrimSpace(filter.Action); v != "" {
		query = query.Where("action = ?", v)
	}
	if v := strings.TrimSpace(filter.Outcome); v != "" {
		query = query.Where("outcome = ?", v)
	}
	if v := strings.TrimSpace(filter.TargetType); v != "" {
		query = query.Where("target_type = ?", v)
	}
	var items []domain.BackupAuditEvent
	err := query.Order("occurred_at DESC, id DESC").Find(&items).Error
	return items, err
}
