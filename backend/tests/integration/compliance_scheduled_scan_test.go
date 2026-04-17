package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestComplianceScheduledScanIntegration_ProfileAcceptsScheduledConfig(t *testing.T) {
	t.Parallel()
	env := newComplianceIntegrationEnv(t, "compliance-scheduled-int")
	baselineID := integrationCreateBaseline(t, env)
	_ = integrationCreateProfile(t, env, baselineID, "scheduled")
	resp := integrationReq(t, env, http.MethodGet, "/api/v1/compliance/scan-profiles", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "scheduled") {
		t.Fatalf("scheduled profile query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
