package contract_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceRechecksContract_CreateListGetComplete(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-rechecks-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-recheck-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "recheck-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)
	findingID := complianceReadID(t, complianceDecodeItems(t, complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10)+"/findings", "").Body.Bytes())[0].(map[string]any), "id")

	createResp := complianceRequest(t, env, http.MethodPost, fmt.Sprintf("/api/v1/compliance/findings/%d/rechecks", findingID), `{"triggerSource":"manual"}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create recheck failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	recheckID := complianceDecodeObject(t, createResp.Body.Bytes())["id"].(string)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/rechecks?workspaceId="+strconv.FormatUint(env.workspaceID, 10), "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list rechecks failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/rechecks/"+recheckID, "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get recheck failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}

	completeResp := complianceRequest(t, env, http.MethodPost, "/api/v1/compliance/rechecks/"+recheckID+"/complete", `{"passed":true,"summary":"passed"}`)
	if completeResp.Code != http.StatusOK {
		t.Fatalf("complete recheck failed status=%d body=%s", completeResp.Code, strings.TrimSpace(completeResp.Body.String()))
	}
}
