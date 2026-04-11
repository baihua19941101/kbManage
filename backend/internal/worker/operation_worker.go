package worker

import (
	"context"
	"errors"
	"strings"
	"sync"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	operationSvc "kbmanage/backend/internal/service/operation"
)

// OperationWorker consumes operation queue and updates operation status.
type OperationWorker struct {
	repo       *repository.OperationRepository
	queue      operationSvc.QueueService
	executor   operationSvc.Executor
	auditWrite *auditSvc.EventWriter

	startOnce sync.Once
}

func NewOperationWorker(
	repo *repository.OperationRepository,
	queue operationSvc.QueueService,
	executor operationSvc.Executor,
	auditWriter *auditSvc.EventWriter,
) *OperationWorker {
	if executor == nil {
		executor = operationSvc.NewExecutor(nil)
	}
	return &OperationWorker{
		repo:       repo,
		queue:      queue,
		executor:   executor,
		auditWrite: auditWriter,
	}
}

func (w *OperationWorker) Start(ctx context.Context) {
	if w == nil || w.repo == nil || w.queue == nil {
		return
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *OperationWorker) run(ctx context.Context) {
	for {
		operationID, err := w.queue.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			continue
		}

		item, err := w.repo.GetByID(ctx, operationID)
		if err != nil {
			continue
		}

		runningItem, transitioned, err := w.repo.TransitionStatus(
			ctx,
			operationID,
			[]domain.OperationStatus{domain.OperationStatusPending},
			domain.OperationStatusRunning,
			"operation is running",
			"",
			"",
		)
		if err != nil {
			continue
		}
		if !transitioned {
			continue
		}
		if runningItem != nil {
			item = runningItem
		}

		_ = w.writeAuditEvent(ctx, item, auditSvc.OperationAuditActionStart, domain.AuditOutcomeSuccess, map[string]any{
			"operationType": item.OperationType,
			"targetRef":     item.TargetRef,
			"status":        item.Status,
		})

		result, execErr := w.executor.Execute(ctx, item)
		if execErr != nil {
			failureReason := strings.TrimSpace(result.FailureReason)
			if failureReason == "" {
				failureReason = strings.TrimSpace(execErr.Error())
			}
			progressMessage := strings.TrimSpace(result.ProgressMessage)
			if progressMessage == "" {
				progressMessage = "operation execution failed"
			}
			resultMessage := strings.TrimSpace(result.ResultMessage)
			if resultMessage == "" {
				resultMessage = failureReason
			}

			failedItem, updated, transitionErr := w.repo.TransitionStatus(
				ctx,
				operationID,
				[]domain.OperationStatus{domain.OperationStatusRunning},
				domain.OperationStatusFailed,
				progressMessage,
				resultMessage,
				failureReason,
			)
			if transitionErr == nil && updated {
				_ = w.writeAuditEvent(ctx, failedItem, auditSvc.OperationAuditActionFailure, domain.AuditOutcomeFailed, map[string]any{
					"operationType":  item.OperationType,
					"targetRef":      item.TargetRef,
					"status":         domain.OperationStatusFailed,
					"failureReason":  failureReason,
					"executorResult": resultMessage,
				})
			}
			continue
		}

		progressMessage := strings.TrimSpace(result.ProgressMessage)
		if progressMessage == "" {
			progressMessage = "operation execution completed"
		}
		resultMessage := strings.TrimSpace(result.ResultMessage)
		if resultMessage == "" {
			resultMessage = buildDefaultOperationResultMessage(item)
		}

		succeededItem, updated, transitionErr := w.repo.TransitionStatus(
			ctx,
			operationID,
			[]domain.OperationStatus{domain.OperationStatusRunning},
			domain.OperationStatusSucceeded,
			progressMessage,
			resultMessage,
			"",
		)
		if transitionErr == nil && updated {
			_ = w.writeAuditEvent(ctx, succeededItem, auditSvc.OperationAuditActionSuccess, domain.AuditOutcomeSuccess, map[string]any{
				"operationType": item.OperationType,
				"targetRef":     item.TargetRef,
				"status":        domain.OperationStatusSucceeded,
				"resultMessage": resultMessage,
			})
		}
	}
}

func buildDefaultOperationResultMessage(item *domain.OperationRequest) string {
	if item == nil {
		return "operation executed successfully"
	}
	if strings.TrimSpace(item.TargetRef) != "" {
		return "operation executed: " + strings.TrimSpace(item.TargetRef)
	}
	return "operation executed successfully"
}

func (w *OperationWorker) writeAuditEvent(
	ctx context.Context,
	item *domain.OperationRequest,
	action string,
	outcome domain.AuditOutcome,
	details map[string]any,
) error {
	if w == nil || w.auditWrite == nil || item == nil {
		return nil
	}
	actorID := item.OperatorID
	if details == nil {
		details = map[string]any{}
	}
	details["operationId"] = item.ID
	details["requestId"] = item.RequestID
	details["operationType"] = item.OperationType
	return w.auditWrite.WriteOperationEvent(
		ctx,
		item.RequestID,
		&actorID,
		item.ID,
		action,
		outcome,
		details,
	)
}
