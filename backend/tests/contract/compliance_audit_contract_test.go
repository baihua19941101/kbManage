package contract_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestComplianceAuditContract_QueryComplianceEvents(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-audit-contract")
	baselineID := createComplianceBaselineContract(t, env, "cis-audit-baseline")
	profileID := createComplianceProfileContract(t, env, baselineID, "audit-scan", "manual")
	scanID := executeComplianceProfileContract(t, env, profileID)
	findingID := complianceReadID(t, complianceDecodeItems(t, complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/scans/"+strconv.FormatUint(scanID, 10)+"/findings", "").Body.Bytes())[0].(map[string]any), "id")
	start := time.Now().UTC().Add(-time.Hour)
	end := time.Now().UTC().Add(time.Hour)
	createResp := complianceRequest(t, env, http.MethodPost, fmt.Sprintf("/api/v1/compliance/findings/%d/exceptions", findingID), createExceptionPayload(start, end))
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create exception failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	resp := complianceRequest(t, env, http.MethodGet, "/api/v1/audit/compliance/events?action=compliance.exception.request&timeFrom="+start.Format(time.RFC3339)+"&timeTo="+end.Format(time.RFC3339), "")
	if resp.Code != http.StatusOK {
		t.Fatalf("query compliance audit failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := complianceDecodeObject(t, resp.Body.Bytes())
	if payload["count"] == nil {
		t.Fatalf("expected count payload=%v", payload)
	}
}
