package worker

import (
	"context"
	"time"

	clusterLifecycleSvc "kbmanage/backend/internal/service/clusterlifecycle"
)

type ClusterLifecycleUpgradeWorker struct {
	svc      *clusterLifecycleSvc.Service
	interval time.Duration
}

func NewClusterLifecycleUpgradeWorker(svc *clusterLifecycleSvc.Service, interval time.Duration) *ClusterLifecycleUpgradeWorker {
	return &ClusterLifecycleUpgradeWorker{svc: svc, interval: interval}
}

func (w *ClusterLifecycleUpgradeWorker) Start(_ context.Context) {}
