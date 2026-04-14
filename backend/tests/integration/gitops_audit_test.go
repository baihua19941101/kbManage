package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	"kbmanage/backend/tests/testutil"
)

func TestGitOpsAuditIntegration_WritesVerifyAndRollbackEvents(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-audit-user",
		Password: "GitOpsAudit@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-audit", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	source := &domain.DeliverySource{
		Name:        "gitops-audit-source",
		SourceType:  domain.DeliverySourceTypeGit,
		Endpoint:    "https://git.example.com/gitops-audit.git",
		WorkspaceID: &access.WorkspaceID,
		ProjectID:   &access.ProjectID,
		Status:      domain.DeliverySourceStatusReady,
	}
	if err := app.DB.WithContext(context.Background()).Create(source).Error; err != nil {
		t.Fatalf("seed source failed: %v", err)
	}

	unit := &domain.ApplicationDeliveryUnit{
		Name:           "gitops-audit-unit",
		WorkspaceID:    access.WorkspaceID,
		ProjectID:      &access.ProjectID,
		SourceID:       source.ID,
		SourcePath:     "apps/demo",
		SyncMode:       domain.DeliverySyncModeManual,
		DeliveryStatus: domain.DeliveryUnitStatusReady,
	}
	if err := app.DB.WithContext(context.Background()).Create(unit).Error; err != nil {
		t.Fatalf("seed delivery unit failed: %v", err)
	}

	verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources/"+u64s(source.ID)+"/verify", nil)
	verifyReq.Header.Set("Authorization", "Bearer "+token)
	verifyResp := httptest.NewRecorder()
	app.Router.ServeHTTP(verifyResp, verifyReq)
	if verifyResp.Code != http.StatusAccepted {
		t.Fatalf("expected verify status=202 got=%d body=%s", verifyResp.Code, strings.TrimSpace(verifyResp.Body.String()))
	}

	actionReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/delivery-units/"+u64s(unit.ID)+"/actions", strings.NewReader(`{
		"actionType":"rollback",
		"targetReleaseId": 1,
		"payload": {"reason":"integration-audit"}
	}`))
	actionReq.Header.Set("Authorization", "Bearer "+token)
	actionReq.Header.Set("Content-Type", "application/json")
	actionResp := httptest.NewRecorder()
	app.Router.ServeHTTP(actionResp, actionReq)
	if actionResp.Code != http.StatusAccepted {
		t.Fatalf("expected action status=202 got=%d body=%s", actionResp.Code, strings.TrimSpace(actionResp.Body.String()))
	}

	var events []domain.AuditEvent
	if err := app.DB.WithContext(context.Background()).Where("action LIKE ?", "gitops.%").Find(&events).Error; err != nil {
		t.Fatalf("query gitops audit events failed: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected gitops audit events generated")
	}

	actions := map[string]bool{}
	for _, event := range events {
		actions[event.Action] = true
		if event.ResourceType != auditSvc.GitOpsAuditResourceType {
			t.Fatalf("expected gitops resource type, got=%s", event.ResourceType)
		}
	}
	if !actions[auditSvc.GitOpsAuditActionSourceVerify] {
		t.Fatalf("expected verify event, got=%v", actions)
	}
	if !actions[auditSvc.GitOpsAuditActionRollbackSubmit] {
		t.Fatalf("expected rollback submit event, got=%v", actions)
	}
}
