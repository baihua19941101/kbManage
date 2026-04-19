package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREContract_UpgradePrecheck(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/upgrades/prechecks", `{
		"workspaceId":`+workspaceID+`,
		"currentVersion":"1.30.0",
		"targetVersion":"1.31.0",
		"scope":"platform"
	}`)
	if resp.Code != http.StatusOK {
		t.Fatalf("upgrade precheck failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeSREObject(t, resp.Body.Bytes())
	if payload["decision"] == nil {
		t.Fatalf("expected decision payload=%v", payload)
	}
}
