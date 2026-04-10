package domain

import "time"

type ScopeType string

const (
	ScopeTypePlatform  ScopeType = "platform"
	ScopeTypeWorkspace ScopeType = "workspace"
	ScopeTypeProject   ScopeType = "project"
)

type PlatformRole struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"size:128;uniqueIndex;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPlatformRole struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"index;not null"`
	RoleID    uint64 `gorm:"index;not null"`
	CreatedAt time.Time
}

type Workspace struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"size:128;uniqueIndex;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Project struct {
	ID          uint64 `gorm:"primaryKey"`
	WorkspaceID uint64 `gorm:"index;not null"`
	Name        string `gorm:"size:128;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Project) TableName() string { return "projects" }
