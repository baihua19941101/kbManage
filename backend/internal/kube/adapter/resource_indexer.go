package adapter

import "context"

// ResourceIndexer is a stub for background Kubernetes inventory indexing.
type ResourceIndexer interface {
	SyncCluster(ctx context.Context, clusterID uint64) error
}

type NoopResourceIndexer struct{}

func NewResourceIndexer() ResourceIndexer {
	return NoopResourceIndexer{}
}

func (NoopResourceIndexer) SyncCluster(_ context.Context, _ uint64) error {
	return nil
}
