package domain

import "time"

type ClusterStatus string

const (
	ClusterStatusUnknown  ClusterStatus = "unknown"
	ClusterStatusHealthy  ClusterStatus = "healthy"
	ClusterStatusDegraded ClusterStatus = "degraded"
)

type Cluster struct {
	ID        uint64        `gorm:"primaryKey"`
	Name      string        `gorm:"size:128;uniqueIndex;not null"`
	APIServer string        `gorm:"size:512;not null"`
	Status    ClusterStatus `gorm:"size:32;not null;default:unknown"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
