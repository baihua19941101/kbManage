package gitops

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type StageExecutionResult struct {
	Environment    string `json:"environment"`
	Status         string `json:"status"`
	TargetCount    int    `json:"targetCount"`
	SucceededCount int    `json:"succeededCount"`
	FailedCount    int    `json:"failedCount"`
	FailureReason  string `json:"failureReason,omitempty"`
}

type StageExecutionService struct {
	units   *repository.DeliveryUnitRepository
	targets *repository.ClusterTargetGroupRepository
}

func NewStageExecutionService(
	units *repository.DeliveryUnitRepository,
	targets *repository.ClusterTargetGroupRepository,
) *StageExecutionService {
	return &StageExecutionService{units: units, targets: targets}
}

func (s *StageExecutionService) ExecuteStage(
	ctx context.Context,
	stage domain.EnvironmentStage,
	payload map[string]any,
) (StageExecutionResult, error) {
	res := StageExecutionResult{
		Environment: strings.TrimSpace(stage.Name),
		Status:      string(domain.EnvironmentStageStatusSucceeded),
		TargetCount: 1,
	}
	if s == nil || s.units == nil {
		return res, ErrGitOpsNotConfigured
	}
	if stage.ID == 0 {
		return res, fmt.Errorf("environment stage id is required")
	}
	if count := s.estimateTargetCount(ctx, stage.TargetGroupID); count > 0 {
		res.TargetCount = count
	}

	now := time.Now()
	if err := s.units.UpdateEnvironmentStage(ctx, stage.ID, map[string]any{
		"status":          domain.EnvironmentStageStatusProgressing,
		"last_entered_at": now,
	}); err != nil {
		return res, err
	}

	if stage.Paused {
		res.Status = string(domain.EnvironmentStageStatusFailed)
		res.FailedCount = res.TargetCount
		res.FailureReason = "stage is paused"
		_ = s.units.UpdateEnvironmentStage(ctx, stage.ID, map[string]any{
			"status":            domain.EnvironmentStageStatusPaused,
			"last_completed_at": time.Now(),
		})
		return res, nil
	}
	if shouldSimulateStageFailure(payload, stage.Name) {
		res.Status = string(domain.EnvironmentStageStatusFailed)
		res.FailedCount = res.TargetCount
		res.FailureReason = "simulated stage execution failure"
		_ = s.units.UpdateEnvironmentStage(ctx, stage.ID, map[string]any{
			"status":            domain.EnvironmentStageStatusFailed,
			"last_completed_at": time.Now(),
		})
		return res, nil
	}

	res.SucceededCount = res.TargetCount
	res.FailedCount = 0
	if err := s.units.UpdateEnvironmentStage(ctx, stage.ID, map[string]any{
		"status":            domain.EnvironmentStageStatusSucceeded,
		"last_completed_at": time.Now(),
	}); err != nil {
		return res, err
	}
	return res, nil
}

func (s *StageExecutionService) MarkStageWaiting(ctx context.Context, stageID uint64) error {
	if s == nil || s.units == nil || stageID == 0 {
		return nil
	}
	stage, err := s.units.GetEnvironmentStageByID(ctx, stageID)
	if err != nil {
		return err
	}
	if stage.Status != domain.EnvironmentStageStatusIdle {
		return nil
	}
	return s.units.UpdateEnvironmentStage(ctx, stageID, map[string]any{"status": domain.EnvironmentStageStatusWaiting})
}

func (s *StageExecutionService) estimateTargetCount(ctx context.Context, targetGroupID uint64) int {
	if s == nil || s.targets == nil || targetGroupID == 0 {
		return 1
	}
	group, err := s.targets.GetByID(ctx, targetGroupID)
	if err != nil || strings.TrimSpace(group.ClusterRefsJSON) == "" {
		return 1
	}
	refs := make([]uint64, 0)
	if err := json.Unmarshal([]byte(group.ClusterRefsJSON), &refs); err != nil {
		return 1
	}
	if len(refs) == 0 {
		return 1
	}
	return len(refs)
}

func shouldSimulateStageFailure(payload map[string]any, environment string) bool {
	if len(payload) == 0 {
		return false
	}
	if v, ok := payload["simulateFailure"].(bool); ok && v {
		return true
	}
	raw, ok := payload["simulateFailureEnvironments"]
	if !ok {
		return false
	}
	items, ok := raw.([]any)
	if !ok {
		return false
	}
	env := strings.ToLower(strings.TrimSpace(environment))
	for _, item := range items {
		if strings.ToLower(strings.TrimSpace(fmt.Sprint(item))) == env {
			return true
		}
	}
	return false
}
