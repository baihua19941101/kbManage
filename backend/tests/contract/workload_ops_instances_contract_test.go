package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestWorkloadOpsInstancesContract_ListInstances(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "workloadops-instances-contract",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-instances-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/instances?clusterId="+strconv.FormatUint(access.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=demo-api", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}

	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if len(payload.Items) == 0 {
		t.Fatalf("expected at least one instance item")
	}
	first := payload.Items[0]
	if strings.TrimSpace(stringValue(first["podName"])) == "" {
		t.Fatalf("podName is required: %v", first)
	}
	if strings.TrimSpace(stringValue(first["phase"])) == "" {
		t.Fatalf("phase is required: %v", first)
	}
	if _, ok := first["ready"].(bool); !ok {
		t.Fatalf("ready must be bool: %v", first)
	}
}

func stringValue(v any) string {
	s, _ := v.(string)
	return s
}
