package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyContract_RoleAssignments(t *testing.T) {
	ctx := newIdentityTenancyContractCtx(t, "workspace-owner")
	roleID := createRoleDefinitionContract(t, ctx, "project-editor", "project", "bounded", true)
	validUntil := time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339)

	createResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"`+strconv.FormatUint(ctx.UserID, 10)+`",
		"roleDefinitionId":`+strconv.FormatUint(roleID, 10)+`,
		"scopeType":"project",
		"scopeRef":"`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`",
		"sourceType":"temporary",
		"validUntil":"`+validUntil+`"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create assignment failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performIdentityTenancyContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/assignments", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), `"scopeType":"project"`) {
		t.Fatalf("list assignments failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
