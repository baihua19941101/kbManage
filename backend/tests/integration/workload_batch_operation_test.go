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

func TestWorkloadBatchOperationIntegration_PartialAndSummary(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-batch-integration",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-batch-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/batches", strings.NewReader(`{
		"actionType": "restart",
		"riskConfirmed": true,
		"targets": [
			{"clusterId":`+strconv.FormatUint(access.ClusterID, 10)+`,"namespace":"default","resourceKind":"Deployment","resourceName":"demo-api-1"},
			{"clusterId":`+strconv.FormatUint(access.ClusterID, 10)+`,"namespace":"default","resourceKind":"Deployment","resourceName":"demo-api-2"}
		]
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
	if payload["totalTargets"].(float64) != 2 {
		t.Fatalf("expected totalTargets=2 payload=%v", payload)
	}
	if payload["status"] == nil {
		t.Fatalf("expected batch status payload=%v", payload)
	}
}
