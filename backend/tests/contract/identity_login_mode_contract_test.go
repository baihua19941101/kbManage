package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_LoginModeListing(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	createIdentitySourceContract(t, ctx, "oidc", "optional")

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list login modes failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if !strings.Contains(resp.Body.String(), `"loginModes"`) || !strings.Contains(resp.Body.String(), `"sourceType":"local"`) {
		t.Fatalf("expected login modes and local fallback body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
