package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

type ComplianceTrendSnapshotWorker struct {
	svc       *complianceSvc.TrendService
	interval  time.Duration
	startOnce sync.Once
}

func NewComplianceTrendSnapshotWorker(svc *complianceSvc.TrendService, interval time.Duration) *ComplianceTrendSnapshotWorker {
	if interval <= 0 {
		interval = 15 * time.Minute
	}
	return &ComplianceTrendSnapshotWorker{svc: svc, interval: interval}
}

func (w *ComplianceTrendSnapshotWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil {
		return
	}
	w.startOnce.Do(func() { go w.run(ctx) })
}

func (w *ComplianceTrendSnapshotWorker) RunOnce(ctx context.Context, filter complianceSvc.TrendFilter) (*complianceSvc.ComplianceTrendPoint, error) {
	if w == nil || w.svc == nil {
		return nil, errors.New("compliance trend snapshot worker is not configured")
	}
	return w.svc.RecordSnapshot(ctx, filter)
}

func (w *ComplianceTrendSnapshotWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = w.RunOnce(ctx, complianceSvc.TrendFilter{})
		}
	}
}
