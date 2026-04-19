package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREContract_CreateUpgradePlan(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades", `{
		"workspaceId":`+workspaceID+`,
		"name":"upgrade-131",
		"currentVersion":"1.30.0",
		"targetVersion":"1.31.0",
		"rolloutStrategy":"rolling"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create upgrade plan failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
