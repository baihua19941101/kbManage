package worker

import (
	"context"

	backupRestoreSvc "kbmanage/backend/internal/service/backuprestore"
)

type MigrationJobWorker struct {
	svc *backupRestoreSvc.Service
}

func NewMigrationJobWorker(svc *backupRestoreSvc.Service) *MigrationJobWorker {
	return &MigrationJobWorker{svc: svc}
}

func (w *MigrationJobWorker) Execute(ctx context.Context, userID uint64, input backupRestoreSvc.CreateMigrationPlanInput) error {
	if w == nil || w.svc == nil {
		return nil
	}
	_, err := w.svc.CreateMigrationPlan(ctx, userID, input)
	return err
}
