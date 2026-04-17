package contract_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/tests/testutil"
)

type complianceContractEnv struct {
	app         *testutil.App
	token       string
	workspaceID uint64
	projectID   uint64
	clusterID   uint64
}

func newComplianceContractEnv(t *testing.T, username string) complianceContractEnv {
	t.Helper()
	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: username, Password: "Compliance@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, username, "workspace-owner")
	return complianceContractEnv{app: app, token: testutil.IssueAccessToken(t, app.Config, seeded.User.ID), workspaceID: access.WorkspaceID, projectID: access.ProjectID, clusterID: access.ClusterID}
}

func complianceRequest(t *testing.T, env complianceContractEnv, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+env.token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	env.app.Router.ServeHTTP(resp, req)
	return resp
}

func complianceDecodeObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func complianceDecodeItems(t *testing.T, body []byte) []any {
	t.Helper()
	payload := complianceDecodeObject(t, body)
	raw, ok := payload["items"]
	if !ok {
		t.Fatalf("missing items payload=%v", payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("items not array type=%T", raw)
	}
	return items
}

func complianceReadID(t *testing.T, payload map[string]any, key string) uint64 {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		raw, ok = payload["ID"]
	}
	if !ok {
		t.Fatalf("missing key %s payload=%v", key, payload)
	}
	switch v := raw.(type) {
	case float64:
		return uint64(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			t.Fatalf("invalid id %s=%v", key, raw)
		}
		return id
	default:
		t.Fatalf("invalid id type %T", raw)
	}
	return 0
}

func createComplianceBaselineContract(t *testing.T, env complianceContractEnv, name string) uint64 {
	t.Helper()
	resp := complianceRequest(t, env, http.MethodPost, "/api/v1/compliance/baselines", fmt.Sprintf(`{"name":%q,"standardType":"cis","version":"1.9.0","description":"baseline"}`, name))
	if resp.Code != http.StatusCreated {
		t.Fatalf("create baseline failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return complianceReadID(t, complianceDecodeObject(t, resp.Body.Bytes()), "id")
}

func createComplianceProfileContract(t *testing.T, env complianceContractEnv, baselineID uint64, name string, scheduleMode string) uint64 {
	t.Helper()
	if scheduleMode == "" {
		scheduleMode = "manual"
	}
	body := fmt.Sprintf(`{"name":%q,"baselineId":%d,"workspaceId":%d,"projectId":%d,"scopeType":"cluster","clusterRefs":[%d],"scheduleMode":%q`, name, baselineID, env.workspaceID, env.projectID, env.clusterID, scheduleMode)
	if scheduleMode == "scheduled" {
		body += `,"cronExpression":"0 */6 * * *"`
	}
	body += `}`
	resp := complianceRequest(t, env, http.MethodPost, "/api/v1/compliance/scan-profiles", body)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create profile failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return complianceReadID(t, complianceDecodeObject(t, resp.Body.Bytes()), "id")
}

func executeComplianceProfileContract(t *testing.T, env complianceContractEnv, profileID uint64) uint64 {
	t.Helper()
	resp := complianceRequest(t, env, http.MethodPost, fmt.Sprintf("/api/v1/compliance/scan-profiles/%d/execute", profileID), `{"triggerSource":"manual"}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("execute scan failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return complianceReadID(t, complianceDecodeObject(t, resp.Body.Bytes()), "id")
}

func createExceptionPayload(start time.Time, end time.Time) string {
	return fmt.Sprintf(`{"reason":"legacy exception","startsAt":%q,"expiresAt":%q}`, start.Format(time.RFC3339), end.Format(time.RFC3339))
}
