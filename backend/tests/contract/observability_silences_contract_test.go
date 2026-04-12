package contract_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_SilencesRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-silences-contract-user",
		Password:    "ObsSilences@123",
		DisplayName: "Obs Silences Contract User",
		Email:       "obs-silences-contract-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-silences-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	now := time.Now().UTC()
	createResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodPost,
		"/api/v1/observability/silences",
		`{"name":"night-maintenance","reason":"release","startsAt":"`+now.Add(-5*time.Minute).Format(time.RFC3339)+`","endsAt":"`+now.Add(30*time.Minute).Format(time.RFC3339)+`"}`,
	)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodGet,
		"/api/v1/observability/silences",
		"",
	)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	payload := mustDecodeObservabilityObject(t, listResp.Body.Bytes())
	assertObservabilityArrayField(t, payload, "items")
}
