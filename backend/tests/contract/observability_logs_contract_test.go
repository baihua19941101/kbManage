package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_LogsRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-logs-user",
		Password:    "ObsLogs@123",
		DisplayName: "Obs Logs User",
		Email:       "obs-logs-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-logs-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	resp := performObservabilityAuthedRequest(t, app.Router, token, http.MethodGet, "/api/v1/observability/logs/query?namespace=default&workload=mock-app&limit=10", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
