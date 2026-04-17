package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

type complianceIntegrationEnv struct {
	app         *testutil.App
	token       string
	workspaceID uint64
	projectID   uint64
	clusterID   uint64
}

func newComplianceIntegrationEnv(t *testing.T, username string) complianceIntegrationEnv {
	t.Helper()
	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: username, Password: "Compliance@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, username, "workspace-owner")
	return complianceIntegrationEnv{app: app, token: testutil.IssueAccessToken(t, app.Config, seeded.User.ID), workspaceID: access.WorkspaceID, projectID: access.ProjectID, clusterID: access.ClusterID}
}

func integrationReq(t *testing.T, env complianceIntegrationEnv, method, path, body string) *httptest.ResponseRecorder {
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

func integrationDecodeObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func integrationReadID(t *testing.T, payload map[string]any) uint64 {
	t.Helper()
	raw, ok := payload["id"]
	if !ok {
		raw, ok = payload["ID"]
	}
	if !ok {
		t.Fatalf("missing id payload=%v", payload)
	}
	switch v := raw.(type) {
	case float64:
		return uint64(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			t.Fatalf("invalid id payload=%v", payload)
		}
		return id
	default:
		t.Fatalf("unsupported id type=%T", raw)
	}
	return 0
}

func integrationCreateBaseline(t *testing.T, env complianceIntegrationEnv) uint64 {
	t.Helper()
	resp := integrationReq(t, env, http.MethodPost, "/api/v1/compliance/baselines", `{"name":"cis-int","standardType":"cis","version":"1.9.0"}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create baseline failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return integrationReadID(t, integrationDecodeObject(t, resp.Body.Bytes()))
}

func integrationCreateProfile(t *testing.T, env complianceIntegrationEnv, baselineID uint64, mode string) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{"name":"int-profile","baselineId":%d,"workspaceId":%d,"projectId":%d,"scopeType":"cluster","clusterRefs":[%d],"scheduleMode":%q`, baselineID, env.workspaceID, env.projectID, env.clusterID, mode)
	if mode == "scheduled" {
		body += `,"cronExpression":"0 */6 * * *"`
	}
	body += `}`
	resp := integrationReq(t, env, http.MethodPost, "/api/v1/compliance/scan-profiles", body)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create profile failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return integrationReadID(t, integrationDecodeObject(t, resp.Body.Bytes()))
}

func integrationFirstFindingID(t *testing.T, env complianceIntegrationEnv) uint64 {
	t.Helper()
	resp := integrationReq(t, env, http.MethodGet, "/api/v1/compliance/findings", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list findings failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := integrationDecodeObject(t, resp.Body.Bytes())
	raw, ok := payload["items"]
	if !ok {
		t.Fatalf("missing items payload=%v", payload)
	}
	items, ok := raw.([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("no finding items payload=%v", payload)
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("invalid finding item type=%T", items[0])
	}
	return integrationReadID(t, first)
}
