package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyContract_AccessRiskQuery(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	roleID := createRoleDefinitionContract(t, ctx, "risk-project-editor", "project", "downward-allowed", true)
	validUntil := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"`+strconv.FormatUint(ctx.UserID, 10)+`",
		"roleDefinitionId":`+strconv.FormatUint(roleID, 10)+`,
		"scopeType":"project",
		"scopeRef":"`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`",
		"sourceType":"temporary",
		"validUntil":"`+validUntil+`"
	}`)

	resp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/access-risks", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"riskType"`) {
		t.Fatalf("access risk query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
