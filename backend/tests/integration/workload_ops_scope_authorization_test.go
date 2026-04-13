package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/internal/repository"
	"kbmanage/backend/tests/testutil"
)

func TestWorkloadOpsScopeAuthorizationIsolationAndRevocation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)

	userA := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-scope-user-a",
		Password: "WopsScopeA@123",
	})
	userB := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-scope-user-b",
		Password: "WopsScopeB@123",
	})

	seedA := testutil.SeedObservabilityAccess(t, app.DB, userA.User.ID, "wops-scope-a", "workspace-owner")
	seedB := testutil.SeedObservabilityAccess(t, app.DB, userB.User.ID, "wops-scope-b", "workspace-owner")

	tokenA := testutil.IssueAccessToken(t, app.Config, userA.User.ID)
	tokenB := testutil.IssueAccessToken(t, app.Config, userB.User.ID)

	reqAllowed := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/context?clusterId="+u64s(seedA.ClusterID)+"&namespace=default&resourceKind=Deployment&resourceName=demo", nil)
	reqAllowed.Header.Set("Authorization", "Bearer "+tokenA)
	respAllowed := httptest.NewRecorder()
	app.Router.ServeHTTP(respAllowed, reqAllowed)
	if respAllowed.Code != http.StatusOK {
		t.Fatalf("expected status=200 got=%d body=%s", respAllowed.Code, strings.TrimSpace(respAllowed.Body.String()))
	}

	reqDenied := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/context?clusterId="+u64s(seedB.ClusterID)+"&namespace=default&resourceKind=Deployment&resourceName=demo", nil)
	reqDenied.Header.Set("Authorization", "Bearer "+tokenA)
	respDenied := httptest.NewRecorder()
	app.Router.ServeHTTP(respDenied, reqDenied)
	if respDenied.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respDenied.Code, strings.TrimSpace(respDenied.Body.String()))
	}

	if err := app.DB.Where("subject_type = ? AND subject_id = ? AND scope_type = ? AND scope_id = ?", "user", userA.User.ID, "workspace", seedA.WorkspaceID).
		Delete(&repository.ScopeRoleBinding{}).Error; err != nil {
		t.Fatalf("delete binding failed: %v", err)
	}

	reqRevoked := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/actions", strings.NewReader(`{
		"clusterId":`+u64s(seedA.ClusterID)+`,
		"namespace":"default",
		"resourceKind":"Deployment",
		"resourceName":"demo",
		"actionType":"restart",
		"riskConfirmed":true
	}`))
	reqRevoked.Header.Set("Authorization", "Bearer "+tokenA)
	reqRevoked.Header.Set("Content-Type", "application/json")
	respRevoked := httptest.NewRecorder()
	app.Router.ServeHTTP(respRevoked, reqRevoked)
	if respRevoked.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respRevoked.Code, strings.TrimSpace(respRevoked.Body.String()))
	}

	reqB := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/resources/instances?clusterId="+u64s(seedA.ClusterID)+"&namespace=default&resourceKind=Deployment&resourceName=demo", nil)
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	respB := httptest.NewRecorder()
	app.Router.ServeHTTP(respB, reqB)
	if respB.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respB.Code, strings.TrimSpace(respB.Body.String()))
	}
}

func u64s(v uint64) string {
	return strconv.FormatUint(v, 10)
}
