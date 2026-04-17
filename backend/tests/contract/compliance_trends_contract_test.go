package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestComplianceTrendsContract_QueryTrends(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-trends-contract")
	resp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/trends", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("trends query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := complianceDecodeObject(t, resp.Body.Bytes())
	if payload["items"] == nil && payload["points"] == nil {
		t.Fatalf("unexpected trends payload=%v", payload)
	}
}
