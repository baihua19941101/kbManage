package contract_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceRemediationContract_CreateListUpdate(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-remediation-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-rem-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "rem-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)
	findings := complianceDecodeItems(t, complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10)+"/findings", "").Body.Bytes())
	findingID := complianceReadID(t, findings[0].(map[string]any), "id")

	createResp := complianceRequest(t, env, http.MethodPost, fmt.Sprintf("/api/v1/compliance/findings/%d/remediation-tasks", findingID), `{"title":"fix audit logging","owner":"alice","priority":"high"}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create remediation failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	taskID := complianceDecodeObject(t, createResp.Body.Bytes())["id"].(string)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/remediation-tasks?workspaceId="+strconv.FormatUint(env.workspaceID, 10), "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list remediation failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	progressResp := complianceRequest(t, env, http.MethodPatch, "/api/v1/compliance/remediation-tasks/"+taskID, `{"status":"in_progress"}`)
	if progressResp.Code != http.StatusOK {
		t.Fatalf("progress remediation failed status=%d body=%s", progressResp.Code, strings.TrimSpace(progressResp.Body.String()))
	}
	updateResp := complianceRequest(t, env, http.MethodPatch, "/api/v1/compliance/remediation-tasks/"+taskID, `{"status":"done","resolutionSummary":"fixed"}`)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update remediation failed status=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
}
