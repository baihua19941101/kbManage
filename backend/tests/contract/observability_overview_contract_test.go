package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_OverviewRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-contract-user",
		Password:    "ObsContract@123",
		DisplayName: "Obs Contract User",
		Email:       "obs-contract-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-overview-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	t.Run("requires bearer token", func(t *testing.T) {
		resp := performObservabilityAuthedRequest(
			t,
			app.Router,
			"",
			http.MethodGet,
			"/api/v1/observability/overview",
			"",
		)
		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("expected status=401, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeObservabilityObject(t, resp.Body.Bytes())
		errorText, _ := payload["error"].(string)
		if errorText != "missing bearer token" && errorText != "invalid token" {
			t.Fatalf("expected auth error for missing/invalid bearer token, payload=%v", payload)
		}
	})

	t.Run("returns overview json shape", func(t *testing.T) {
		query := url.Values{}
		query.Set("startAt", "2026-01-01T00:00:00Z")
		query.Set("endAt", "2026-01-01T01:00:00Z")
		resp := performObservabilityAuthedRequest(
			t,
			app.Router,
			token,
			http.MethodGet,
			"/api/v1/observability/overview?"+query.Encode(),
			"",
		)
		if resp.Code != http.StatusOK {
			t.Fatalf("expected status=200, got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeObservabilityObject(t, resp.Body.Bytes())
		assertObservabilityArrayField(t, payload, "cards")
		assertObservabilityArrayField(t, payload, "hotAlerts")
		assertObservabilityArrayField(t, payload, "topEvents")
		assertObservabilityArrayField(t, payload, "metricHighlights")
	})
}

func performObservabilityAuthedRequest(t *testing.T, h http.Handler, token, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	return resp
}

func mustDecodeObservabilityObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid json response: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func assertObservabilityArrayField(t *testing.T, payload map[string]any, field string) {
	t.Helper()
	raw, ok := payload[field]
	if !ok {
		t.Fatalf("missing %s field, payload=%v", field, payload)
	}
	if _, ok := raw.([]any); !ok {
		t.Fatalf("field %s should be array, got=%T payload=%v", field, raw, payload)
	}
}
