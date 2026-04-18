package contract_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestClusterLifecycleContract_ImportRegisterListDetail(t *testing.T) {
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "cl-contract", Password: "Contract@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "cluster-lifecycle-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	importResp := performClusterLifecycleRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"imported-contract",
		"displayName":"Imported Contract",
		"workspaceId":`+uintToString(access.WorkspaceID)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://imported.example.test"
	}`)
	if importResp.Code != http.StatusAccepted {
		t.Fatalf("import failed status=%d body=%s", importResp.Code, strings.TrimSpace(importResp.Body.String()))
	}

	registerResp := performClusterLifecycleRequest(t, app.Router, token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/register", `{
		"name":"registered-contract",
		"displayName":"Registered Contract",
		"workspaceId":`+uintToString(access.WorkspaceID)+`,
		"projectId":`+uintToString(access.ProjectID)+`,
		"infrastructureType":"registered",
		"driverKey":"register-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1"
	}`)
	if registerResp.Code != http.StatusCreated {
		t.Fatalf("register failed status=%d body=%s", registerResp.Code, strings.TrimSpace(registerResp.Body.String()))
	}

	listResp := performClusterLifecycleRequest(t, app.Router, token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if !strings.Contains(listResp.Body.String(), "registered-contract") {
		t.Fatalf("expected cluster list to contain registered-contract body=%s", strings.TrimSpace(listResp.Body.String()))
	}

	getResp := performClusterLifecycleRequest(t, app.Router, token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/2", "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get detail failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
}

func performClusterLifecycleRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	return resp
}

func uintToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}
