package integration_test

import (
	"encoding/json"
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

func TestSecurityPolicyAuditQueryIntegration_FilterAndScope(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-audit-int",
		Password: "SecurityPolicyAuditInt@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-audit-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyBody := fmt.Sprintf(`{
		"name":"audit-int-policy",
		"workspaceId":%d,
		"projectId":%d,
		"scopeLevel":"project",
		"category":"label",
		"ruleTemplate":{"required":["owner"]},
		"defaultEnforcementMode":"warn",
		"riskLevel":"medium"
	}`, access.WorkspaceID, access.ProjectID)
	createPolicyReq := httptest.NewRequest(http.MethodPost, "/api/v1/security-policies", strings.NewReader(policyBody))
	createPolicyReq.Header.Set("Authorization", "Bearer "+token)
	createPolicyReq.Header.Set("Content-Type", "application/json")
	createPolicyResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createPolicyResp, createPolicyReq)
	if createPolicyResp.Code != http.StatusCreated {
		t.Fatalf("expected create policy status=201, got=%d body=%s", createPolicyResp.Code, strings.TrimSpace(createPolicyResp.Body.String()))
	}
	createdPolicy := mustDecodeSecurityPolicyAuditQueryObject(t, createPolicyResp.Body.Bytes())
	policyID := mustReadSecurityPolicyAuditQueryID(t, createdPolicy, "id")

	assignmentID := createSecurityPolicyAssignmentForAuditQuery(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)
	hit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &assignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Deployment",
		ResourceName:      "orders-api",
		HitResult:         domain.PolicyHitResultWarn,
		RiskLevel:         domain.PolicyRiskLevelHigh,
		Message:           "missing owner label",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        time.Now().UTC(),
	}
	if err := app.DB.Create(hit).Error; err != nil {
		t.Fatalf("seed policy hit failed: %v", err)
	}

	updateReq := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(hit.ID, 10)+"/remediation",
		strings.NewReader(`{"status":"in_progress","comment":"开始处理"}`),
	)
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	app.Router.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected remediation update status=200, got=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}

	queryReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/audit/security-policies/events?action=securitypolicy.hit.remediation.update&outcome=success&limit=20",
		nil,
	)
	queryReq.Header.Set("Authorization", "Bearer "+token)
	queryResp := httptest.NewRecorder()
	app.Router.ServeHTTP(queryResp, queryReq)
	if queryResp.Code != http.StatusOK {
		t.Fatalf("expected security policy audit query status=200, got=%d body=%s", queryResp.Code, strings.TrimSpace(queryResp.Body.String()))
	}
	payload := mustDecodeSecurityPolicyAuditQueryObject(t, queryResp.Body.Bytes())
	items := mustReadSecurityPolicyAuditQueryArray(t, payload, "items")
	if len(items) == 0 {
		t.Fatalf("expected remediation audit event returned payload=%v", payload)
	}

	for i := range items {
		item, ok := items[i].(map[string]any)
		if !ok {
			t.Fatalf("unexpected audit item type=%T", items[i])
		}
		action := strings.TrimSpace(mustReadSecurityPolicyAuditQueryString(t, item, "Action"))
		if action != "securitypolicy.hit.remediation.update" {
			t.Fatalf("expected remediation action, got action=%q item=%v", action, item)
		}
		resourceType := strings.TrimSpace(mustReadSecurityPolicyAuditQueryString(t, item, "ResourceType"))
		if resourceType != "securitypolicy" {
			t.Fatalf("expected resourceType=securitypolicy, got resourceType=%q item=%v", resourceType, item)
		}
	}
}

func createSecurityPolicyAssignmentForAuditQuery(
	t *testing.T,
	r http.Handler,
	token string,
	policyID uint64,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
) uint64 {
	t.Helper()
	body := `{"workspaceId":` + strconv.FormatUint(workspaceID, 10) +
		`,"projectId":` + strconv.FormatUint(projectID, 10) +
		`,"clusterRefs":["` + strconv.FormatUint(clusterID, 10) +
		`"],"namespaceRefs":["orders"],"resourceKinds":["Deployment"],"enforcementMode":"warn","rolloutStage":"canary"}`
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

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list assignment failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	payload := mustDecodeSecurityPolicyAuditQueryObject(t, listResp.Body.Bytes())
	items := mustReadSecurityPolicyAuditQueryArray(t, payload, "items")
	if len(items) == 0 {
		t.Fatalf("expected assignment list not empty payload=%v", payload)
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected assignment item type=%T", items[0])
	}
	return mustReadSecurityPolicyAuditQueryID(t, item, "id")
}

func mustDecodeSecurityPolicyAuditQueryObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadSecurityPolicyAuditQueryArray(t *testing.T, payload map[string]any, key string) []any {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	arr, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q must be array, got=%T value=%v", key, raw, raw)
	}
	return arr
}

func mustReadSecurityPolicyAuditQueryString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	text, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q must be string, got=%T value=%v", key, raw, raw)
	}
	return text
}

func mustReadSecurityPolicyAuditQueryID(t *testing.T, payload map[string]any, key string) uint64 {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	number, ok := raw.(float64)
	if !ok || number <= 0 {
		t.Fatalf("field %q must be positive number, got=%T value=%v", key, raw, raw)
	}
	return uint64(number)
}
