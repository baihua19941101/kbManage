package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

type ComplianceRecheckWorker struct {
	svc       *complianceSvc.RecheckService
	interval  time.Duration
	batchSize int
	startOnce sync.Once
}

func NewComplianceRecheckWorker(svc *complianceSvc.RecheckService, interval time.Duration, batchSize int) *ComplianceRecheckWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	if batchSize <= 0 {
		batchSize = 20
	}
	return &ComplianceRecheckWorker{svc: svc, interval: interval, batchSize: batchSize}
}

func (w *ComplianceRecheckWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil {
		return
	}
	w.startOnce.Do(func() { go w.run(ctx) })
}

func (w *ComplianceRecheckWorker) RunOnce(ctx context.Context) (int, error) {
	if w == nil || w.svc == nil {
		return 0, errors.New("compliance recheck worker is not configured")
	}
	return w.svc.RunPending(ctx, w.batchSize)
}

func (w *ComplianceRecheckWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = w.RunOnce(ctx)
		}
	}
}
