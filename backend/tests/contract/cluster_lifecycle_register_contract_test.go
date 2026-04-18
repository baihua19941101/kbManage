package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Register(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/register", `{
		"name":"registered-contract",
		"displayName":"Registered Contract",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.access.ProjectID, 10)+`,
		"infrastructureType":"registered",
		"driverKey":"register-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("register failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if !strings.Contains(resp.Body.String(), `"registrationToken"`) {
		t.Fatalf("expected registration token body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
