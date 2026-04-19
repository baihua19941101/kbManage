package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyContract_DelegationGrant(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	createDelegationGrantContract(t, ctx, "delegate-user")

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/delegations", `{
		"grantorRef":"`+strconv.FormatUint(ctx.UserID, 10)+`",
		"delegateRef":"delegate-two",
		"allowedRoleLevels":["workspace"],
		"validFrom":"`+time.Now().UTC().Format(time.RFC3339)+`",
		"validUntil":"`+time.Now().UTC().Add(time.Hour).Format(time.RFC3339)+`",
		"reason":"second grant"
	}`)
	if resp.Code != http.StatusCreated || !strings.Contains(resp.Body.String(), `"delegateRef":"delegate-two"`) {
		t.Fatalf("delegation contract failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
