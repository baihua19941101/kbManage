package worker

import (
	"context"

	backupRestoreSvc "kbmanage/backend/internal/service/backuprestore"
)

type RestoreJobWorker struct {
	svc *backupRestoreSvc.Service
}

func NewRestoreJobWorker(svc *backupRestoreSvc.Service) *RestoreJobWorker {
	return &RestoreJobWorker{svc: svc}
}

func (w *RestoreJobWorker) Execute(ctx context.Context, userID uint64, input backupRestoreSvc.CreateRestoreJobInput) error {
	if w == nil || w.svc == nil {
		return nil
	}
	_, err := w.svc.CreateRestoreJob(ctx, userID, input)
	return err
}
