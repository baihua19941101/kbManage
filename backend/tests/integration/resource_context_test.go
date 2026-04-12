package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityResourceContextFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-context-int-user",
		Password:    "ObsContext@123",
		DisplayName: "Obs Context Int User",
		Email:       "obs-context-int-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-context-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/observability/resources/context?clusterId="+strconv.FormatUint(seed.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=mock-app",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	body := resp.Body.String()
	if !strings.Contains(body, "resourceRef") || !strings.Contains(body, "logSummary") {
		t.Fatalf("expected resource context fields in response, body=%s", strings.TrimSpace(body))
	}
}
