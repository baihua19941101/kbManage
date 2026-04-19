package worker

import (
	"context"

	backupRestoreSvc "kbmanage/backend/internal/service/backuprestore"
)

type BackupRunWorker struct {
	svc *backupRestoreSvc.Service
}

func NewBackupRunWorker(svc *backupRestoreSvc.Service) *BackupRunWorker {
	return &BackupRunWorker{svc: svc}
}

func (w *BackupRunWorker) Execute(ctx context.Context, userID, policyID uint64) error {
	if w == nil || w.svc == nil {
		return nil
	}
	_, err := w.svc.RunPolicy(ctx, userID, policyID)
	return err
}
