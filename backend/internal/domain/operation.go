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
	ID              uint64          `gorm:"primaryKey"`
	RequestID       string          `gorm:"size:64;uniqueIndex;not null"`
	OperatorID      uint64          `gorm:"index;not null"`
	OperationType   string          `gorm:"size:128;not null"`
	TargetRef       string          `gorm:"size:255;not null"`
	Status          OperationStatus `gorm:"size:32;not null"`
	RiskLevel       RiskLevel       `gorm:"size:32;not null"`
	ProgressMessage string          `gorm:"type:text"`
	ResultMessage   string          `gorm:"type:text"`
	FailureReason   string          `gorm:"type:text"`
	CompletedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (OperationRequest) TableName() string { return "operation_requests" }

func (s OperationStatus) IsTerminal() bool {
	return s == OperationStatusSucceeded || s == OperationStatusFailed
}

func (s OperationStatus) CanTransitTo(next OperationStatus) bool {
	if s == next {
		return true
	}
	switch s {
	case OperationStatusPending:
		return next == OperationStatusRunning
	case OperationStatusRunning:
		return next == OperationStatusSucceeded || next == OperationStatusFailed
	default:
		return false
	}
}
