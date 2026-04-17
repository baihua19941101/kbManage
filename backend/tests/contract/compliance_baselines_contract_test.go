package contract_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceBaselinesContract_CreateListGetUpdate(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-baselines-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-control-plane")

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/baselines?standardType=cis&status=draft", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list baselines failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if len(complianceDecodeItems(t, listResp.Body.Bytes())) == 0 {
		t.Fatalf("expected baselines in list")
	}

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/baselines/"+strconv.FormatUint(baselineID, 10), "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get baseline failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}

	updateResp := complianceRequest(t, env, http.MethodPatch, "/api/v1/compliance/baselines/"+strconv.FormatUint(baselineID, 10), fmt.Sprintf(`{"name":%q,"status":"active"}`, "cis-control-plane-v2"))
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update baseline failed status=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
	payload := complianceDecodeObject(t, updateResp.Body.Bytes())
	if payload["status"] != "active" && payload["Status"] != "active" {
		t.Fatalf("expected active status payload=%v", payload)
	}
}
