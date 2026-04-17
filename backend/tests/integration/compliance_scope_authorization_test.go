package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestComplianceScopeAuthorizationIntegration_DeniesUserWithoutScope(t *testing.T) {
	t.Parallel()
	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{Username: "compliance-denied-user", Password: "Compliance@123"})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/compliance/findings?workspaceId=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
