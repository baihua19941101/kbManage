package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ResourceInventory is a denormalized indexed resource record.
type ResourceInventory struct {
	ID        uint64 `gorm:"primaryKey"`
	ClusterID uint64 `gorm:"index;not null"`
	Namespace string `gorm:"size:128;index"`
	Kind      string `gorm:"size:128;index;not null"`
	Name      string `gorm:"size:255;index;not null"`
	Health    string `gorm:"size:32;index;not null;default:unknown"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResourceListFilter struct {
	ClusterID  uint64
	ClusterIDs []uint64
	Namespace  string
	Kind       string
	Health     string
	Keyword    string
	Limit      int
	Offset     int
}

type ResourceDetailFilter struct {
	ClusterID uint64
	Namespace string
	Kind      string
	Name      string
}

type ResourceInventoryRepository struct {
	db *gorm.DB
}

func NewResourceInventoryRepository(db *gorm.DB) *ResourceInventoryRepository {
	return &ResourceInventoryRepository{db: db}
}

func (r *ResourceInventoryRepository) List(ctx context.Context, filter ResourceListFilter) ([]ResourceInventory, error) {
	q := r.db.WithContext(ctx).Model(&ResourceInventory{})
	if filter.ClusterID > 0 {
		q = q.Where("cluster_id = ?", filter.ClusterID)
	} else if len(filter.ClusterIDs) > 0 {
		q = q.Where("cluster_id IN ?", filter.ClusterIDs)
	}
	if filter.Namespace != "" {
		q = q.Where("namespace = ?", filter.Namespace)
	}
	if filter.Kind != "" {
		q = q.Where("kind = ?", filter.Kind)
	}
	if filter.Health != "" {
		q = q.Where("health = ?", strings.ToLower(filter.Health))
	}
	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		q = q.Where("name LIKE ?", kw)
	}
	if filter.Limit <= 0 || filter.Limit > 200 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var items []ResourceInventory
	err := q.Order("updated_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error
	return items, err
}

func (r *ResourceInventoryRepository) CountByHealth(ctx context.Context, clusterID uint64) (map[string]int64, error) {
	counts := map[string]int64{"healthy": 0, "degraded": 0, "unknown": 0}
	type row struct {
		Health string
		Total  int64
	}

	var rows []row
	err := r.db.WithContext(ctx).
		Model(&ResourceInventory{}).
		Select("health, COUNT(1) AS total").
		Where("cluster_id = ?", clusterID).
		Group("health").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		counts[strings.ToLower(item.Health)] = item.Total
	}
	return counts, nil
}

func (r *ResourceInventoryRepository) GetDetail(ctx context.Context, filter ResourceDetailFilter) (*ResourceInventory, error) {
	if filter.ClusterID == 0 {
		return nil, errors.New("clusterID is required")
	}
	if strings.TrimSpace(filter.Kind) == "" {
		return nil, errors.New("kind is required")
	}
	if strings.TrimSpace(filter.Name) == "" {
		return nil, errors.New("name is required")
	}

	var item ResourceInventory
	err := r.db.WithContext(ctx).
		Model(&ResourceInventory{}).
		Where("cluster_id = ? AND namespace = ? AND kind = ? AND name = ?",
			filter.ClusterID,
			strings.TrimSpace(filter.Namespace),
			strings.TrimSpace(filter.Kind),
			strings.TrimSpace(filter.Name),
		).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ReplaceClusterSnapshot atomically replaces all indexed resources for a cluster.
func (r *ResourceInventoryRepository) ReplaceClusterSnapshot(ctx context.Context, clusterID uint64, items []ResourceInventory) error {
	if clusterID == 0 {
		return errors.New("clusterID is required")
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("cluster_id = ?", clusterID).Delete(&ResourceInventory{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(items) > 0 {
		for i := range items {
			items[i].ID = 0
			items[i].ClusterID = clusterID
			items[i].Health = strings.ToLower(strings.TrimSpace(items[i].Health))
			if items[i].Health == "" {
				items[i].Health = "unknown"
			}
		}
		if err := tx.CreateInBatches(items, 200).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
