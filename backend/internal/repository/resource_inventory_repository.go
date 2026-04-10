package repository

import (
	"context"
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
	ClusterID uint64
	Namespace string
	Kind      string
	Health    string
	Keyword   string
	Limit     int
	Offset    int
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
