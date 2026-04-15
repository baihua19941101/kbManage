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

func TestSecurityPolicyContract_HitsQueryAndRemediation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-hits-contract",
		Password: "SecurityPolicyHits@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-hits-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	warnAssignmentID := createSecurityPolicyContractAssignment(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)
	enforceAssignmentID := createSecurityPolicyContractAssignmentWithMode(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID, "enforce", "canary")

	now := time.Now().UTC()
	targetHit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &enforceAssignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Deployment",
		ResourceName:      "orders-api",
		HitResult:         domain.PolicyHitResultBlock,
		RiskLevel:         domain.PolicyRiskLevelCritical,
		Message:           "critical violation",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        now,
	}
	if err := app.DB.Create(targetHit).Error; err != nil {
		t.Fatalf("seed target hit failed: %v", err)
	}
	otherHit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &warnAssignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "billing",
		ResourceKind:      "Deployment",
		ResourceName:      "billing-api",
		HitResult:         domain.PolicyHitResultWarn,
		RiskLevel:         domain.PolicyRiskLevelLow,
		Message:           "low risk warning",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        now.Add(-2 * time.Hour),
	}
	if err := app.DB.Create(otherHit).Error; err != nil {
		t.Fatalf("seed other hit failed: %v", err)
	}

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)
	hitsReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/hits?policyId="+strconv.FormatUint(policyID, 10)+
			"&workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&enforcementMode=enforce&riskLevel=critical&remediationStatus=open"+
			"&from="+from+"&to="+to,
		nil,
	)
	hitsReq.Header.Set("Authorization", "Bearer "+token)
	hitsResp := httptest.NewRecorder()
	app.Router.ServeHTTP(hitsResp, hitsReq)
	if hitsResp.Code != http.StatusOK {
		t.Fatalf("expected hits query status=200, got=%d body=%s", hitsResp.Code, strings.TrimSpace(hitsResp.Body.String()))
	}
	hitsPayload := mustDecodeSecurityPolicyContractObject(t, hitsResp.Body.Bytes())
	hits := mustReadSecurityPolicyContractArray(t, hitsPayload, "items")
	if len(hits) != 1 {
		t.Fatalf("expected exactly one hit after filters, got=%d payload=%v", len(hits), hitsPayload)
	}
	hitObject, ok := hits[0].(map[string]any)
	if !ok {
		t.Fatalf("unexpected hit payload item type=%T", hits[0])
	}
	if got := mustReadSecurityPolicyContractString(t, hitObject, "resourceName"); strings.TrimSpace(got) != "orders-api" {
		t.Fatalf("expected filtered hit resourceName=orders-api, got=%q payload=%v", got, hitObject)
	}

	updateReq := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(targetHit.ID, 10)+"/remediation",
		strings.NewReader(`{"status":"in_progress","comment":"开始整改"}`),
	)
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	app.Router.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected remediation update status=200, got=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
	updated := mustDecodeSecurityPolicyContractObject(t, updateResp.Body.Bytes())
	if got := strings.TrimSpace(mustReadSecurityPolicyContractString(t, updated, "remediationStatus")); got != "in_progress" {
		t.Fatalf("expected remediationStatus=in_progress, got=%q payload=%v", got, updated)
	}

	queryUpdatedReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/hits?policyId="+strconv.FormatUint(policyID, 10)+"&remediationStatus=in_progress",
		nil,
	)
	queryUpdatedReq.Header.Set("Authorization", "Bearer "+token)
	queryUpdatedResp := httptest.NewRecorder()
	app.Router.ServeHTTP(queryUpdatedResp, queryUpdatedReq)
	if queryUpdatedResp.Code != http.StatusOK {
		t.Fatalf("expected updated hits query status=200, got=%d body=%s", queryUpdatedResp.Code, strings.TrimSpace(queryUpdatedResp.Body.String()))
	}
	if !strings.Contains(queryUpdatedResp.Body.String(), "orders-api") {
		t.Fatalf("expected updated hit appears in remediation filter body=%s", strings.TrimSpace(queryUpdatedResp.Body.String()))
	}
}

func createSecurityPolicyContractAssignmentWithMode(
	t *testing.T,
	r http.Handler,
	token string,
	policyID uint64,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
	enforcementMode string,
	rolloutStage string,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":["%d"],
		"namespaceRefs":["orders"],
		"resourceKinds":["Deployment"],
		"enforcementMode":"%s",
		"rolloutStage":"%s"
	}`, workspaceID, projectID, clusterID, enforcementMode, rolloutStage)
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
	for i := range items {
		item, ok := items[i].(map[string]any)
		if !ok {
			continue
		}
		if mode, ok := item["enforcementMode"].(string); ok && strings.EqualFold(strings.TrimSpace(mode), enforcementMode) {
			return mustReadSecurityPolicyContractID(t, item, "id")
		}
	}
	t.Fatalf("assignment with mode=%s not found payload=%v", enforcementMode, payload)
	return 0
}
