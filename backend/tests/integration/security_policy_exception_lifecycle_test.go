package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	securityPolicySvc "kbmanage/backend/internal/service/securitypolicy"
	"kbmanage/backend/internal/worker"
	"kbmanage/backend/tests/testutil"
)

func TestSecurityPolicyExceptionLifecycleIntegration_ExpireByWorker(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-exception-int",
		Password: "SecurityPolicyExceptionInt@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-exception-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	assignmentID := createSecurityPolicyExceptionLifecycleAssignment(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)

	hit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &assignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Pod",
		ResourceName:      "orders-api-0",
		HitResult:         domain.PolicyHitResultBlock,
		RiskLevel:         domain.PolicyRiskLevelHigh,
		Message:           "forbidden hostPath",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        time.Now().UTC(),
	}
	if err := app.DB.Create(hit).Error; err != nil {
		t.Fatalf("seed hit record failed: %v", err)
	}

	startsAt := time.Now().UTC().Add(-2 * time.Minute).Format(time.RFC3339)
	expiresAt := time.Now().UTC().Add(1 * time.Minute).Format(time.RFC3339)
	createBody := fmt.Sprintf(`{"reason":"临时排障窗口","startsAt":"%s","expiresAt":"%s"}`, startsAt, expiresAt)
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(hit.ID, 10)+"/exceptions",
		strings.NewReader(createBody),
	)
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create exception status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeGitOpsDeliveryUnitModelingObject(t, createResp.Body.Bytes())
	exceptionID := mustReadGitOpsDeliveryUnitModelingID(t, created, "id")

	reviewReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/exceptions/"+strconv.FormatUint(exceptionID, 10)+"/review",
		strings.NewReader(`{"decision":"approve","comment":"批准 30 分钟"}`),
	)
	reviewReq.Header.Set("Authorization", "Bearer "+token)
	reviewReq.Header.Set("Content-Type", "application/json")
	reviewResp := httptest.NewRecorder()
	app.Router.ServeHTTP(reviewResp, reviewReq)
	if reviewResp.Code != http.StatusOK {
		t.Fatalf("expected review status=200, got=%d body=%s", reviewResp.Code, strings.TrimSpace(reviewResp.Body.String()))
	}
	if !strings.Contains(reviewResp.Body.String(), `"status":"active"`) && !strings.Contains(reviewResp.Body.String(), `"status":"approved"`) {
		t.Fatalf("expected status active/approved body=%s", strings.TrimSpace(reviewResp.Body.String()))
	}

	if err := app.DB.Model(&domain.PolicyExceptionRequest{}).
		Where("id = ?", exceptionID).
		Updates(map[string]any{"status": domain.PolicyExceptionActive, "expires_at": time.Now().UTC().Add(-1 * time.Minute)}).Error; err != nil {
		t.Fatalf("force set expired timestamp failed: %v", err)
	}

	exceptionRepo := repository.NewPolicyExceptionRepository(app.DB)
	exceptionCache := securityPolicySvc.NewExceptionCache(nil, 0)
	expiryWorker := worker.NewPolicyExceptionExpiryWorker(exceptionRepo, exceptionCache, time.Second)
	changed, err := expiryWorker.RunOnce(context.Background(), time.Now().UTC())
	if err != nil {
		t.Fatalf("run expiry worker failed: %v", err)
	}
	if changed == 0 {
		t.Fatalf("expected at least one expired exception to be recycled")
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/exceptions?policyId="+strconv.FormatUint(policyID, 10)+"&status=expired&workspaceId="+strconv.FormatUint(access.WorkspaceID, 10),
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list exceptions status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if !strings.Contains(listResp.Body.String(), `"status":"expired"`) {
		t.Fatalf("expected expired exception in list body=%s", strings.TrimSpace(listResp.Body.String()))
	}
}

func createSecurityPolicyExceptionLifecycleAssignment(
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
		"resourceKinds":["Pod"],
		"enforcementMode":"warn",
		"rolloutStage":"pilot"
	}`, workspaceID, projectID, clusterID)
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
	payload := mustDecodeGitOpsDeliveryUnitModelingObject(t, listResp.Body.Bytes())
	items, ok := payload["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("missing assignment items payload=%v", payload)
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("invalid assignment item type=%T", items[0])
	}
	return mustReadGitOpsDeliveryUnitModelingID(t, first, "id")
}
