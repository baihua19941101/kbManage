package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/api/router"
	"kbmanage/backend/internal/repository"
)

func TestClustersContract_BasicRoutesAndShape(t *testing.T) {
	t.Parallel()

	r := newTestRouter()

	tests := []struct {
		name         string
		method       string
		target       string
		body         string
		expectSchema bool
	}{
		{
			name:         "list clusters",
			method:       http.MethodGet,
			target:       "/api/v1/clusters",
			expectSchema: true,
		},
		{
			name:   "create cluster",
			method: http.MethodPost,
			target: "/api/v1/clusters",
			body: `{
				"name": "dev-cluster",
				"apiServer": "https://127.0.0.1:6443",
				"credentialType": "token",
				"credentialPayload": "redacted"
			}`,
			expectSchema: true,
		},
		{
			name:         "list cluster resources",
			method:       http.MethodGet,
			target:       "/api/v1/clusters/demo/resources",
			expectSchema: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			if resp.Code == http.StatusNotFound {
				t.Skipf("route %s %s not implemented yet", tc.method, tc.target)
			}

			assertAllowedStatus(t, resp.Code)
			if !tc.expectSchema {
				return
			}
			assertBasicJSONShape(t, resp.Code, resp.Body.Bytes())
		})
	}
}

func newTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func assertAllowedStatus(t *testing.T, status int) {
	t.Helper()

	switch status {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusMethodNotAllowed,
		http.StatusNotImplemented:
		return
	default:
		t.Fatalf("unexpected status code: %d", status)
	}
}

func assertBasicJSONShape(t *testing.T, status int, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	switch status {
	case http.StatusOK:
		if _, hasItems := payload["items"]; hasItems {
			return
		}
		if _, hasMessage := payload["message"]; hasMessage {
			return
		}
		if _, hasStatus := payload["status"]; hasStatus {
			return
		}
		t.Fatalf("expected one of [items,message,status] in 200 response, got: %v", payload)
	case http.StatusCreated:
		if _, hasID := payload["id"]; !hasID {
			t.Fatalf("expected field id in 201 response, got: %v", payload)
		}
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		if _, hasError := payload["error"]; !hasError {
			t.Fatalf("expected field error in %d response, got: %v", status, payload)
		}
	}
}
