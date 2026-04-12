package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_MetricsRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-metrics-user",
		Password:    "ObsMetrics@123",
		DisplayName: "Obs Metrics User",
		Email:       "obs-metrics-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-metrics-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	resp := performObservabilityAuthedRequest(t, app.Router, token, http.MethodGet, "/api/v1/observability/metrics/series?subjectType=workload&subjectRef=mock-app&metricKey=cpu_usage", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
