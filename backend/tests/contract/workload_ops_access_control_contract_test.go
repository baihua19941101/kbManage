package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestWorkloadOpsContract_AccessControlDeniedWithoutScopeBinding(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-access-denied-contract-user",
		Password: "WorkloadOps@123",
	})
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	cases := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{name: "context", method: http.MethodGet, target: "/api/v1/workload-ops/resources/context?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=demo"},
		{name: "instances", method: http.MethodGet, target: "/api/v1/workload-ops/resources/instances?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=demo"},
		{name: "revisions", method: http.MethodGet, target: "/api/v1/workload-ops/resources/revisions?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=demo"},
		{name: "actions", method: http.MethodPost, target: "/api/v1/workload-ops/actions", body: `{"clusterId":1,"namespace":"default","resourceKind":"Deployment","resourceName":"demo","actionType":"restart","riskConfirmed":true}`},
		{name: "batches", method: http.MethodPost, target: "/api/v1/workload-ops/batches", body: `{"actionType":"restart","riskConfirmed":true,"targets":[{"clusterId":1,"namespace":"default","resourceKind":"Deployment","resourceName":"demo"}]}`},
		{name: "terminal", method: http.MethodPost, target: "/api/v1/workload-ops/terminal/sessions", body: `{"clusterId":1,"namespace":"default","podName":"demo-pod","containerName":"app"}`},
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
