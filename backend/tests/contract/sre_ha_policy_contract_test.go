package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREContract_HAPolicyAndHealthOverview(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	createResp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/ha-policies", `{
		"workspaceId":`+workspaceID+`,
		"name":"platform-ha",
		"controlPlaneScope":"platform",
		"deploymentMode":"active-active",
		"replicaExpectation":3,
		"failoverTriggerPolicy":"node-down-30s"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create ha policy failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	healthResp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/health/overview?workspaceId="+workspaceID, "")
	if healthResp.Code != http.StatusOK {
		t.Fatalf("health overview failed status=%d body=%s", healthResp.Code, strings.TrimSpace(healthResp.Body.String()))
	}
	payload := mustDecodeSREObject(t, healthResp.Body.Bytes())
	if payload["overallStatus"] == nil {
		t.Fatalf("expected overallStatus payload=%v", payload)
	}
}
