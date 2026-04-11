package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"kbmanage/backend/internal/kube/adapter"
)

const (
	defaultClusterSyncQueueSize = 64
	defaultClusterSyncTimeout   = 20 * time.Second
)

// ClusterSyncWorker serializes cluster snapshot indexing jobs in background.
type ClusterSyncWorker struct {
	indexer     adapter.ResourceIndexer
	queue       chan uint64
	syncTimeout time.Duration

	startOnce sync.Once
}

func NewClusterSyncWorker(indexer adapter.ResourceIndexer, queueSize int, syncTimeout time.Duration) *ClusterSyncWorker {
	if queueSize <= 0 {
		queueSize = defaultClusterSyncQueueSize
	}
	if syncTimeout <= 0 {
		syncTimeout = defaultClusterSyncTimeout
	}
	return &ClusterSyncWorker{
		indexer:     indexer,
		queue:       make(chan uint64, queueSize),
		syncTimeout: syncTimeout,
	}
}

func (w *ClusterSyncWorker) Start(ctx context.Context) {
	if w == nil || w.indexer == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	w.startOnce.Do(func() {
		go w.run(ctx)
	})
}

func (w *ClusterSyncWorker) Enqueue(clusterID uint64) error {
	if w == nil || w.indexer == nil {
		return errors.New("cluster sync worker is not configured")
	}
	if clusterID == 0 {
		return errors.New("clusterID is required")
	}

	select {
	case w.queue <- clusterID:
		return nil
	default:
		return errors.New("cluster sync queue is full")
	}
}

func (w *ClusterSyncWorker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case clusterID := <-w.queue:
			syncCtx, cancel := context.WithTimeout(ctx, w.syncTimeout)
			err := w.indexer.SyncCluster(syncCtx, clusterID)
			cancel()
			if err != nil {
				log.Printf("cluster sync failed for cluster=%d: %v", clusterID, err)
			}
		}
	}
}
