package worker

import (
	"context"
	"errors"
	"sync"

	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
)

// AuditExportWorker consumes audit export tasks and updates task status.
type AuditExportWorker struct {
	svc  *auditSvc.Service
	repo *repository.AuditExportRepository

	startOnce sync.Once
}

func NewAuditExportWorker(svc *auditSvc.Service, repo *repository.AuditExportRepository) *AuditExportWorker {
	return &AuditExportWorker{svc: svc, repo: repo}
}

func (w *AuditExportWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil || w.repo == nil {
		return
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *AuditExportWorker) run(ctx context.Context) {
	for {
		taskID, err := w.repo.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			continue
		}
		_ = w.svc.ProcessExportTask(ctx, taskID)
	}
}
