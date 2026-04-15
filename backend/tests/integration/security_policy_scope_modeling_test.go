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

func TestSecurityPolicyScopeModelingIntegration_CreateAssignmentsAndHierarchyView(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-scope-int",
		Password: "SecurityPolicyScope@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-scope-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(
		t,
		app.Router,
		token,
		access.WorkspaceID,
		access.ProjectID,
	)

	assignBody := fmt.Sprintf(`{
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":["%d"],
		"namespaceRefs":["orders"],
		"resourceKinds":["Deployment"],
		"enforcementMode":"warn",
		"rolloutStage":"pilot"
	}`, access.WorkspaceID, access.ProjectID, access.ClusterID)

	assignReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments",
		strings.NewReader(assignBody),
	)
	assignReq.Header.Set("Authorization", "Bearer "+token)
	assignReq.Header.Set("Content-Type", "application/json")
	assignResp := httptest.NewRecorder()
	app.Router.ServeHTTP(assignResp, assignReq)
	if assignResp.Code != http.StatusAccepted {
		t.Fatalf("expected create assignment status=202, got=%d body=%s", assignResp.Code, strings.TrimSpace(assignResp.Body.String()))
	}

	listAssignmentReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments",
		nil,
	)
	listAssignmentReq.Header.Set("Authorization", "Bearer "+token)
	listAssignmentResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listAssignmentResp, listAssignmentReq)
	if listAssignmentResp.Code != http.StatusOK {
		t.Fatalf("expected list assignment status=200, got=%d body=%s", listAssignmentResp.Code, strings.TrimSpace(listAssignmentResp.Body.String()))
	}
	if !strings.Contains(listAssignmentResp.Body.String(), "\"items\"") {
		t.Fatalf("expected assignment list payload contains items, body=%s", strings.TrimSpace(listAssignmentResp.Body.String()))
	}

	listPolicyReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+"&projectId="+strconv.FormatUint(access.ProjectID, 10),
		nil,
	)
	listPolicyReq.Header.Set("Authorization", "Bearer "+token)
	listPolicyResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listPolicyResp, listPolicyReq)
	if listPolicyResp.Code != http.StatusOK {
		t.Fatalf("expected policy list status=200, got=%d body=%s", listPolicyResp.Code, strings.TrimSpace(listPolicyResp.Body.String()))
	}
	if !strings.Contains(listPolicyResp.Body.String(), "scope-modeling-policy") {
		t.Fatalf("expected created policy appears in list payload=%s", strings.TrimSpace(listPolicyResp.Body.String()))
	}
}

func createSecurityPolicyScopeModelingPolicy(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"scope-modeling-policy",
		"workspaceId":%d,
		"projectId":%d,
		"scopeLevel":"project",
		"category":"label",
		"ruleTemplate":{"required":["team"]},
		"defaultEnforcementMode":"warn",
		"riskLevel":"medium"
	}`, workspaceID, projectID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/security-policies", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create security-policy fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsDeliveryUnitModelingObject(t, resp.Body.Bytes())
	return mustReadGitOpsDeliveryUnitModelingID(t, payload, "id")
}
