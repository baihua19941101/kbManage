package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_AccessControlDeniedWithoutScopeBinding(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-access-denied-contract-user",
		Password:    "ObsAccessDenied@123",
		DisplayName: "Obs Access Denied Contract User",
		Email:       "obs-access-denied-contract-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	cases := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{name: "logs", method: http.MethodGet, target: "/api/v1/observability/logs/query?clusterId=1&namespace=default"},
		{name: "events", method: http.MethodGet, target: "/api/v1/observability/events?clusterId=1&namespace=default&resourceKind=Pod&resourceName=demo"},
		{name: "metrics", method: http.MethodGet, target: "/api/v1/observability/metrics/series?clusterId=1&subjectType=pod&subjectRef=demo&metricKey=cpu_usage"},
		{name: "alerts", method: http.MethodGet, target: "/api/v1/observability/alerts"},
		{name: "governance", method: http.MethodPost, target: "/api/v1/observability/alert-rules", body: `{"name":"blocked","severity":"warning","conditionExpression":"cpu > 90"}`},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := performObservabilityAuthedRequest(t, app.Router, token, tc.method, tc.target, tc.body)
			if resp.Code != http.StatusForbidden {
				t.Fatalf("expected status=403, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
		})
	}
}
