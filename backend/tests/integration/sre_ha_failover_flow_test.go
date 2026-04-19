package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREIntegration_HAPolicyAndHealthOverviewFlow(t *testing.T) {
	ctx := newSREIntegrationCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	createResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/ha-policies", `{
		"workspaceId":`+workspaceID+`,
		"name":"integration-ha",
		"controlPlaneScope":"platform",
		"deploymentMode":"active-active",
		"replicaExpectation":3,
		"failoverTriggerPolicy":"node-down-30s"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create ha policy failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	overviewResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/health/overview?workspaceId="+workspaceID, "")
	if overviewResp.Code != http.StatusOK {
		t.Fatalf("health overview failed status=%d body=%s", overviewResp.Code, strings.TrimSpace(overviewResp.Body.String()))
	}
}
