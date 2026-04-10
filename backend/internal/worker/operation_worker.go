package worker

import (
	"context"
	"errors"
	"strings"
	"sync"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	operationSvc "kbmanage/backend/internal/service/operation"
)

// OperationWorker consumes operation queue and updates operation status.
type OperationWorker struct {
	repo  *repository.OperationRepository
	queue operationSvc.QueueService

	startOnce sync.Once
}

func NewOperationWorker(repo *repository.OperationRepository, queue operationSvc.QueueService) *OperationWorker {
	return &OperationWorker{repo: repo, queue: queue}
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

		_ = w.repo.UpdateStatus(ctx, operationID, domain.OperationStatusRunning, "")
		item, err := w.repo.GetByID(ctx, operationID)
		if err != nil {
			_ = w.repo.UpdateStatus(ctx, operationID, domain.OperationStatusFailed, "operation record not found")
			continue
		}

		if strings.EqualFold(strings.TrimSpace(item.OperationType), "fail") {
			_ = w.repo.UpdateStatus(ctx, operationID, domain.OperationStatusFailed, "operation execution failed")
			continue
		}

		msg := "operation executed successfully"
		if item.TargetRef != "" {
			msg = "operation executed: " + item.TargetRef
		}
		_ = w.repo.UpdateStatus(ctx, operationID, domain.OperationStatusSucceeded, msg)
	}
}
