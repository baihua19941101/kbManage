package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/repository"
)

const (
	defaultAuditRetentionPeriod = 180 * 24 * time.Hour
	defaultAuditRetentionTick   = 24 * time.Hour
)

// AuditRetentionWorker periodically purges audit records older than 180 days.
type AuditRetentionWorker struct {
	repo            *repository.AuditRepository
	retentionPeriod time.Duration
	tickInterval    time.Duration
	now             func() time.Time

	startOnce sync.Once
}

func NewAuditRetentionWorker(repo *repository.AuditRepository) *AuditRetentionWorker {
	return &AuditRetentionWorker{
		repo:            repo,
		retentionPeriod: defaultAuditRetentionPeriod,
		tickInterval:    defaultAuditRetentionTick,
		now:             time.Now,
	}
}

func (w *AuditRetentionWorker) Start(ctx context.Context) {
	if w == nil || w.repo == nil {
		return
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *AuditRetentionWorker) run(ctx context.Context) {
	// Run one cleanup immediately on startup so long-running instances do not wait for the first tick.
	w.cleanup(ctx)

	ticker := time.NewTicker(w.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return
			}
			return
		case <-ticker.C:
			w.cleanup(ctx)
		}
	}
}

func (w *AuditRetentionWorker) cleanup(ctx context.Context) {
	if w == nil || w.repo == nil || w.now == nil {
		return
	}
	cutoff := w.now().Add(-w.retentionPeriod)
	_, _ = w.repo.DeleteOlderThan(ctx, cutoff)
}
