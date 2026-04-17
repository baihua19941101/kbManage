package domain

import "time"

type ComplianceStandardType string

type ComplianceBaselineStatus string

type ComplianceScopeType string

type ComplianceScheduleMode string

type ComplianceScanProfileStatus string

type ComplianceTriggerSource string

type ComplianceScanStatus string

type ComplianceCoverageStatus string

type ComplianceFindingResult string

type ComplianceRiskLevel string

type ComplianceRemediationStatus string

type ComplianceEvidenceType string

type ComplianceEvidenceConfidence string

type ComplianceEvidenceRedactionStatus string

type ComplianceRemediationTaskStatus string

type ComplianceExceptionStatus string

type ComplianceRecheckStatus string

type ComplianceRecheckTriggerSource string

type ComplianceArchiveExportScope string

type ComplianceArchiveExportStatus string

const (
	ComplianceStandardTypeCIS              ComplianceStandardType = "cis"
	ComplianceStandardTypeSTIG             ComplianceStandardType = "stig"
	ComplianceStandardTypePlatformBaseline ComplianceStandardType = "platform-baseline"

	ComplianceBaselineStatusDraft    ComplianceBaselineStatus = "draft"
	ComplianceBaselineStatusActive   ComplianceBaselineStatus = "active"
	ComplianceBaselineStatusDisabled ComplianceBaselineStatus = "disabled"
	ComplianceBaselineStatusArchived ComplianceBaselineStatus = "archived"

	ComplianceScopeTypeCluster     ComplianceScopeType = "cluster"
	ComplianceScopeTypeNode        ComplianceScopeType = "node"
	ComplianceScopeTypeNamespace   ComplianceScopeType = "namespace"
	ComplianceScopeTypeResourceSet ComplianceScopeType = "resource-set"

	ComplianceScheduleModeManual    ComplianceScheduleMode = "manual"
	ComplianceScheduleModeScheduled ComplianceScheduleMode = "scheduled"

	ComplianceScanProfileStatusDraft    ComplianceScanProfileStatus = "draft"
	ComplianceScanProfileStatusActive   ComplianceScanProfileStatus = "active"
	ComplianceScanProfileStatusPaused   ComplianceScanProfileStatus = "paused"
	ComplianceScanProfileStatusArchived ComplianceScanProfileStatus = "archived"

	ComplianceTriggerSourceManual   ComplianceTriggerSource = "manual"
	ComplianceTriggerSourceSchedule ComplianceTriggerSource = "schedule"
	ComplianceTriggerSourceRecheck  ComplianceTriggerSource = "recheck"

	ComplianceScanStatusPending            ComplianceScanStatus = "pending"
	ComplianceScanStatusRunning            ComplianceScanStatus = "running"
	ComplianceScanStatusPartiallySucceeded ComplianceScanStatus = "partially_succeeded"
	ComplianceScanStatusSucceeded          ComplianceScanStatus = "succeeded"
	ComplianceScanStatusFailed             ComplianceScanStatus = "failed"
	ComplianceScanStatusCanceled           ComplianceScanStatus = "canceled"

	ComplianceCoverageStatusFull        ComplianceCoverageStatus = "full"
	ComplianceCoverageStatusPartial     ComplianceCoverageStatus = "partial"
	ComplianceCoverageStatusUnavailable ComplianceCoverageStatus = "unavailable"

	ComplianceFindingResultPass    ComplianceFindingResult = "pass"
	ComplianceFindingResultFail    ComplianceFindingResult = "fail"
	ComplianceFindingResultWarn    ComplianceFindingResult = "warn"
	ComplianceFindingResultSkipped ComplianceFindingResult = "skipped"
	ComplianceFindingResultError   ComplianceFindingResult = "error"

	ComplianceRiskLevelLow      ComplianceRiskLevel = "low"
	ComplianceRiskLevelMedium   ComplianceRiskLevel = "medium"
	ComplianceRiskLevelHigh     ComplianceRiskLevel = "high"
	ComplianceRiskLevelCritical ComplianceRiskLevel = "critical"

	ComplianceRemediationStatusOpen            ComplianceRemediationStatus = "open"
	ComplianceRemediationStatusInProgress      ComplianceRemediationStatus = "in_progress"
	ComplianceRemediationStatusExceptionActive ComplianceRemediationStatus = "exception_active"
	ComplianceRemediationStatusReadyForRecheck ComplianceRemediationStatus = "ready_for_recheck"
	ComplianceRemediationStatusClosed          ComplianceRemediationStatus = "closed"

	ComplianceEvidenceTypeConfiguration ComplianceEvidenceType = "configuration"
	ComplianceEvidenceTypeWorkloadState ComplianceEvidenceType = "workload-state"
	ComplianceEvidenceTypeNodeState     ComplianceEvidenceType = "node-state"
	ComplianceEvidenceTypePermission    ComplianceEvidenceType = "permission"
	ComplianceEvidenceTypeNetwork       ComplianceEvidenceType = "network"
	ComplianceEvidenceTypeScannerOutput ComplianceEvidenceType = "scanner-output"

	ComplianceEvidenceConfidenceHigh   ComplianceEvidenceConfidence = "high"
	ComplianceEvidenceConfidenceMedium ComplianceEvidenceConfidence = "medium"
	ComplianceEvidenceConfidenceLow    ComplianceEvidenceConfidence = "low"

	ComplianceEvidenceRedactionRaw    ComplianceEvidenceRedactionStatus = "raw"
	ComplianceEvidenceRedactionMasked ComplianceEvidenceRedactionStatus = "masked"

	ComplianceRemediationTaskStatusTodo       ComplianceRemediationTaskStatus = "todo"
	ComplianceRemediationTaskStatusInProgress ComplianceRemediationTaskStatus = "in_progress"
	ComplianceRemediationTaskStatusBlocked    ComplianceRemediationTaskStatus = "blocked"
	ComplianceRemediationTaskStatusDone       ComplianceRemediationTaskStatus = "done"
	ComplianceRemediationTaskStatusCanceled   ComplianceRemediationTaskStatus = "canceled"

	ComplianceExceptionStatusPending  ComplianceExceptionStatus = "pending"
	ComplianceExceptionStatusApproved ComplianceExceptionStatus = "approved"
	ComplianceExceptionStatusRejected ComplianceExceptionStatus = "rejected"
	ComplianceExceptionStatusActive   ComplianceExceptionStatus = "active"
	ComplianceExceptionStatusExpired  ComplianceExceptionStatus = "expired"
	ComplianceExceptionStatusRevoked  ComplianceExceptionStatus = "revoked"

	ComplianceRecheckStatusPending  ComplianceRecheckStatus = "pending"
	ComplianceRecheckStatusRunning  ComplianceRecheckStatus = "running"
	ComplianceRecheckStatusPassed   ComplianceRecheckStatus = "passed"
	ComplianceRecheckStatusFailed   ComplianceRecheckStatus = "failed"
	ComplianceRecheckStatusCanceled ComplianceRecheckStatus = "canceled"

	ComplianceRecheckTriggerSourceManual           ComplianceRecheckTriggerSource = "manual"
	ComplianceRecheckTriggerSourceRemediationDone  ComplianceRecheckTriggerSource = "remediation_done"
	ComplianceRecheckTriggerSourceExceptionExpired ComplianceRecheckTriggerSource = "exception_expired"

	ComplianceArchiveExportScopeScans    ComplianceArchiveExportScope = "scans"
	ComplianceArchiveExportScopeFindings ComplianceArchiveExportScope = "findings"
	ComplianceArchiveExportScopeTrends   ComplianceArchiveExportScope = "trends"
	ComplianceArchiveExportScopeAudit    ComplianceArchiveExportScope = "audit"
	ComplianceArchiveExportScopeBundle   ComplianceArchiveExportScope = "bundle"

	ComplianceArchiveExportStatusPending   ComplianceArchiveExportStatus = "pending"
	ComplianceArchiveExportStatusRunning   ComplianceArchiveExportStatus = "running"
	ComplianceArchiveExportStatusSucceeded ComplianceArchiveExportStatus = "succeeded"
	ComplianceArchiveExportStatusFailed    ComplianceArchiveExportStatus = "failed"
	ComplianceArchiveExportStatusExpired   ComplianceArchiveExportStatus = "expired"
)

type ComplianceBaseline struct {
	ID               uint64                   `gorm:"primaryKey"`
	Name             string                   `gorm:"size:128;not null"`
	StandardType     ComplianceStandardType   `gorm:"size:32;not null"`
	Version          string                   `gorm:"size:64;not null"`
	Description      string                   `gorm:"type:text"`
	TargetLevelsJSON string                   `gorm:"type:longtext"`
	RulesJSON        string                   `gorm:"type:longtext"`
	RuleCount        int                      `gorm:"not null;default:0"`
	Status           ComplianceBaselineStatus `gorm:"size:32;not null;default:draft"`
	CreatedBy        *uint64                  `gorm:"index"`
	UpdatedBy        *uint64                  `gorm:"index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ScanProfile struct {
	ID                uint64                      `gorm:"primaryKey"`
	Name              string                      `gorm:"size:128;not null"`
	BaselineID        uint64                      `gorm:"index;not null"`
	WorkspaceID       *uint64                     `gorm:"index"`
	ProjectID         *uint64                     `gorm:"index"`
	ScopeType         ComplianceScopeType         `gorm:"size:32;not null"`
	ClusterRefsJSON   string                      `gorm:"type:longtext"`
	NodeSelectorsJSON string                      `gorm:"type:longtext"`
	NamespaceRefsJSON string                      `gorm:"type:longtext"`
	ResourceKindsJSON string                      `gorm:"type:longtext"`
	ScheduleMode      ComplianceScheduleMode      `gorm:"size:32;not null;default:manual"`
	CronExpression    string                      `gorm:"size:128"`
	Status            ComplianceScanProfileStatus `gorm:"size:32;not null;default:draft"`
	LastRunAt         *time.Time
	CreatedBy         *uint64 `gorm:"index"`
	UpdatedBy         *uint64 `gorm:"index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ScanExecution struct {
	ID                   uint64                   `gorm:"primaryKey"`
	ProfileID            uint64                   `gorm:"index;not null"`
	BaselineID           uint64                   `gorm:"index;not null"`
	WorkspaceID          *uint64                  `gorm:"index"`
	ProjectID            *uint64                  `gorm:"index"`
	BaselineSnapshotJSON string                   `gorm:"type:longtext"`
	BaselineVersionLabel string                   `gorm:"size:128"`
	TriggerSource        ComplianceTriggerSource  `gorm:"size:32;not null"`
	Status               ComplianceScanStatus     `gorm:"size:32;not null;default:pending"`
	CoverageStatus       ComplianceCoverageStatus `gorm:"size:32;not null;default:unavailable"`
	ScopeSnapshotJSON    string                   `gorm:"type:longtext"`
	StartedAt            *time.Time
	CompletedAt          *time.Time
	Score                float64 `gorm:"not null;default:0"`
	PassCount            int     `gorm:"not null;default:0"`
	FailCount            int     `gorm:"not null;default:0"`
	WarningCount         int     `gorm:"not null;default:0"`
	ErrorSummary         string  `gorm:"type:text"`
	CreatedBy            *uint64 `gorm:"index"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ComplianceFinding struct {
	ID                uint64                      `gorm:"primaryKey"`
	ScanExecutionID   uint64                      `gorm:"index;not null"`
	ControlID         string                      `gorm:"size:128;not null"`
	ControlTitle      string                      `gorm:"size:255;not null"`
	Result            ComplianceFindingResult     `gorm:"size:32;not null"`
	RiskLevel         ComplianceRiskLevel         `gorm:"size:32;not null"`
	ClusterID         *uint64                     `gorm:"index"`
	NodeName          string                      `gorm:"size:255"`
	Namespace         string                      `gorm:"size:255"`
	ResourceKind      string                      `gorm:"size:128"`
	ResourceName      string                      `gorm:"size:255"`
	ResourceUID       string                      `gorm:"size:255"`
	Summary           string                      `gorm:"type:text"`
	RemediationStatus ComplianceRemediationStatus `gorm:"size:32;not null;default:open"`
	DetectedAt        time.Time                   `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type EvidenceRecord struct {
	ID              uint64                            `gorm:"primaryKey"`
	FindingID       uint64                            `gorm:"index;not null"`
	EvidenceType    ComplianceEvidenceType            `gorm:"size:32;not null"`
	SourceRef       string                            `gorm:"size:512"`
	CollectedAt     time.Time                         `gorm:"not null"`
	Confidence      ComplianceEvidenceConfidence      `gorm:"size:32;not null"`
	Summary         string                            `gorm:"type:text"`
	ArtifactRef     string                            `gorm:"size:1024"`
	RedactionStatus ComplianceEvidenceRedactionStatus `gorm:"size:32;not null;default:masked"`
	PayloadJSON     string                            `gorm:"type:longtext"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RemediationTask struct {
	ID                uint64                          `gorm:"primaryKey"`
	FindingID         uint64                          `gorm:"index;not null"`
	Title             string                          `gorm:"size:255;not null"`
	Owner             string                          `gorm:"size:255"`
	Priority          ComplianceRiskLevel             `gorm:"size:32;not null;default:medium"`
	Status            ComplianceRemediationTaskStatus `gorm:"size:32;not null;default:todo"`
	DueAt             *time.Time
	ResolutionSummary string  `gorm:"type:text"`
	CreatedBy         *uint64 `gorm:"index"`
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ComplianceExceptionRequest struct {
	ID                uint64                    `gorm:"primaryKey"`
	FindingID         uint64                    `gorm:"index;not null"`
	ScopeSnapshotJSON string                    `gorm:"type:longtext"`
	Reason            string                    `gorm:"type:text"`
	RequestedBy       *uint64                   `gorm:"index"`
	ReviewedBy        *uint64                   `gorm:"index"`
	Status            ComplianceExceptionStatus `gorm:"size:32;not null;default:pending"`
	StartsAt          *time.Time
	ExpiresAt         *time.Time
	ReviewComment     string `gorm:"type:text"`
	ReviewedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type RecheckTask struct {
	ID                      uint64                         `gorm:"primaryKey"`
	FindingID               uint64                         `gorm:"index;not null"`
	TriggerSource           ComplianceRecheckTriggerSource `gorm:"size:32;not null"`
	Status                  ComplianceRecheckStatus        `gorm:"size:32;not null;default:pending"`
	TargetScopeSnapshotJSON string                         `gorm:"type:longtext"`
	ResultScanExecutionID   *uint64                        `gorm:"index"`
	RequestedBy             *uint64                        `gorm:"index"`
	StartedAt               *time.Time
	CompletedAt             *time.Time
	Summary                 string `gorm:"type:text"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type ComplianceTrendSnapshot struct {
	ID                        uint64              `gorm:"primaryKey"`
	WorkspaceID               *uint64             `gorm:"index"`
	ProjectID                 *uint64             `gorm:"index"`
	ScopeType                 ComplianceScopeType `gorm:"size:32;not null"`
	ScopeRef                  string              `gorm:"size:255;not null"`
	BaselineID                *uint64             `gorm:"index"`
	BaselineVersion           string              `gorm:"size:128"`
	WindowStart               time.Time           `gorm:"not null"`
	WindowEnd                 time.Time           `gorm:"not null"`
	CoverageRate              float64             `gorm:"not null;default:0"`
	ScoreAvg                  float64             `gorm:"not null;default:0"`
	OpenFindingsCount         int                 `gorm:"not null;default:0"`
	HighRiskOpenCount         int                 `gorm:"not null;default:0"`
	RemediationCompletionRate float64             `gorm:"not null;default:0"`
	ExceptionActiveCount      int                 `gorm:"not null;default:0"`
	GeneratedAt               time.Time           `gorm:"not null"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type ArchiveExportTask struct {
	ID                  uint64                        `gorm:"primaryKey"`
	WorkspaceID         *uint64                       `gorm:"index"`
	ProjectID           *uint64                       `gorm:"index"`
	BaselineID          *uint64                       `gorm:"index"`
	ExportScope         ComplianceArchiveExportScope  `gorm:"size:32;not null"`
	FiltersSnapshotJSON string                        `gorm:"type:longtext"`
	Status              ComplianceArchiveExportStatus `gorm:"size:32;not null;default:pending"`
	ArtifactRef         string                        `gorm:"size:1024"`
	RequestedBy         *uint64                       `gorm:"index"`
	StartedAt           *time.Time
	CompletedAt         *time.Time
	FailureReason       string `gorm:"type:text"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ComplianceBaselineSnapshot struct {
	BaselineID   uint64                 `json:"baselineId"`
	Name         string                 `json:"name"`
	StandardType ComplianceStandardType `json:"standardType"`
	Version      string                 `json:"version"`
	VersionLabel string                 `json:"versionLabel"`
	RuleCount    int                    `json:"ruleCount"`
	TargetLevels []string               `json:"targetLevels"`
	Rules        map[string]any         `json:"rules,omitempty"`
}

func (ComplianceBaseline) TableName() string         { return "compliance_baselines" }
func (ScanProfile) TableName() string                { return "compliance_scan_profiles" }
func (ScanExecution) TableName() string              { return "compliance_scan_executions" }
func (ComplianceFinding) TableName() string          { return "compliance_findings" }
func (EvidenceRecord) TableName() string             { return "compliance_evidence_records" }
func (RemediationTask) TableName() string            { return "compliance_remediation_tasks" }
func (ComplianceExceptionRequest) TableName() string { return "compliance_exception_requests" }
func (RecheckTask) TableName() string                { return "compliance_recheck_tasks" }
func (ComplianceTrendSnapshot) TableName() string    { return "compliance_trend_snapshots" }
func (ArchiveExportTask) TableName() string          { return "compliance_archive_export_tasks" }
