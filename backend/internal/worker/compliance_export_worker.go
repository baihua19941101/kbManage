package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

type ComplianceExportWorker struct {
	svc       *complianceSvc.ArchiveExportService
	interval  time.Duration
	startOnce sync.Once
}

func NewComplianceExportWorker(svc *complianceSvc.ArchiveExportService, interval time.Duration) *ComplianceExportWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	return &ComplianceExportWorker{svc: svc, interval: interval}
}

func (w *ComplianceExportWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil {
		return
	}
	w.startOnce.Do(func() { go w.run(ctx) })
}

func (w *ComplianceExportWorker) RunOnce(ctx context.Context) (int, error) {
	if w == nil || w.svc == nil {
		return 0, errors.New("compliance export worker is not configured")
	}
	items, err := w.svc.ListExports(ctx, complianceSvc.ArchiveExportFilter{Status: "pending"})
	if err != nil {
		return 0, err
	}
	processed := 0
	for _, item := range items {
		if _, err := w.svc.ProcessExport(ctx, item.ID); err == nil {
			processed++
		}
	}
	return processed, nil
}

func (w *ComplianceExportWorker) run(ctx context.Context) {
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
