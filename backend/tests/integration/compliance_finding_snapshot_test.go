package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceFindingSnapshotIntegration_DetailIncludesEvidence(t *testing.T) {
	t.Parallel()
	env := newComplianceIntegrationEnv(t, "compliance-finding-int")
	baselineID := integrationCreateBaseline(t, env)
	profileID := integrationCreateProfile(t, env, baselineID, "manual")
	_ = integrationReq(t, env, http.MethodPost, "/api/v1/compliance/scan-profiles/"+strconv.FormatUint(profileID, 10)+"/execute", `{"triggerSource":"manual"}`)
	findingID := integrationFirstFindingID(t, env)
	resp := integrationReq(t, env, http.MethodGet, "/api/v1/compliance/findings/"+strconv.FormatUint(findingID, 10), "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "evidences") {
		t.Fatalf("finding detail failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
