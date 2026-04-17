package contract_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestComplianceExceptionsContract_CreateListReview(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-exceptions-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-ex-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "ex-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)
	findingID := complianceReadID(t, complianceDecodeItems(t, complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10)+"/findings", "").Body.Bytes())[0].(map[string]any), "id")

	start := time.Now().UTC()
	end := start.Add(2 * time.Hour)
	createResp := complianceRequest(t, env, http.MethodPost, fmt.Sprintf("/api/v1/compliance/findings/%d/exceptions", findingID), createExceptionPayload(start, end))
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create exception failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	exceptionID := complianceDecodeObject(t, createResp.Body.Bytes())["id"].(string)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/exceptions?workspaceId="+strconv.FormatUint(env.workspaceID, 10), "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list exceptions failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	reviewResp := complianceRequest(t, env, http.MethodPost, "/api/v1/compliance/exceptions/"+exceptionID+"/review", `{"decision":"approve","reviewComment":"approved"}`)
	if reviewResp.Code != http.StatusOK {
		t.Fatalf("review exception failed status=%d body=%s", reviewResp.Code, strings.TrimSpace(reviewResp.Body.String()))
	}
}
