package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestIdentityTenancyContract_MembershipQuery(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	unitID := createOrganizationContract(t, ctx, "org-members", 0)

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/organizations/"+strconv.FormatUint(unitID, 10)+"/memberships", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"membershipRole":"owner"`) {
		t.Fatalf("list memberships failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
