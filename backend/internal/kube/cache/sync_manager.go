package cache

import (
	"context"
	"fmt"
	"time"

	kcache "k8s.io/client-go/tools/cache"
)

type SyncManager struct {
	timeout time.Duration
}

func NewSyncManager(timeout time.Duration) *SyncManager {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &SyncManager{timeout: timeout}
}

func (m *SyncManager) WaitForSync(ctx context.Context, cacheSyncs ...kcache.InformerSynced) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if ok := kcache.WaitForCacheSync(ctxWithTimeout.Done(), cacheSyncs...); !ok {
		return fmt.Errorf("informer cache sync failed")
	}
	return nil
}
