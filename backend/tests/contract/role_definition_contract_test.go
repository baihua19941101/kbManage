package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_RoleDefinitions(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	createRoleDefinitionContract(t, ctx, "workspace-admin", "workspace", "bounded", true)

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/roles", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"roleLevel":"workspace"`) {
		t.Fatalf("list roles failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
