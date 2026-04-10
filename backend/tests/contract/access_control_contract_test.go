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
	"kbmanage/backend/internal/service/auth"
)

func TestAccessControlContract_WorkspaceProjectRoleBindingRoutes(t *testing.T) {
	t.Parallel()

	r := newAccessControlTestRouter()
	token := mustIssueAccessToken(t, 1001)

	tests := []struct {
		name   string
		method string
		target string
		body   string
	}{
		{
			name:   "list workspaces",
			method: http.MethodGet,
			target: "/api/v1/workspaces",
		},
		{
			name:   "create workspace",
			method: http.MethodPost,
			target: "/api/v1/workspaces",
			body: `{
				"name": "team-alpha",
				"code": "TEAM_ALPHA",
				"description": "contract test workspace"
			}`,
		},
		{
			name:   "list projects in workspace",
			method: http.MethodGet,
			target: "/api/v1/workspaces/1/projects",
		},
		{
			name:   "create role binding",
			method: http.MethodPost,
			target: "/api/v1/role-bindings",
			body: `{
				"subjectType": "user",
				"subjectId": 1002,
				"scopeType": "workspace",
				"scopeId": 1,
				"roleKey": "workspace-owner"
			}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer "+token)
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			if resp.Code == http.StatusNotFound {
				t.Skipf("route %s %s not implemented yet", tc.method, tc.target)
			}

			assertAccessControlAllowedStatus(t, resp.Code)
			assertAccessControlJSON(t, resp.Body.Bytes())
		})
	}
}

func newAccessControlTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func assertAccessControlAllowedStatus(t *testing.T, status int) {
	t.Helper()

	switch status {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusUnprocessableEntity,
		http.StatusMethodNotAllowed,
		http.StatusNotImplemented:
		return
	default:
		t.Fatalf("unexpected status code: %d", status)
	}
}

func assertAccessControlJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}
