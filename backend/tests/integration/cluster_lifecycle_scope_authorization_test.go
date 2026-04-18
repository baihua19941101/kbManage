package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestClusterLifecycleIntegration_ScopeAuthorization(t *testing.T) {
	owner := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	importResp := performClusterLifecycleIntegrationRequest(t, owner.app.Router, owner.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"scope-int",
		"displayName":"Scope Int",
		"workspaceId":`+strconv.FormatUint(owner.access.WorkspaceID, 10)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://scope.example.test"
	}`)
	if importResp.Code != http.StatusAccepted {
		t.Fatalf("owner import failed status=%d body=%s", importResp.Code, strings.TrimSpace(importResp.Body.String()))
	}

	viewer := testutil.SeedUser(t, owner.app.DB, testutil.SeedUserInput{Username: "scope-viewer", Password: "Viewer@123"})
	viewerToken := testutil.IssueAccessToken(t, owner.app.Config, viewer.User.ID)

	listResp := performClusterLifecycleIntegrationRequest(t, owner.app.Router, viewerToken, http.MethodGet, "/api/v1/cluster-lifecycle/clusters", "")
	if listResp.Code != http.StatusOK || strings.Contains(listResp.Body.String(), "scope-int") {
		t.Fatalf("viewer should not see cluster status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	detailResp := performClusterLifecycleIntegrationRequest(t, owner.app.Router, viewerToken, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/1", "")
	if detailResp.Code != http.StatusForbidden {
		t.Fatalf("viewer detail should be forbidden status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
}
