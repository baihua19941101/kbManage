package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceOverviewContract_QueryOverview(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-overview-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-overview-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "overview-scan", "manual")
	_ = executeComplianceProfileContract(t, env, profileID)

	resp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/overview?workspaceId="+strconv.FormatUint(env.workspaceID, 10)+"&projectId="+strconv.FormatUint(env.projectID, 10)+"&groupBy=cluster", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("overview query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := complianceDecodeObject(t, resp.Body.Bytes())
	if payload["coverageRate"] == nil {
		t.Fatalf("expected coverageRate payload=%v", payload)
	}
}
