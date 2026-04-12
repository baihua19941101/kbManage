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

func TestObservabilityScopeAuthorizationIsolationAndRevocation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)

	userA := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-scope-user-a",
		Password:    "ObsScopeA@123",
		DisplayName: "Obs Scope User A",
		Email:       "obs-scope-user-a@example.test",
	})
	userB := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-scope-user-b",
		Password:    "ObsScopeB@123",
		DisplayName: "Obs Scope User B",
		Email:       "obs-scope-user-b@example.test",
	})

	seedA := testutil.SeedObservabilityAccess(t, app.DB, userA.User.ID, "scope-a", "workspace-owner")
	seedB := testutil.SeedObservabilityAccess(t, app.DB, userB.User.ID, "scope-b", "workspace-owner")

	tokenA := testutil.IssueAccessToken(t, app.Config, userA.User.ID)
	tokenB := testutil.IssueAccessToken(t, app.Config, userB.User.ID)

	t.Run("cross-workspace cluster isolation", func(t *testing.T) {
		reqAllowed := httptest.NewRequest(http.MethodGet, "/api/v1/observability/logs/query?clusterId="+uint64ToString(seedA.ClusterID)+"&namespace=default", nil)
		reqAllowed.Header.Set("Authorization", "Bearer "+tokenA)
		respAllowed := httptest.NewRecorder()
		app.Router.ServeHTTP(respAllowed, reqAllowed)
		if respAllowed.Code != http.StatusOK {
			t.Fatalf("expected status=200, got=%d body=%s", respAllowed.Code, strings.TrimSpace(respAllowed.Body.String()))
		}

		reqDenied := httptest.NewRequest(http.MethodGet, "/api/v1/observability/logs/query?clusterId="+uint64ToString(seedB.ClusterID)+"&namespace=default", nil)
		reqDenied.Header.Set("Authorization", "Bearer "+tokenA)
		respDenied := httptest.NewRecorder()
		app.Router.ServeHTTP(respDenied, reqDenied)
		if respDenied.Code != http.StatusForbidden {
			t.Fatalf("expected status=403, got=%d body=%s", respDenied.Code, strings.TrimSpace(respDenied.Body.String()))
		}

		reqDeniedB := httptest.NewRequest(http.MethodGet, "/api/v1/observability/metrics/series?clusterId="+uint64ToString(seedA.ClusterID)+"&subjectType=pod&subjectRef=demo&metricKey=cpu_usage", nil)
		reqDeniedB.Header.Set("Authorization", "Bearer "+tokenB)
		respDeniedB := httptest.NewRecorder()
		app.Router.ServeHTTP(respDeniedB, reqDeniedB)
		if respDeniedB.Code != http.StatusForbidden {
			t.Fatalf("expected status=403, got=%d body=%s", respDeniedB.Code, strings.TrimSpace(respDeniedB.Body.String()))
		}
	})

	t.Run("permission revocation takes effect immediately", func(t *testing.T) {
		if err := app.DB.
			Where("subject_type = ? AND subject_id = ? AND scope_type = ? AND scope_id = ?", "user", userA.User.ID, "workspace", seedA.WorkspaceID).
			Delete(&repository.ScopeRoleBinding{}).Error; err != nil {
			t.Fatalf("delete role binding failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/observability/events?clusterId="+uint64ToString(seedA.ClusterID)+"&namespace=default&resourceKind=Pod&resourceName=demo", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)
		if resp.Code != http.StatusForbidden {
			t.Fatalf("expected status=403, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
	})
}

func TestObservabilityScopeAuthorizationReadOnlyCannotGovern(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-readonly-user",
		Password:    "ObsReadonly@123",
		DisplayName: "Obs Readonly User",
		Email:       "obs-readonly-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "readonly", "workspace-viewer")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/observability/alerts", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	body := `{"name":"readonly-blocked","severity":"warning","conditionExpression":"cpu > 80","scopeSnapshot":"{\"workspaceIds\":[` + uint64ToString(seed.WorkspaceID) + `]}"} `
	writeReq := httptest.NewRequest(http.MethodPost, "/api/v1/observability/alert-rules", strings.NewReader(body))
	writeReq.Header.Set("Authorization", "Bearer "+token)
	writeReq.Header.Set("Content-Type", "application/json")
	writeResp := httptest.NewRecorder()
	app.Router.ServeHTTP(writeResp, writeReq)
	if writeResp.Code != http.StatusForbidden {
		t.Fatalf("expected status=403, got=%d body=%s", writeResp.Code, strings.TrimSpace(writeResp.Body.String()))
	}
}

func uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}
