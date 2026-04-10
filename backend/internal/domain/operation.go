package domain

import "time"

type OperationStatus string

type RiskLevel string

const (
	OperationStatusPending   OperationStatus = "pending"
	OperationStatusRunning   OperationStatus = "running"
	OperationStatusSucceeded OperationStatus = "succeeded"
	OperationStatusFailed    OperationStatus = "failed"

	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

type OperationRequest struct {
	ID            uint64          `gorm:"primaryKey"`
	RequestID     string          `gorm:"size:64;uniqueIndex;not null"`
	OperatorID    uint64          `gorm:"index;not null"`
	OperationType string          `gorm:"size:128;not null"`
	TargetRef     string          `gorm:"size:255;not null"`
	Status        OperationStatus `gorm:"size:32;not null"`
	RiskLevel     RiskLevel       `gorm:"size:32;not null"`
	ResultMessage string          `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (OperationRequest) TableName() string { return "operation_requests" }
