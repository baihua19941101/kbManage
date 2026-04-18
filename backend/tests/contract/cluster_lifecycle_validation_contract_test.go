package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Validate(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	importResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"validate-contract",
		"displayName":"Validate Contract",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://validate.example.test"
	}`)
	if importResp.Code != http.StatusAccepted {
		t.Fatalf("import failed status=%d body=%s", importResp.Code, strings.TrimSpace(importResp.Body.String()))
	}
	validateResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/validate", `{
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"infrastructureType":"existing"
	}`)
	if validateResp.Code != http.StatusOK {
		t.Fatalf("validate failed status=%d body=%s", validateResp.Code, strings.TrimSpace(validateResp.Body.String()))
	}
}
