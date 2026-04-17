package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceRemediationLifecycleIntegration_CreateAndClose(t *testing.T) {
	t.Parallel()
	env := newComplianceIntegrationEnv(t, "compliance-remediation-int")
	baselineID := integrationCreateBaseline(t, env)
	profileID := integrationCreateProfile(t, env, baselineID, "manual")
	_ = integrationReq(t, env, http.MethodPost, "/api/v1/compliance/scan-profiles/"+strconv.FormatUint(profileID, 10)+"/execute", `{"triggerSource":"manual"}`)
	findingID := integrationFirstFindingID(t, env)
	createResp := integrationReq(t, env, http.MethodPost, "/api/v1/compliance/findings/"+strconv.FormatUint(findingID, 10)+"/remediation-tasks", `{"title":"fix","owner":"alice","priority":"high"}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create remediation failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	if !strings.Contains(createResp.Body.String(), `"status":"todo"`) {
		t.Fatalf("expected todo body=%s", createResp.Body.String())
	}
}
