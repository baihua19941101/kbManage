package worker

import "context"

type WorkloadBatchRunner interface {
	Run(ctx context.Context) error
}

type WorkloadBatchWorker struct {
	runner WorkloadBatchRunner
}

func NewWorkloadBatchWorker(runner WorkloadBatchRunner) *WorkloadBatchWorker {
	return &WorkloadBatchWorker{runner: runner}
}

func (w *WorkloadBatchWorker) Start(ctx context.Context) error {
	if w == nil || w.runner == nil {
		return nil
	}
	return w.runner.Run(ctx)
}
