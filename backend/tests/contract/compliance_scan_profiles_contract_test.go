package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceScanProfilesContract_CreateListGetUpdate(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-profiles-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-profile-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "daily-cluster-scan", "scheduled")

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scan-profiles?workspaceId="+strconv.FormatUint(env.workspaceID, 10)+"&projectId="+strconv.FormatUint(env.projectID, 10), "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list profiles failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if len(complianceDecodeItems(t, listResp.Body.Bytes())) == 0 {
		t.Fatalf("expected profiles in list")
	}

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scan-profiles/"+strconv.FormatUint(profileID, 10), "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get profile failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}

	updateResp := complianceRequest(t, env, http.MethodPatch, "/api/v1/compliance/scan-profiles/"+strconv.FormatUint(profileID, 10), `{"status":"active"}`)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update profile failed status=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
}
