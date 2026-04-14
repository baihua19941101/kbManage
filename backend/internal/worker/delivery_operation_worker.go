package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	gitopsSvc "kbmanage/backend/internal/service/gitops"
)

// DeliveryOperationWorker consumes gitops operation queue and executes actions.
type DeliveryOperationWorker struct {
	operations *repository.DeliveryOperationRepository
	units      *repository.DeliveryUnitRepository
	queue      gitopsSvc.OperationQueue
	executor   *gitopsSvc.Executor
	progress   *gitopsSvc.ProgressCache

	startOnce sync.Once
}

func NewDeliveryOperationWorker(
	operations *repository.DeliveryOperationRepository,
	units *repository.DeliveryUnitRepository,
	queue gitopsSvc.OperationQueue,
	executor *gitopsSvc.Executor,
	progress *gitopsSvc.ProgressCache,
) *DeliveryOperationWorker {
	return &DeliveryOperationWorker{
		operations: operations,
		units:      units,
		queue:      queue,
		executor:   executor,
		progress:   progress,
	}
}

func (w *DeliveryOperationWorker) Start(ctx context.Context) {
	if w == nil || w.operations == nil || w.units == nil || w.queue == nil || w.executor == nil {
		return
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *DeliveryOperationWorker) run(ctx context.Context) {
	for {
		operationID, err := w.queue.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			continue
		}
		if operationID == 0 {
			continue
		}

		item, err := w.operations.GetByID(ctx, operationID)
		if err != nil || item == nil {
			continue
		}
		if item.Status != domain.DeliveryOperationStatusPending {
			continue
		}

		_ = w.operations.UpdateStatus(
			ctx,
			item.ID,
			domain.DeliveryOperationStatusRunning,
			15,
			"operation is running",
			"",
		)
		w.setProgress(ctx, item.ID, 15, "running", "")

		item, err = w.operations.GetByID(ctx, operationID)
		if err != nil || item == nil {
			continue
		}

		payload, err := decodeOperationPayload(item.PayloadJSON)
		if err != nil {
			w.markFailed(ctx, item, "invalid payload", err)
			continue
		}

		unit, err := w.units.GetByID(ctx, item.DeliveryUnitID)
		if err != nil {
			w.markFailed(ctx, item, "load delivery unit failed", err)
			continue
		}
		detail, err := w.units.GetDetailByID(ctx, item.DeliveryUnitID)
		if err != nil {
			w.markFailed(ctx, item, "load delivery unit detail failed", err)
			continue
		}

		result, execErr := w.executor.Execute(ctx, gitopsSvc.ExecuteInput{
			Operation: item,
			Unit:      unit,
			Detail:    detail,
			Payload:   payload,
		})
		if execErr != nil {
			message := strings.TrimSpace(result.ResultSummary)
			if message == "" {
				message = "delivery action execution failed"
			}
			w.markFailed(ctx, item, message, errors.Join(execErr, errors.New(strings.TrimSpace(result.FailureReason))))
			continue
		}

		if err := w.applyResult(ctx, item, payload, result); err != nil {
			w.markFailed(ctx, item, "persist operation result failed", err)
			continue
		}
	}
}

func (w *DeliveryOperationWorker) applyResult(
	ctx context.Context,
	item *domain.DeliveryOperation,
	payload map[string]any,
	result gitopsSvc.ExecuteResult,
) error {
	if item == nil {
		return errors.New("operation is required")
	}

	if len(result.UnitUpdates) > 0 {
		if err := w.units.UpdateFields(ctx, item.DeliveryUnitID, result.UnitUpdates); err != nil {
			return err
		}
	}

	mergedPayload := mergePayload(payload, result.OperationPatch)
	encodedPayload, err := encodePayload(mergedPayload)
	if err != nil {
		return err
	}
	if encodedPayload != strings.TrimSpace(item.PayloadJSON) {
		if err := w.operations.UpdatePayload(ctx, item.ID, encodedPayload); err != nil {
			return err
		}
	}

	status := normalizeDeliveryOperationStatus(result.Status)
	progress := normalizeProgress(result.Progress)
	summary := strings.TrimSpace(result.ResultSummary)
	if summary == "" {
		summary = fmt.Sprintf("%s completed", strings.TrimSpace(string(item.ActionType)))
	}
	reason := strings.TrimSpace(result.FailureReason)
	if status == domain.DeliveryOperationStatusFailed && reason == "" {
		reason = "delivery action failed"
	}

	if err := w.operations.UpdateStatus(ctx, item.ID, status, progress, summary, reason); err != nil {
		return err
	}
	w.setProgress(ctx, item.ID, progress, summary, reason)
	return nil
}

func (w *DeliveryOperationWorker) markFailed(ctx context.Context, item *domain.DeliveryOperation, summary string, execErr error) {
	if item == nil {
		return
	}
	reason := "operation execution failed"
	if execErr != nil {
		reason = strings.TrimSpace(execErr.Error())
	}
	summary = strings.TrimSpace(summary)
	if summary == "" {
		summary = "delivery action failed"
	}
	_ = w.operations.UpdateStatus(
		ctx,
		item.ID,
		domain.DeliveryOperationStatusFailed,
		normalizeProgress(item.ProgressPercent),
		summary,
		reason,
	)
	w.setProgress(ctx, item.ID, normalizeProgress(item.ProgressPercent), summary, reason)
}

func (w *DeliveryOperationWorker) setProgress(ctx context.Context, operationID uint64, percent int, message string, reason string) {
	if w == nil || w.progress == nil || operationID == 0 {
		return
	}
	_ = w.progress.SetOperationProgress(ctx, operationID, gitopsSvc.OperationProgressSnapshot{
		Percent:   normalizeProgress(percent),
		Message:   strings.TrimSpace(composeProgressMessage(message, reason)),
		UpdatedAt: time.Now(),
	})
}

func composeProgressMessage(message string, reason string) string {
	base := strings.TrimSpace(message)
	detail := strings.TrimSpace(reason)
	if base == "" {
		return detail
	}
	if detail == "" {
		return base
	}
	return base + " | " + detail
}

func normalizeDeliveryOperationStatus(status domain.DeliveryOperationStatus) domain.DeliveryOperationStatus {
	switch status {
	case domain.DeliveryOperationStatusPending,
		domain.DeliveryOperationStatusRunning,
		domain.DeliveryOperationStatusPartiallySucceeded,
		domain.DeliveryOperationStatusSucceeded,
		domain.DeliveryOperationStatusFailed,
		domain.DeliveryOperationStatusCanceled:
		return status
	default:
		return domain.DeliveryOperationStatusSucceeded
	}
}

func normalizeProgress(progress int) int {
	if progress < 0 {
		return 0
	}
	if progress > 100 {
		return 100
	}
	if progress == 0 {
		return 100
	}
	return progress
}

func decodeOperationPayload(raw string) (map[string]any, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return map[string]any{}, nil
	}
	payload := make(map[string]any)
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func encodePayload(payload map[string]any) (string, error) {
	if len(payload) == 0 {
		return "", nil
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func mergePayload(base map[string]any, patch map[string]any) map[string]any {
	if len(base) == 0 && len(patch) == 0 {
		return map[string]any{}
	}
	merged := make(map[string]any, len(base)+len(patch))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range patch {
		merged[k] = v
	}
	return merged
}
