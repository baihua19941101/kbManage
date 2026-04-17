package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceScanExecutionIntegration_ExecuteProducesFindings(t *testing.T) {
	t.Parallel()
	env := newComplianceIntegrationEnv(t, "compliance-scan-int")
	baselineID := integrationCreateBaseline(t, env)
	profileID := integrationCreateProfile(t, env, baselineID, "manual")
	resp := integrationReq(t, env, http.MethodPost, "/api/v1/compliance/scan-profiles/"+strconv.FormatUint(profileID, 10)+"/execute", `{"triggerSource":"manual"}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("execute scan failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	listResp := integrationReq(t, env, http.MethodGet, "/api/v1/compliance/findings?workspaceId="+strconv.FormatUint(env.workspaceID, 10), "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), "ControlID") {
		t.Fatalf("findings query failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
