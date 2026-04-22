package contract_test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestEnterprisePolishContract_GovernanceReports(t *testing.T) {
	ctx := newEnterpriseContractCtx(t, "workspace-owner")
	seedEnterpriseAuditData(t, ctx)
	resp := performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/governance/coverage", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("coverage status=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/enterprise/reports", `{"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,"reportType":"management","title":"周报","audienceType":"leadership","timeRange":"7d"}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create report status=%d body=%s", resp.Code, resp.Body.String())
	}
}
