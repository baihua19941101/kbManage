package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	"kbmanage/backend/tests/testutil"
)

func TestAlertGovernanceFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-alert-governance-int-user",
		Password:    "ObsAlertGovernance@123",
		DisplayName: "Obs Alert Governance Int User",
		Email:       "obs-alert-governance-int-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-alert-governance-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	createRuleReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observability/alert-rules",
		strings.NewReader(`{"name":"mem high","severity":"critical","conditionExpression":"memory_usage > 90","scopeSnapshot":"{\"workspaceIds\":[`+strconv.FormatUint(seed.WorkspaceID, 10)+`]}"} `),
	)
	createRuleReq.Header.Set("Authorization", "Bearer "+token)
	createRuleReq.Header.Set("Content-Type", "application/json")
	createRuleResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createRuleResp, createRuleReq)
	if createRuleResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createRuleResp.Code, strings.TrimSpace(createRuleResp.Body.String()))
	}

	createTargetReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observability/notification-targets",
		strings.NewReader(`{"name":"ops webhook","targetType":"webhook","configRef":"secret://ops/webhook"}`),
	)
	createTargetReq.Header.Set("Authorization", "Bearer "+token)
	createTargetReq.Header.Set("Content-Type", "application/json")
	createTargetResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createTargetResp, createTargetReq)
	if createTargetResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createTargetResp.Code, strings.TrimSpace(createTargetResp.Body.String()))
	}

	now := time.Now().UTC()
	createSilenceReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observability/silences",
		strings.NewReader(`{"name":"release","startsAt":"`+now.Add(-1*time.Minute).Format(time.RFC3339)+`","endsAt":"`+now.Add(20*time.Minute).Format(time.RFC3339)+`"}`),
	)
	createSilenceReq.Header.Set("Authorization", "Bearer "+token)
	createSilenceReq.Header.Set("Content-Type", "application/json")
	createSilenceResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createSilenceResp, createSilenceReq)
	if createSilenceResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createSilenceResp.Code, strings.TrimSpace(createSilenceResp.Body.String()))
	}

	listRuleReq := httptest.NewRequest(http.MethodGet, "/api/v1/observability/alert-rules", nil)
	listRuleReq.Header.Set("Authorization", "Bearer "+token)
	listRuleResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listRuleResp, listRuleReq)
	if listRuleResp.Code != http.StatusOK || !strings.Contains(listRuleResp.Body.String(), "mem high") {
		t.Fatalf("unexpected alert rule list response status=%d body=%s", listRuleResp.Code, strings.TrimSpace(listRuleResp.Body.String()))
	}

	listTargetReq := httptest.NewRequest(http.MethodGet, "/api/v1/observability/notification-targets", nil)
	listTargetReq.Header.Set("Authorization", "Bearer "+token)
	listTargetResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listTargetResp, listTargetReq)
	if listTargetResp.Code != http.StatusOK || !strings.Contains(listTargetResp.Body.String(), "ops webhook") {
		t.Fatalf("unexpected target list response status=%d body=%s", listTargetResp.Code, strings.TrimSpace(listTargetResp.Body.String()))
	}

	listSilenceReq := httptest.NewRequest(http.MethodGet, "/api/v1/observability/silences", nil)
	listSilenceReq.Header.Set("Authorization", "Bearer "+token)
	listSilenceResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listSilenceResp, listSilenceReq)
	if listSilenceResp.Code != http.StatusOK || !strings.Contains(listSilenceResp.Body.String(), "release") {
		t.Fatalf("unexpected silence list response status=%d body=%s", listSilenceResp.Code, strings.TrimSpace(listSilenceResp.Body.String()))
	}

	var auditEvents []domain.AuditEvent
	if err := app.DB.
		Where(
			"action IN ?",
			[]string{
				auditSvc.ObservabilityAuditActionAlertRuleCreate,
				auditSvc.ObservabilityAuditActionNotificationTargetCreate,
				auditSvc.ObservabilityAuditActionSilenceCreate,
			},
		).
		Find(&auditEvents).Error; err != nil {
		t.Fatalf("query audit events failed: %v", err)
	}
	if len(auditEvents) < 3 {
		t.Fatalf("expected at least 3 observability governance audit events, got=%d", len(auditEvents))
	}
}
