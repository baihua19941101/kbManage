package contract_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestSecurityPolicyContract_ModeSwitchAndHits(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-enforce-contract",
		Password: "SecurityPolicyEnforce@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-enforce-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	assignmentID := createSecurityPolicyContractAssignment(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)

	hit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &assignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Deployment",
		ResourceName:      "api",
		HitResult:         domain.PolicyHitResultWarn,
		RiskLevel:         domain.PolicyRiskLevelHigh,
		Message:           "image tag latest is forbidden",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        time.Now().UTC(),
	}
	if err := app.DB.Create(hit).Error; err != nil {
		t.Fatalf("seed hit record failed: %v", err)
	}

	switchReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/mode-switch",
		strings.NewReader(fmt.Sprintf(`{"targetMode":"enforce","assignmentIds":[%d],"reason":"promote canary"}`, assignmentID)),
	)
	switchReq.Header.Set("Authorization", "Bearer "+token)
	switchReq.Header.Set("Content-Type", "application/json")
	switchResp := httptest.NewRecorder()
	app.Router.ServeHTTP(switchResp, switchReq)
	if switchResp.Code != http.StatusAccepted {
		t.Fatalf("expected mode switch status=202, got=%d body=%s", switchResp.Code, strings.TrimSpace(switchResp.Body.String()))
	}
	switchPayload := mustDecodeSecurityPolicyContractObject(t, switchResp.Body.Bytes())
	if strings.TrimSpace(mustReadSecurityPolicyContractString(t, switchPayload, "operation")) != "mode-switch" {
		t.Fatalf("expected operation=mode-switch payload=%v", switchPayload)
	}

	assignmentListReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments",
		nil,
	)
	assignmentListReq.Header.Set("Authorization", "Bearer "+token)
	assignmentListResp := httptest.NewRecorder()
	app.Router.ServeHTTP(assignmentListResp, assignmentListReq)
	if assignmentListResp.Code != http.StatusOK {
		t.Fatalf("expected assignment list status=200, got=%d body=%s", assignmentListResp.Code, strings.TrimSpace(assignmentListResp.Body.String()))
	}
	if !strings.Contains(assignmentListResp.Body.String(), `"enforcementMode":"enforce"`) {
		t.Fatalf("expected assignment mode switched to enforce body=%s", strings.TrimSpace(assignmentListResp.Body.String()))
	}

	hitsReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/hits?policyId="+strconv.FormatUint(policyID, 10)+
			"&workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&riskLevel=high",
		nil,
	)
	hitsReq.Header.Set("Authorization", "Bearer "+token)
	hitsResp := httptest.NewRecorder()
	app.Router.ServeHTTP(hitsResp, hitsReq)
	if hitsResp.Code != http.StatusOK {
		t.Fatalf("expected hits list status=200, got=%d body=%s", hitsResp.Code, strings.TrimSpace(hitsResp.Body.String()))
	}
	if !strings.Contains(hitsResp.Body.String(), "image tag latest is forbidden") {
		t.Fatalf("expected hit record in response body=%s", strings.TrimSpace(hitsResp.Body.String()))
	}
}

func createSecurityPolicyContractAssignment(
	t *testing.T,
	r http.Handler,
	token string,
	policyID uint64,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":["%d"],
		"namespaceRefs":["orders"],
		"resourceKinds":["Deployment"],
		"enforcementMode":"warn",
		"rolloutStage":"pilot"
	}`, workspaceID, projectID, clusterID)
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments",
		strings.NewReader(body),
	)
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusAccepted {
		t.Fatalf("create assignment failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list assignment failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	payload := mustDecodeSecurityPolicyContractObject(t, listResp.Body.Bytes())
	items := mustReadSecurityPolicyContractArray(t, payload, "items")
	if len(items) == 0 {
		t.Fatalf("expected assignment items in payload=%v", payload)
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected assignment item type=%T", items[0])
	}
	return mustReadSecurityPolicyContractID(t, first, "id")
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
	payload := mustDecodeSecurityPolicyContractObject(t, resp.Body.Bytes())
	return mustReadSecurityPolicyContractID(t, payload, "id")
}
