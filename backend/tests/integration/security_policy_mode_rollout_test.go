package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestSecurityPolicyModeRolloutIntegration_SwitchModeForAllAssignments(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-rollout-int",
		Password: "SecurityPolicyRollout@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-rollout-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)

	createSecurityPolicyAssignmentIntegration(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID, "orders")
	createSecurityPolicyAssignmentIntegration(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID, "billing")

	switchReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/mode-switch",
		strings.NewReader(`{"targetMode":"enforce","reason":"full rollout"}`),
	)
	switchReq.Header.Set("Authorization", "Bearer "+token)
	switchReq.Header.Set("Content-Type", "application/json")
	switchResp := httptest.NewRecorder()
	app.Router.ServeHTTP(switchResp, switchReq)
	if switchResp.Code != http.StatusAccepted {
		t.Fatalf("expected mode switch status=202, got=%d body=%s", switchResp.Code, strings.TrimSpace(switchResp.Body.String()))
	}
	if !strings.Contains(switchResp.Body.String(), `"targetCount":2`) {
		t.Fatalf("expected mode switch targetCount=2 body=%s", strings.TrimSpace(switchResp.Body.String()))
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected assignments status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if strings.Count(listResp.Body.String(), `"enforcementMode":"enforce"`) < 2 {
		t.Fatalf("expected all assignments switched to enforce body=%s", strings.TrimSpace(listResp.Body.String()))
	}
}

func createSecurityPolicyAssignmentIntegration(
	t *testing.T,
	r http.Handler,
	token string,
	policyID uint64,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
	namespace string,
) {
	t.Helper()
	body := fmt.Sprintf(`{
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":["%d"],
		"namespaceRefs":["%s"],
		"resourceKinds":["Deployment"],
		"enforcementMode":"warn",
		"rolloutStage":"canary"
	}`, workspaceID, projectID, clusterID, namespace)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments",
		strings.NewReader(body),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("create assignment failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
