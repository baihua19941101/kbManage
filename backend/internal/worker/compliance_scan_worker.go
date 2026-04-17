package worker

import (
	"context"
	"sync"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

type ComplianceScanWorker struct {
	svc       *complianceSvc.ScanExecutionService
	interval  time.Duration
	startOnce sync.Once
}

func NewComplianceScanWorker(svc *complianceSvc.ScanExecutionService, intervalSeconds int) *ComplianceScanWorker {
	interval := 30 * time.Second
	if intervalSeconds > 0 {
		interval = time.Duration(intervalSeconds) * time.Second
	}
	return &ComplianceScanWorker{svc: svc, interval: interval}
}

func (w *ComplianceScanWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil {
		return
	}
	w.startOnce.Do(func() { go w.run(ctx) })
}

func (w *ComplianceScanWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = w.svc.RunPending(ctx, 10)
		}
	}
}
