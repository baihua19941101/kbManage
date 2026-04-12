package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestAlertCenterFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-alert-center-int-user",
		Password:    "ObsAlertCenter@123",
		DisplayName: "Obs Alert Center Int User",
		Email:       "obs-alert-center-int-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-alert-center-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	incident := &domain.AlertIncidentSnapshot{
		SourceIncidentKey: "inc-alert-center-001",
		ClusterID:         &seed.ClusterID,
		WorkspaceID:       &seed.WorkspaceID,
		ProjectID:         &seed.ProjectID,
		Severity:          domain.AlertSeverityWarning,
		Status:            domain.AlertIncidentStatusFiring,
		Summary:           "integration incident",
	}
	if err := app.DB.Create(incident).Error; err != nil {
		t.Fatalf("seed incident failed: %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/observability/alerts", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if !strings.Contains(listResp.Body.String(), "integration incident") {
		t.Fatalf("expected incident in list, body=%s", strings.TrimSpace(listResp.Body.String()))
	}

	ackReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/observability/alerts/1/acknowledge",
		strings.NewReader(`{"note":"ack from integration"}`),
	)
	ackReq.Header.Set("Authorization", "Bearer "+token)
	ackReq.Header.Set("Content-Type", "application/json")
	ackResp := httptest.NewRecorder()
	app.Router.ServeHTTP(ackResp, ackReq)
	if ackResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", ackResp.Code, strings.TrimSpace(ackResp.Body.String()))
	}
	if !strings.Contains(ackResp.Body.String(), "acknowledged") {
		t.Fatalf("expected acknowledged state, body=%s", strings.TrimSpace(ackResp.Body.String()))
	}
}
