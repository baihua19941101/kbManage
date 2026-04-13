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

func TestWorkloadOpsContextContract_GetContext(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "workloadops-context-contract",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-context-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/context?clusterId="+strconv.FormatUint(access.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=demo-api", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeWopsObject(t, resp.Body.Bytes())
	assertWopsStringField(t, payload, "namespace", "default")
	assertWopsStringField(t, payload, "resourceKind", "Deployment")
	assertWopsStringField(t, payload, "resourceName", "demo-api")
	assertWopsStringField(t, payload, "healthStatus", "unknown")
	assertWopsStringField(t, payload, "rolloutStatus", "unknown")
}

func mustDecodeWopsObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid object: %v", err)
	}
	return payload
}

func assertWopsStringField(t *testing.T, payload map[string]any, key, expected string) {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	got, _ := raw.(string)
	if strings.TrimSpace(got) != expected {
		t.Fatalf("expected %s=%q got=%q payload=%v", key, expected, got, payload)
	}
}
