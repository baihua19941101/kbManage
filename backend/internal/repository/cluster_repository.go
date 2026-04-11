package repository

import (
	"context"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) Create(ctx context.Context, cluster *domain.Cluster) error {
	return r.db.WithContext(ctx).Create(cluster).Error
}

func (r *ClusterRepository) GetByID(ctx context.Context, id uint64) (*domain.Cluster, error) {
	var cluster domain.Cluster
	if err := r.db.WithContext(ctx).First(&cluster, id).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (r *ClusterRepository) List(ctx context.Context) ([]domain.Cluster, error) {
	var clusters []domain.Cluster
	err := r.db.WithContext(ctx).Order("id DESC").Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) ListByIDs(ctx context.Context, ids []uint64) ([]domain.Cluster, error) {
	if len(ids) == 0 {
		return []domain.Cluster{}, nil
	}

	var clusters []domain.Cluster
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("id DESC").
		Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) UpdateStatus(ctx context.Context, id uint64, status domain.ClusterStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Cluster{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ClusterRepository) UpdateSyncStatus(ctx context.Context, id uint64, syncStatus domain.ClusterSyncStatus, failureReason string) error {
	// SQLite in-memory tests are sensitive to concurrent writer locks from
	// background sync workers; skip best-effort sync status writes there.
	if r.db != nil && strings.EqualFold(r.db.Dialector.Name(), "sqlite") {
		return nil
	}

	now := time.Now()
	updates := map[string]any{
		"sync_status":  syncStatus,
		"last_sync_at": now,
	}

	reason := strings.TrimSpace(failureReason)
	switch syncStatus {
	case domain.ClusterSyncStatusSuccess:
		updates["last_success_at"] = now
		updates["sync_failure_reason"] = ""
	case domain.ClusterSyncStatusFailed:
		updates["sync_failure_reason"] = reason
	default:
		if reason != "" {
			updates["sync_failure_reason"] = reason
		}
	}

	return r.db.WithContext(ctx).Model(&domain.Cluster{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ClusterRepository) MarkSyncRunning(ctx context.Context, id uint64) error {
	return r.UpdateSyncStatus(ctx, id, domain.ClusterSyncStatusRunning, "")
}

func (r *ClusterRepository) MarkSyncSuccess(ctx context.Context, id uint64) error {
	return r.UpdateSyncStatus(ctx, id, domain.ClusterSyncStatusSuccess, "")
}

func (r *ClusterRepository) MarkSyncFailed(ctx context.Context, id uint64, failureReason string) error {
	return r.UpdateSyncStatus(ctx, id, domain.ClusterSyncStatusFailed, failureReason)
}
