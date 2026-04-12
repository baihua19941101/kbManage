package worker

import (
	"context"
	"sync"
	"time"
)

const defaultObservabilitySyncInterval = 30 * time.Second

type ObservabilitySyncRunner interface {
	Sync(ctx context.Context) error
}

type ObservabilitySyncWorker struct {
	runner   ObservabilitySyncRunner
	interval time.Duration

	startOnce sync.Once
}

func NewObservabilitySyncWorker(runner ObservabilitySyncRunner, interval time.Duration) *ObservabilitySyncWorker {
	if interval <= 0 {
		interval = defaultObservabilitySyncInterval
	}
	return &ObservabilitySyncWorker{
		runner:   runner,
		interval: interval,
	}
}

func (w *ObservabilitySyncWorker) Start(ctx context.Context) {
	if w == nil || w.runner == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *ObservabilitySyncWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.runner.Sync(ctx)
		}
	}
}
