package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceFindingsContract_ListDetailEvidence(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-findings-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-findings-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "findings-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10)+"/findings", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list findings failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	items := complianceDecodeItems(t, listResp.Body.Bytes())
	if len(items) == 0 {
		t.Fatalf("expected findings in list")
	}
	finding := items[0].(map[string]any)
	findingID := complianceReadID(t, finding, "id")

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/findings/"+strconv.FormatUint(findingID, 10), "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get finding failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	payload := complianceDecodeObject(t, getResp.Body.Bytes())
	if payload["controlId"] == nil {
		t.Fatalf("expected controlId payload=%v", payload)
	}

	evidenceResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/findings/"+strconv.FormatUint(findingID, 10)+"/evidence", "")
	if evidenceResp.Code != http.StatusOK {
		t.Fatalf("list evidence failed status=%d body=%s", evidenceResp.Code, strings.TrimSpace(evidenceResp.Body.String()))
	}
}
