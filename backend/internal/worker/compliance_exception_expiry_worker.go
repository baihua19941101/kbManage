package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

type ComplianceExceptionExpiryWorker struct {
	svc       *complianceSvc.ExceptionService
	interval  time.Duration
	startOnce sync.Once
}

func NewComplianceExceptionExpiryWorker(svc *complianceSvc.ExceptionService, interval time.Duration) *ComplianceExceptionExpiryWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	return &ComplianceExceptionExpiryWorker{svc: svc, interval: interval}
}

func (w *ComplianceExceptionExpiryWorker) Start(ctx context.Context) {
	if w == nil || w.svc == nil {
		return
	}
	w.startOnce.Do(func() { go w.run(ctx) })
}

func (w *ComplianceExceptionExpiryWorker) RunOnce(ctx context.Context, now time.Time) (int, error) {
	if w == nil || w.svc == nil {
		return 0, errors.New("compliance exception expiry worker is not configured")
	}
	return w.svc.ExpireDueExceptions(ctx, now)
}

func (w *ComplianceExceptionExpiryWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = w.RunOnce(ctx, time.Now())
		}
	}
}
