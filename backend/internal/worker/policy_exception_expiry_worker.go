package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"kbmanage/backend/internal/repository"
	securityPolicySvc "kbmanage/backend/internal/service/securitypolicy"
)

// PolicyExceptionExpiryWorker recycles expired policy exceptions.
type PolicyExceptionExpiryWorker struct {
	exceptions *repository.PolicyExceptionRepository
	cache      *securityPolicySvc.ExceptionCache
	interval   time.Duration

	startOnce sync.Once
}

func NewPolicyExceptionExpiryWorker(
	exceptions *repository.PolicyExceptionRepository,
	cache *securityPolicySvc.ExceptionCache,
	interval time.Duration,
) *PolicyExceptionExpiryWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	return &PolicyExceptionExpiryWorker{
		exceptions: exceptions,
		cache:      cache,
		interval:   interval,
	}
}

func (w *PolicyExceptionExpiryWorker) Start(ctx context.Context) {
	if w == nil || w.exceptions == nil {
		return
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *PolicyExceptionExpiryWorker) RunOnce(ctx context.Context, now time.Time) (int, error) {
	if w == nil || w.exceptions == nil {
		return 0, errors.New("policy exception expiry worker is not configured")
	}
	ids, err := w.exceptions.ExpireActiveBefore(ctx, now)
	if err != nil {
		return 0, err
	}
	for _, id := range ids {
		if w.cache != nil {
			_ = w.cache.SetExceptionStatus(ctx, id, "expired")
		}
	}
	return len(ids), nil
}

func (w *PolicyExceptionExpiryWorker) run(ctx context.Context) {
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
