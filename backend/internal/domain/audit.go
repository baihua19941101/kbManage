package domain

import (
	"encoding/json"
	"time"
)

type AuditOutcome string

const (
	AuditOutcomeSuccess AuditOutcome = "success"
	AuditOutcomeDenied  AuditOutcome = "denied"
	AuditOutcomeFailed  AuditOutcome = "failed"
)

type AuditEvent struct {
	ID            uint64          `gorm:"primaryKey"`
	RequestID     string          `gorm:"size:64;index;not null"`
	ActorID       *uint64         `gorm:"index"`
	ClusterID     *uint64         `gorm:"index"`
	WorkspaceID   *uint64         `gorm:"index"`
	ProjectID     *uint64         `gorm:"index"`
	AuditCategory string          `gorm:"size:64;index"`
	ActionScope   string          `gorm:"size:64;index"`
	ScopeSnapshot json.RawMessage `gorm:"type:json"`
	SearchTags    string          `gorm:"size:512;index"`
	Action        string          `gorm:"size:128;not null"`
	ResourceType  string          `gorm:"size:128"`
	ResourceID    string          `gorm:"size:128"`
	Outcome       AuditOutcome    `gorm:"size:32;not null"`
	Details       json.RawMessage `gorm:"type:json"`
	CreatedAt     time.Time
}
