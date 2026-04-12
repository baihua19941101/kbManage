package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_EventsRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-events-user",
		Password:    "ObsEvents@123",
		DisplayName: "Obs Events User",
		Email:       "obs-events-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-events-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	resp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodGet,
		"/api/v1/observability/events?clusterId="+strconv.FormatUint(seed.ClusterID, 10)+"&namespace=default&resourceKind=Deployment&resourceName=mock-app",
		"",
	)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
