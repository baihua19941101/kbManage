package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREIntegration_UpgradeRolloutFlow(t *testing.T) {
	ctx := newSREIntegrationCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	precheckResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades/prechecks", `{
		"workspaceId":`+workspaceID+`,
		"currentVersion":"1.30.0",
		"targetVersion":"1.31.0",
		"scope":"platform"
	}`)
	if precheckResp.Code != http.StatusOK {
		t.Fatalf("upgrade precheck failed status=%d body=%s", precheckResp.Code, strings.TrimSpace(precheckResp.Body.String()))
	}
	createResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades", `{
		"workspaceId":`+workspaceID+`,
		"name":"integration-upgrade",
		"currentVersion":"1.30.0",
		"targetVersion":"1.31.0",
		"rolloutStrategy":"rolling"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create upgrade plan failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
}
