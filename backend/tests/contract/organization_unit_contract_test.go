package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_OrganizationUnits(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	createOrganizationContract(t, ctx, "org-root", 0)

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/organizations", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"name":"org-root"`) {
		t.Fatalf("list organizations failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
