package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"
)

func TestSecurityPolicyViolationLifecycleIntegration_RemediationStateFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-violation-int",
		Password: "SecurityPolicyViolation@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-violation-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	assignmentID := createSecurityPolicyAssignmentForViolationLifecycle(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)

	detectedAt := time.Now().UTC()
	hit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &assignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Deployment",
		ResourceName:      "orders-api",
		HitResult:         domain.PolicyHitResultBlock,
		RiskLevel:         domain.PolicyRiskLevelHigh,
		Message:           "missing required label",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        detectedAt,
	}
	if err := app.DB.Create(hit).Error; err != nil {
		t.Fatalf("seed policy hit failed: %v", err)
	}

	mustUpdateRemediationStatus(t, app.Router, token, hit.ID, "in_progress", "开始排查", http.StatusOK)
	mitigatedPayload := mustUpdateRemediationStatus(t, app.Router, token, hit.ID, "mitigated", "完成修复", http.StatusOK)
	if mitigatedPayload["detectedAt"] == nil {
		t.Fatalf("expected remediation response includes detectedAt payload=%v", mitigatedPayload)
	}
	mustUpdateRemediationStatus(t, app.Router, token, hit.ID, "closed", "关闭工单", http.StatusOK)

	invalidReq := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(hit.ID, 10)+"/remediation",
		strings.NewReader(`{"status":"mitigated"}`),
	)
	invalidReq.Header.Set("Authorization", "Bearer "+token)
	invalidReq.Header.Set("Content-Type", "application/json")
	invalidResp := httptest.NewRecorder()
	app.Router.ServeHTTP(invalidResp, invalidReq)
	if invalidResp.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid transition status=400, got=%d body=%s", invalidResp.Code, strings.TrimSpace(invalidResp.Body.String()))
	}
	mustUpdateRemediationStatus(t, app.Router, token, hit.ID, "open", "复测复开", http.StatusOK)

	from := detectedAt.Add(-5 * time.Minute).Format(time.RFC3339)
	to := detectedAt.Add(5 * time.Minute).Format(time.RFC3339)
	queryReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/hits?policyId="+strconv.FormatUint(policyID, 10)+
			"&workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&clusterId="+strconv.FormatUint(access.ClusterID, 10)+
			"&namespace=orders&enforcementMode=warn&riskLevel=high&remediationStatus=open"+
			"&from="+from+"&to="+to,
		nil,
	)
	queryReq.Header.Set("Authorization", "Bearer "+token)
	queryResp := httptest.NewRecorder()
	app.Router.ServeHTTP(queryResp, queryReq)
	if queryResp.Code != http.StatusOK {
		t.Fatalf("expected filtered hits query status=200, got=%d body=%s", queryResp.Code, strings.TrimSpace(queryResp.Body.String()))
	}
	if !strings.Contains(queryResp.Body.String(), "orders-api") {
		t.Fatalf("expected filtered hit body=%s", strings.TrimSpace(queryResp.Body.String()))
	}
}

func createSecurityPolicyAssignmentForViolationLifecycle(
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
		`"],"namespaceRefs":["orders"],"resourceKinds":["Deployment"],"enforcementMode":"warn","rolloutStage":"pilot"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10)+"/assignments", strings.NewReader(body))
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
	payload := mustDecodeViolationLifecycleObject(t, listResp.Body.Bytes())
	items := mustReadViolationLifecycleArray(t, payload, "items")
	if len(items) == 0 {
		t.Fatalf("expected assignment items payload=%v", payload)
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("invalid assignment item type=%T", items[0])
	}
	return mustReadViolationLifecycleID(t, first, "id")
}

func mustUpdateRemediationStatus(
	t *testing.T,
	r http.Handler,
	token string,
	hitID uint64,
	status string,
	comment string,
	expectedStatus int,
) map[string]any {
	t.Helper()
	body := `{"status":"` + status + `","comment":"` + comment + `"}`
	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(hitID, 10)+"/remediation",
		strings.NewReader(body),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != expectedStatus {
		t.Fatalf("update remediation status=%s expected=%d got=%d body=%s", status, expectedStatus, resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustDecodeViolationLifecycleObject(t, resp.Body.Bytes())
}

func mustDecodeViolationLifecycleObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadViolationLifecycleArray(t *testing.T, payload map[string]any, key string) []any {
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

func mustReadViolationLifecycleID(t *testing.T, payload map[string]any, key string) uint64 {
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
