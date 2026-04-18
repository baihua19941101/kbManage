package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_AuditQuery(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	importResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"audit-contract",
		"displayName":"Audit Contract",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://audit.example.test"
	}`)
	if importResp.Code != http.StatusAccepted {
		t.Fatalf("import failed status=%d body=%s", importResp.Code, strings.TrimSpace(importResp.Body.String()))
	}
	auditResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/audit/cluster-lifecycle/events", "")
	if auditResp.Code != http.StatusOK {
		t.Fatalf("audit query failed status=%d body=%s", auditResp.Code, strings.TrimSpace(auditResp.Body.String()))
	}
	if !strings.Contains(auditResp.Body.String(), "clusterlifecycle") {
		t.Fatalf("expected lifecycle audit body=%s", strings.TrimSpace(auditResp.Body.String()))
	}
}
