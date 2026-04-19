package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_IdentitySourceLifecycle(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	sourceID := createIdentitySourceContract(t, ctx, "oidc", "optional")

	listResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), `"sourceType":"oidc"`) {
		t.Fatalf("list identity sources failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	getResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources/"+strconv.FormatUint(sourceID, 10), "")
	if getResp.Code != http.StatusOK || !strings.Contains(getResp.Body.String(), `"loginMode":"optional"`) {
		t.Fatalf("get identity source failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
}
