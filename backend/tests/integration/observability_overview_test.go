package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityOverviewFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-overview-int-user",
		Password:    "ObsOverview@123",
		DisplayName: "Obs Overview Int User",
		Email:       "obs-overview-int-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-overview-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/observability/overview", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if !strings.Contains(resp.Body.String(), "cards") {
		t.Fatalf("expected cards in response, body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
