package workloadops

import (
	"context"

	"kbmanage/backend/internal/domain"
)

func (s *Service) executeBatchTask(
	ctx context.Context,
	operatorID uint64,
	task *domain.BatchOperationTask,
	items []domain.BatchOperationItem,
	targets []WorkloadReference,
	actionType domain.WorkloadActionType,
	riskLevel domain.RiskLevel,
	riskConfirmed bool,
	payloadJSON string,
) error {
	if task == nil || s.batches == nil {
		return nil
	}

	_ = s.batches.UpdateTaskSummary(ctx, task.ID, domain.BatchOperationStatusRunning, 0, 0, 0, 0)

	succeeded := 0
	failed := 0
	for idx := range targets {
		if idx >= len(items) {
			break
		}
		itemID := items[idx].ID
		_ = s.batches.UpdateItemResult(ctx, task.ID, itemID, domain.BatchOperationItemStatusRunning, nil, "", "")

		action, err := s.createAndExecuteAction(ctx, SubmitWorkloadActionRequest{
			RequestID:     normalizeRequestID(task.RequestID, operatorID),
			OperatorID:    operatorID,
			Target:        targets[idx],
			ActionType:    actionType,
			RiskLevel:     riskLevel,
			RiskConfirmed: riskConfirmed,
			PayloadJSON:   payloadJSON,
			BatchID:       &task.ID,
		})
		if err != nil {
			failed++
			_ = s.batches.UpdateItemResult(ctx, task.ID, itemID, domain.BatchOperationItemStatusFailed, nil, "", err.Error())
			continue
		}

		actionID := action.ID
		succeeded++
		_ = s.batches.UpdateItemResult(ctx, task.ID, itemID, domain.BatchOperationItemStatusSucceeded, &actionID, action.ResultMessage, "")
	}

	status := domain.BatchOperationStatusSucceeded
	if failed > 0 && succeeded > 0 {
		status = domain.BatchOperationStatusPartiallySucceeded
	}
	if failed > 0 && succeeded == 0 {
		status = domain.BatchOperationStatusFailed
	}
	progress := 100
	return s.batches.UpdateTaskSummary(ctx, task.ID, status, succeeded, failed, 0, progress)
}
