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

func TestWorkloadOpsContextIntegration_ContextAndInstances(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "workloadops-context-integration",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-context-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	t.Run("get context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/context?clusterId="+strconv.FormatUint(access.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=api-server", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Fatalf("expected context status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		var payload map[string]any
		if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode context failed: %v", err)
		}
		actions, ok := payload["availableActions"].([]any)
		if !ok || len(actions) == 0 {
			t.Fatalf("expected availableActions not empty, payload=%v", payload)
		}
	})

	t.Run("get instances", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/instances?clusterId="+strconv.FormatUint(access.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=api-server", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Fatalf("expected instances status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		var payload struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode instances failed: %v", err)
		}
		if len(payload.Items) == 0 {
			t.Fatalf("expected non-empty instances")
		}
	})
}
