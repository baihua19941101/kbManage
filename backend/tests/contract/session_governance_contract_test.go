package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_SessionGovernance(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	createIdentitySourceContract(t, ctx, "ldap", "fallback")

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sessions", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"status":"active"`) {
		t.Fatalf("list sessions failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
