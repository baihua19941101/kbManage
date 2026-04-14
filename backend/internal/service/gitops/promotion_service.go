package gitops

import (
	"context"
	"fmt"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type PromotionService struct {
	units    *repository.DeliveryUnitRepository
	stageOps *StageExecutionService
}

func NewPromotionService(
	units *repository.DeliveryUnitRepository,
	stageOps *StageExecutionService,
) *PromotionService {
	return &PromotionService{units: units, stageOps: stageOps}
}

func (s *PromotionService) Promote(
	ctx context.Context,
	unitID uint64,
	environment string,
	payload map[string]any,
) ([]StageExecutionResult, error) {
	if s == nil || s.units == nil || s.stageOps == nil {
		return nil, ErrGitOpsNotConfigured
	}
	stages, err := s.units.ListEnvironmentStages(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if len(stages) == 0 {
		return nil, fmt.Errorf("delivery unit has no environment stages")
	}

	targetIdx := -1
	env := strings.ToLower(strings.TrimSpace(environment))
	if env != "" {
		for i := range stages {
			if strings.ToLower(strings.TrimSpace(stages[i].Name)) == env {
				targetIdx = i
				break
			}
		}
		if targetIdx < 0 {
			return nil, fmt.Errorf("environment %q not found", environment)
		}
		for i := 0; i < targetIdx; i++ {
			if stages[i].Status != domain.EnvironmentStageStatusSucceeded {
				return nil, fmt.Errorf("previous stage %q is not completed", stages[i].Name)
			}
		}
	} else {
		for i := range stages {
			if stages[i].Status != domain.EnvironmentStageStatusSucceeded {
				targetIdx = i
				break
			}
		}
		if targetIdx < 0 {
			return []StageExecutionResult{}, nil
		}
	}

	result, err := s.stageOps.ExecuteStage(ctx, stages[targetIdx], payload)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(result.Status, string(domain.EnvironmentStageStatusSucceeded)) && targetIdx+1 < len(stages) {
		_ = s.stageOps.MarkStageWaiting(ctx, stages[targetIdx+1].ID)
	}
	return []StageExecutionResult{result}, nil
}

func (s *PromotionService) RollbackStages(
	ctx context.Context,
	unitID uint64,
	environment string,
	payload map[string]any,
) ([]StageExecutionResult, error) {
	if s == nil || s.units == nil || s.stageOps == nil {
		return nil, ErrGitOpsNotConfigured
	}
	stages, err := s.units.ListEnvironmentStages(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if len(stages) == 0 {
		return []StageExecutionResult{}, nil
	}

	limit := len(stages)
	env := strings.ToLower(strings.TrimSpace(environment))
	if env != "" {
		limit = 0
		for i := range stages {
			if strings.ToLower(strings.TrimSpace(stages[i].Name)) == env {
				limit = i + 1
				break
			}
		}
		if limit == 0 {
			return nil, fmt.Errorf("environment %q not found", environment)
		}
	}

	results := make([]StageExecutionResult, 0, limit)
	for i := 0; i < limit; i++ {
		result, err := s.stageOps.ExecuteStage(ctx, stages[i], payload)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
		if strings.EqualFold(result.Status, string(domain.EnvironmentStageStatusFailed)) {
			break
		}
	}
	return results, nil
}
