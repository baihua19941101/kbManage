package contract_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

type clusterLifecycleContractCtx struct {
	app    *testutil.App
	token  string
	access testutil.ObservabilityAccessSeed
}

func newClusterLifecycleContractCtx(t *testing.T, roleKey string) clusterLifecycleContractCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "cl-contract-" + roleKey, Password: "Contract@123"})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "cluster-lifecycle-contract", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	return clusterLifecycleContractCtx{app: app, token: token, access: access}
}

func performClusterLifecycleContractRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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

func createClusterLifecycleDriverForContract(t *testing.T, ctx clusterLifecycleContractCtx, capabilities string) {
	t.Helper()
	body := `{
		"driverKey":"generic-driver",
		"version":"v1",
		"displayName":"Generic Driver",
		"providerType":"generic",
		"status":"active",
		"capabilityProfileVersion":"v1",
		"schemaVersion":"v1"` + capabilities + `
	}`
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/drivers", body)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create driver failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}

func createClusterLifecycleClusterForContract(t *testing.T, ctx clusterLifecycleContractCtx) {
	t.Helper()
	createClusterLifecycleDriverForContract(t, ctx, "")
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters", `{
		"name":"cluster-contract",
		"displayName":"Cluster Contract",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.access.ProjectID, 10)+`,
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"nodePools":[{"name":"workers","role":"worker","desiredCount":3,"minCount":1,"maxCount":5,"version":"v1.30.1","zoneRefs":["zone-a"]}]
	}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("create cluster failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
