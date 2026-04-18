package worker

import (
	"context"
	"time"

	clusterLifecycleSvc "kbmanage/backend/internal/service/clusterlifecycle"
)

type ClusterLifecycleRetireWorker struct {
	svc      *clusterLifecycleSvc.Service
	interval time.Duration
}

func NewClusterLifecycleRetireWorker(svc *clusterLifecycleSvc.Service, interval time.Duration) *ClusterLifecycleRetireWorker {
	return &ClusterLifecycleRetireWorker{svc: svc, interval: interval}
}

func (w *ClusterLifecycleRetireWorker) Start(_ context.Context) {}
