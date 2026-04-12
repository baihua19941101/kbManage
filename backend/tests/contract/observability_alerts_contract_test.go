package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_AlertsRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-alerts-contract-user",
		Password:    "ObsAlerts@123",
		DisplayName: "Obs Alerts Contract User",
		Email:       "obs-alerts-contract-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-alerts-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	seed := &domain.AlertIncidentSnapshot{
		SourceIncidentKey: "inc-contract-001",
		Severity:          domain.AlertSeverityWarning,
		Status:            domain.AlertIncidentStatusFiring,
		Summary:           "contract incident",
	}
	if err := app.DB.Create(seed).Error; err != nil {
		t.Fatalf("seed alert incident failed: %v", err)
	}

	resp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodGet,
		"/api/v1/observability/alerts",
		"",
	)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}

	payload := mustDecodeObservabilityObject(t, resp.Body.Bytes())
	assertObservabilityArrayField(t, payload, "items")
}
