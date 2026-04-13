package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestWorkloadOpsAuditIntegration_WritesActionAndTerminalEvents(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-audit-user",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-audit", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	actionReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/actions", strings.NewReader(`{
		"clusterId": `+u64s(access.ClusterID)+`,
		"namespace": "default",
		"resourceKind": "Deployment",
		"resourceName": "demo",
		"actionType": "rollback",
		"riskConfirmed": true,
		"payload": {"revision": 2}
	}`))
	actionReq.Header.Set("Authorization", "Bearer "+token)
	actionReq.Header.Set("Content-Type", "application/json")
	actionResp := httptest.NewRecorder()
	app.Router.ServeHTTP(actionResp, actionReq)
	if actionResp.Code != http.StatusAccepted {
		t.Fatalf("expected action status=202 got=%d body=%s", actionResp.Code, strings.TrimSpace(actionResp.Body.String()))
	}

	sessionReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/terminal/sessions", strings.NewReader(`{
		"clusterId": `+u64s(access.ClusterID)+`,
		"namespace": "default",
		"podName": "demo-pod",
		"containerName": "app",
		"workloadKind": "Deployment",
		"workloadName": "demo"
	}`))
	sessionReq.Header.Set("Authorization", "Bearer "+token)
	sessionReq.Header.Set("Content-Type", "application/json")
	sessionResp := httptest.NewRecorder()
	app.Router.ServeHTTP(sessionResp, sessionReq)
	if sessionResp.Code != http.StatusCreated {
		t.Fatalf("expected session create status=201 got=%d body=%s", sessionResp.Code, strings.TrimSpace(sessionResp.Body.String()))
	}

	var events []domain.AuditEvent
	if err := app.DB.Where("action LIKE ?", "workloadops.%").Find(&events).Error; err != nil {
		t.Fatalf("query audit events failed: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected workloadops audit events generated")
	}

	actions := map[string]bool{}
	for _, event := range events {
		actions[event.Action] = true
	}
	if !actions["workloadops.rollback.submit"] {
		t.Fatalf("expected rollback submit event, got=%v", actions)
	}
	if !actions["workloadops.terminal.open"] {
		t.Fatalf("expected terminal open event, got=%v", actions)
	}
}
