package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_OrganizationMapping(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	unitID := createOrganizationContract(t, ctx, "org-map", 0)

	createResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations/"+strconv.FormatUint(unitID, 10)+"/mappings", `{
		"scopeType":"workspace",
		"scopeRef":"`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`",
		"inheritanceMode":"direct"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create mapping failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/organizations/"+strconv.FormatUint(unitID, 10)+"/mappings", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), `"scopeType":"workspace"`) {
		t.Fatalf("list mapping failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
