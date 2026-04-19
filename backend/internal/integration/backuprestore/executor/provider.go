package executor

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type BackupExecutionRequest struct {
	PolicyID         uint64
	PolicyName       string
	ScopeType        string
	ScopeRef         string
	ConsistencyLevel string
}

type BackupExecutionResult struct {
	Result             string
	ConsistencySummary string
	FailureReason      string
	StorageRef         string
	DurationSeconds    int
	ExpiresAt          *time.Time
}

type RestoreExecutionRequest struct {
	JobID             uint64
	JobType           string
	RestorePointID    uint64
	TargetEnvironment string
}

type RestoreExecutionResult struct {
	Status            string
	ResultSummary     string
	FailureReason     string
	ConflictSummary   string
	ConsistencyNotice string
}

type MigrationExecutionRequest struct {
	PlanID          uint64
	Name            string
	SourceClusterID uint64
	TargetClusterID uint64
}

type MigrationExecutionResult struct {
	Status        string
	ResultSummary string
	FailureReason string
}

type DrillExecutionRequest struct {
	PlanID   uint64
	PlanName string
}

type DrillExecutionResult struct {
	Status            string
	ActualRPOMinutes  int
	ActualRTOMinutes  int
	StepResults       []string
	ValidationResults []string
	IncidentNotes     string
}

type Provider interface {
	RunBackup(context.Context, BackupExecutionRequest) (BackupExecutionResult, error)
	RunRestore(context.Context, RestoreExecutionRequest) (RestoreExecutionResult, error)
	RunMigration(context.Context, MigrationExecutionRequest) (MigrationExecutionResult, error)
	RunDrill(context.Context, DrillExecutionRequest) (DrillExecutionResult, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return &StaticProvider{}
}

func (p *StaticProvider) RunBackup(_ context.Context, req BackupExecutionRequest) (BackupExecutionResult, error) {
	expires := time.Now().Add(7 * 24 * time.Hour)
	return BackupExecutionResult{
		Result:             "succeeded",
		ConsistencySummary: fmt.Sprintf("%s 范围备份完成，可用于恢复。", strings.TrimSpace(req.ScopeType)),
		StorageRef:         fmt.Sprintf("snapshot://policy/%d/%d", req.PolicyID, time.Now().Unix()),
		DurationSeconds:    12,
		ExpiresAt:          &expires,
	}, nil
}

func (p *StaticProvider) RunRestore(_ context.Context, req RestoreExecutionRequest) (RestoreExecutionResult, error) {
	return RestoreExecutionResult{
		Status:            "succeeded",
		ResultSummary:     fmt.Sprintf("%s 已恢复到 %s", strings.TrimSpace(req.JobType), strings.TrimSpace(req.TargetEnvironment)),
		ConflictSummary:   "未发现阻断冲突",
		ConsistencyNotice: "恢复结果基于最近一次成功恢复点，需在业务侧完成最终验收",
	}, nil
}

func (p *StaticProvider) RunMigration(_ context.Context, req MigrationExecutionRequest) (MigrationExecutionResult, error) {
	return MigrationExecutionResult{
		Status:        "succeeded",
		ResultSummary: fmt.Sprintf("迁移计划 %s 已从集群 %d 切换到集群 %d", strings.TrimSpace(req.Name), req.SourceClusterID, req.TargetClusterID),
	}, nil
}

func (p *StaticProvider) RunDrill(_ context.Context, req DrillExecutionRequest) (DrillExecutionResult, error) {
	return DrillExecutionResult{
		Status:            "succeeded",
		ActualRPOMinutes:  8,
		ActualRTOMinutes:  16,
		StepResults:       []string{fmt.Sprintf("演练计划 %s 切换步骤执行完成", strings.TrimSpace(req.PlanName))},
		ValidationResults: []string{"关键命名空间健康检查通过", "恢复后权限抽样验证通过"},
	}, nil
}
