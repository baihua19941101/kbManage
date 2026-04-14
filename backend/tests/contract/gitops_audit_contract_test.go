package contract_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	"kbmanage/backend/tests/testutil"
)

func TestGitOpsContract_AuditEventsForVerifyAndAction(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-audit-contract-user",
		Password: "GitOpsAudit@123",
	})
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	seed := seedGitOpsScopeFixture(t, app, user.User.ID, "workspace-owner")

	verifyResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodPost,
		"/api/v1/gitops/sources/"+u64s(seed.SourceID)+"/verify",
		"",
	)
	if verifyResp.Code != http.StatusAccepted {
		t.Fatalf("expected verify status=202 got=%d body=%s", verifyResp.Code, strings.TrimSpace(verifyResp.Body.String()))
	}

	actionResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodPost,
		"/api/v1/gitops/delivery-units/"+u64s(seed.UnitID)+"/actions",
		`{"actionType":"sync","payload":{"reason":"audit-contract"}}`,
	)
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
			t.Fatalf("expected resourceType=%s, got=%s", auditSvc.GitOpsAuditResourceType, event.ResourceType)
		}
	}
	if !actions[auditSvc.GitOpsAuditActionSourceVerify] {
		t.Fatalf("expected source verify audit event, got=%v", actions)
	}
	if !actions[auditSvc.GitOpsAuditActionSyncSubmit] {
		t.Fatalf("expected sync submit audit event, got=%v", actions)
	}
}
