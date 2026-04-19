package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREContract_RollbackValidation(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	createResp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades", `{
		"workspaceId":`+workspaceID+`,
		"name":"upgrade-rollback",
		"currentVersion":"1.30.0",
		"targetVersion":"1.31.0",
		"rolloutStrategy":"rolling"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create upgrade plan failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades/1/rollback-validations", `{
		"validationScope":"platform",
		"result":"passed",
		"remainingRisk":"none"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create rollback validation failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
