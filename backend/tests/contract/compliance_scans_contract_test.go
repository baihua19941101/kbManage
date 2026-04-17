package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceScansContract_ExecuteListGet(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-scans-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-scan-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "manual-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans?workspaceId="+strconv.FormatUint(env.workspaceID, 10)+"&projectId="+strconv.FormatUint(env.projectID, 10), "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list scans failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if len(complianceDecodeItems(t, listResp.Body.Bytes())) == 0 {
		t.Fatalf("expected scans in list")
	}

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10), "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get scan failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
}
