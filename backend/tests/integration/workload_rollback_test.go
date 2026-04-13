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

func TestWorkloadRollbackIntegration_SubmitRollback(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-rollback-integration",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-rollback-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/actions", strings.NewReader(`{
		"clusterId": `+strconv.FormatUint(access.ClusterID, 10)+`,
		"namespace": "default",
		"resourceKind": "Deployment",
		"resourceName": "demo-api",
		"actionType": "rollback",
		"riskConfirmed": true,
		"payload": {"revision": 2}
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("expected status=202 got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	var payload map[string]any
	_ = json.Unmarshal(resp.Body.Bytes(), &payload)
	if strings.TrimSpace(payload["status"].(string)) != "succeeded" {
		t.Fatalf("expected succeeded rollback payload=%v", payload)
	}
	result := strings.TrimSpace(payload["resultMessage"].(string))
	if !strings.Contains(result, "revision 2") {
		t.Fatalf("expected result contains revision 2, got=%q payload=%v", result, payload)
	}
}
