package integration_test

import (
	"context"
	"testing"

	"kbmanage/backend/internal/domain"
	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	obsSvc "kbmanage/backend/internal/service/observability"
	"kbmanage/backend/tests/testutil"
)

func TestObservabilityAlertSyncWritesSnapshotsAndAudit(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	incidentRepo := repository.NewAlertIncidentRepository(app.DB)
	auditRepo := repository.NewAuditRepository(app.DB)

	syncSvc := obsSvc.NewAlertSyncService(alertProvider.NewMockProvider(), incidentRepo)
	syncSvc.SetAuditWriter(auditSvc.NewEventWriter(auditRepo))

	if err := syncSvc.Sync(context.Background()); err != nil {
		t.Fatalf("sync alerts failed: %v", err)
	}

	var incidents []domain.AlertIncidentSnapshot
	if err := app.DB.Find(&incidents).Error; err != nil {
		t.Fatalf("query incidents failed: %v", err)
	}
	if len(incidents) == 0 {
		t.Fatalf("expected synced incident snapshots")
	}

	var auditEvents []domain.AuditEvent
	if err := app.DB.
		Where("action = ?", auditSvc.ObservabilityAuditActionAlertSync).
		Find(&auditEvents).Error; err != nil {
		t.Fatalf("query audit events failed: %v", err)
	}
	if len(auditEvents) == 0 {
		t.Fatalf("expected observability alert sync audit event")
	}
}
