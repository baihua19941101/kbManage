package domain

import "time"

type ClusterStatus string

const (
	ClusterStatusUnknown  ClusterStatus = "unknown"
	ClusterStatusHealthy  ClusterStatus = "healthy"
	ClusterStatusDegraded ClusterStatus = "degraded"
)

type ClusterSyncStatus string

const (
	ClusterSyncStatusIdle    ClusterSyncStatus = "idle"
	ClusterSyncStatusRunning ClusterSyncStatus = "running"
	ClusterSyncStatusSuccess ClusterSyncStatus = "success"
	ClusterSyncStatusFailed  ClusterSyncStatus = "failed"
)

type Cluster struct {
	ID                uint64            `gorm:"primaryKey"`
	Name              string            `gorm:"size:128;uniqueIndex;not null"`
	APIServer         string            `gorm:"size:512;not null"`
	Status            ClusterStatus     `gorm:"size:32;not null;default:unknown"`
	SyncStatus        ClusterSyncStatus `gorm:"size:32;not null;default:idle"`
	LastSyncAt        *time.Time
	LastSuccessAt     *time.Time
	SyncFailureReason string `gorm:"size:1024"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
