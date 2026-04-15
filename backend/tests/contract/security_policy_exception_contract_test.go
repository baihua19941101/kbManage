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

func TestSecurityPolicyContract_CreateAndReviewException(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-exception-contract",
		Password: "SecurityPolicyException@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-exception-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	policyID := createSecurityPolicyScopeModelingPolicy(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	assignmentID := createSecurityPolicyContractAssignment(t, app.Router, token, policyID, access.WorkspaceID, access.ProjectID, access.ClusterID)

	hit := &domain.PolicyHitRecord{
		PolicyID:          policyID,
		AssignmentID:      &assignmentID,
		ClusterID:         &access.ClusterID,
		Namespace:         "orders",
		ResourceKind:      "Pod",
		ResourceName:      "orders-api-0",
		HitResult:         domain.PolicyHitResultBlock,
		RiskLevel:         domain.PolicyRiskLevelCritical,
		Message:           "privileged container detected",
		RemediationStatus: domain.PolicyRemediationOpen,
		DetectedAt:        time.Now().UTC(),
	}
	if err := app.DB.Create(hit).Error; err != nil {
		t.Fatalf("seed hit record failed: %v", err)
	}

	startsAt := time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339)
	expiresAt := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	createExceptionBody := fmt.Sprintf(`{"reason":"业务高峰期豁免","startsAt":"%s","expiresAt":"%s"}`, startsAt, expiresAt)
	createExceptionReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/hits/"+strconv.FormatUint(hit.ID, 10)+"/exceptions",
		strings.NewReader(createExceptionBody),
	)
	createExceptionReq.Header.Set("Authorization", "Bearer "+token)
	createExceptionReq.Header.Set("Content-Type", "application/json")
	createExceptionResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createExceptionResp, createExceptionReq)
	if createExceptionResp.Code != http.StatusCreated {
		t.Fatalf("expected create exception status=201, got=%d body=%s", createExceptionResp.Code, strings.TrimSpace(createExceptionResp.Body.String()))
	}
	created := mustDecodeSecurityPolicyContractObject(t, createExceptionResp.Body.Bytes())
	exceptionID := mustReadSecurityPolicyContractID(t, created, "id")
	if strings.TrimSpace(mustReadSecurityPolicyContractString(t, created, "status")) != "pending" {
		t.Fatalf("expected created exception status pending payload=%v", created)
	}

	reviewReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/security-policies/exceptions/"+strconv.FormatUint(exceptionID, 10)+"/review",
		strings.NewReader(`{"decision":"approve","comment":"批准临时例外"}`),
	)
	reviewReq.Header.Set("Authorization", "Bearer "+token)
	reviewReq.Header.Set("Content-Type", "application/json")
	reviewResp := httptest.NewRecorder()
	app.Router.ServeHTTP(reviewResp, reviewReq)
	if reviewResp.Code != http.StatusOK {
		t.Fatalf("expected review status=200, got=%d body=%s", reviewResp.Code, strings.TrimSpace(reviewResp.Body.String()))
	}
	reviewed := mustDecodeSecurityPolicyContractObject(t, reviewResp.Body.Bytes())
	reviewedStatus := strings.TrimSpace(mustReadSecurityPolicyContractString(t, reviewed, "status"))
	if reviewedStatus != "active" && reviewedStatus != "approved" {
		t.Fatalf("expected reviewed status active/approved payload=%v", reviewed)
	}

	exceptionsReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies/exceptions?policyId="+strconv.FormatUint(policyID, 10)+"&workspaceId="+strconv.FormatUint(access.WorkspaceID, 10),
		nil,
	)
	exceptionsReq.Header.Set("Authorization", "Bearer "+token)
	exceptionsResp := httptest.NewRecorder()
	app.Router.ServeHTTP(exceptionsResp, exceptionsReq)
	if exceptionsResp.Code != http.StatusOK {
		t.Fatalf("expected list exceptions status=200, got=%d body=%s", exceptionsResp.Code, strings.TrimSpace(exceptionsResp.Body.String()))
	}
	if !strings.Contains(exceptionsResp.Body.String(), "业务高峰期豁免") {
		t.Fatalf("expected exception appears in list body=%s", strings.TrimSpace(exceptionsResp.Body.String()))
	}
}
