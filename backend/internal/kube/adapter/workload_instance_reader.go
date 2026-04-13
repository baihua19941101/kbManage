package adapter

import "context"

type WorkloadInstanceReader interface {
	ListByWorkload(ctx context.Context, clusterID uint64, namespace, resourceKind, resourceName string) ([]map[string]any, error)
}

type NoopWorkloadInstanceReader struct{}

func (NoopWorkloadInstanceReader) ListByWorkload(_ context.Context, _ uint64, _ string, _ string, _ string) ([]map[string]any, error) {
	return []map[string]any{}, nil
}
