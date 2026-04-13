package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestWorkloadActionExecutionIntegration_ScaleRestart(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-action-exec-integration",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-action-exec-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	for _, actionType := range []string{"scale", "restart"} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/actions", strings.NewReader(`{
			"clusterId": `+strconv.FormatUint(access.ClusterID, 10)+`,
			"namespace": "default",
			"resourceKind": "Deployment",
			"resourceName": "demo-api",
			"actionType": "`+actionType+`",
			"riskConfirmed": true
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)
		if resp.Code != http.StatusAccepted {
			t.Fatalf("actionType=%s expected status=202 got=%d body=%s", actionType, resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		var payload map[string]any
		_ = json.Unmarshal(resp.Body.Bytes(), &payload)
		if strings.TrimSpace(payload["status"].(string)) != "succeeded" {
			t.Fatalf("actionType=%s expected succeeded payload=%v", actionType, payload)
		}
	}
}
