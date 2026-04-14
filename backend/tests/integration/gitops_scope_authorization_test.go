package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	"kbmanage/backend/tests/testutil"
)

func TestGitOpsScopeAuthorizationIsolationAndRevocation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)

	userA := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-scope-user-a",
		Password: "GitOpsScopeA@123",
	})
	userB := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-scope-user-b",
		Password: "GitOpsScopeB@123",
	})

	seedA := testutil.SeedObservabilityAccess(t, app.DB, userA.User.ID, "gitops-scope-a", "workspace-owner")
	seedB := testutil.SeedObservabilityAccess(t, app.DB, userB.User.ID, "gitops-scope-b", "workspace-owner")

	sourceA := seedGitOpsSourceForScope(t, app, seedA.WorkspaceID, seedA.ProjectID, "scope-source-a")
	sourceB := seedGitOpsSourceForScope(t, app, seedB.WorkspaceID, seedB.ProjectID, "scope-source-b")

	tokenA := testutil.IssueAccessToken(t, app.Config, userA.User.ID)
	tokenB := testutil.IssueAccessToken(t, app.Config, userB.User.ID)

	reqAllowed := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/sources/"+u64s(sourceA.ID), nil)
	reqAllowed.Header.Set("Authorization", "Bearer "+tokenA)
	respAllowed := httptest.NewRecorder()
	app.Router.ServeHTTP(respAllowed, reqAllowed)
	if respAllowed.Code != http.StatusOK {
		t.Fatalf("expected status=200 got=%d body=%s", respAllowed.Code, strings.TrimSpace(respAllowed.Body.String()))
	}

	reqDenied := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/sources/"+u64s(sourceB.ID), nil)
	reqDenied.Header.Set("Authorization", "Bearer "+tokenA)
	respDenied := httptest.NewRecorder()
	app.Router.ServeHTTP(respDenied, reqDenied)
	if respDenied.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respDenied.Code, strings.TrimSpace(respDenied.Body.String()))
	}

	if err := app.DB.WithContext(context.Background()).Where("subject_type = ? AND subject_id = ? AND scope_type = ? AND scope_id = ?", "user", userA.User.ID, "workspace", seedA.WorkspaceID).
		Delete(&repository.ScopeRoleBinding{}).Error; err != nil {
		t.Fatalf("delete binding failed: %v", err)
	}

	reqRevoked := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/sources/"+u64s(sourceA.ID), nil)
	reqRevoked.Header.Set("Authorization", "Bearer "+tokenA)
	respRevoked := httptest.NewRecorder()
	app.Router.ServeHTTP(respRevoked, reqRevoked)
	if respRevoked.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respRevoked.Code, strings.TrimSpace(respRevoked.Body.String()))
	}

	reqB := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/sources/"+u64s(sourceA.ID), nil)
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	respB := httptest.NewRecorder()
	app.Router.ServeHTTP(respB, reqB)
	if respB.Code != http.StatusForbidden {
		t.Fatalf("expected status=403 got=%d body=%s", respB.Code, strings.TrimSpace(respB.Body.String()))
	}
}

func seedGitOpsSourceForScope(t *testing.T, app *testutil.App, workspaceID uint64, projectID uint64, name string) *domain.DeliverySource {
	t.Helper()
	source := &domain.DeliverySource{
		Name:        name,
		SourceType:  domain.DeliverySourceTypeGit,
		Endpoint:    "https://git.example.com/" + name + ".git",
		WorkspaceID: &workspaceID,
		ProjectID:   &projectID,
		Status:      domain.DeliverySourceStatusReady,
	}
	if err := app.DB.WithContext(context.Background()).Create(source).Error; err != nil {
		t.Fatalf("seed source failed: %v", err)
	}
	return source
}
