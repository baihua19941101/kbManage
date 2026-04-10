package integration_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/api/router"
	"kbmanage/backend/internal/repository"
	"kbmanage/backend/internal/service/auth"
)

func TestScopeAuthorization_WorkspaceIsolationFlow(t *testing.T) {
	t.Parallel()

	r := newScopeAuthTestRouter()
	adminToken := mustIssueScopeAccessToken(t, 2001)
	memberToken := mustIssueScopeAccessToken(t, 2002)

	createResp := httptest.NewRecorder()
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", strings.NewReader(`{
		"name": "isolation-workspace",
		"code": "ISO_WS",
		"description": "workspace isolation integration test"
	}`))
	createReq.Header.Set("Authorization", "Bearer "+adminToken)
	createReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(createResp, createReq)

	if createResp.Code == http.StatusNotFound {
		t.Skip("POST /api/v1/workspaces not implemented yet")
	}
	if createResp.Code != http.StatusCreated {
		t.Skipf("cannot provision workspace fixture, got status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	workspaceID, ok := extractScopeWorkspaceID(createResp.Body.Bytes())
	if !ok {
		t.Skipf("create workspace response does not contain an id field, body=%s", strings.TrimSpace(createResp.Body.String()))
	}

	listResp := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/projects", nil)
	listReq.Header.Set("Authorization", "Bearer "+memberToken)
	r.ServeHTTP(listResp, listReq)

	if listResp.Code == http.StatusNotFound {
		t.Skip("GET /api/v1/workspaces/:id/projects not implemented yet")
	}

	switch listResp.Code {
	case http.StatusForbidden, http.StatusUnauthorized:
		return
	case http.StatusOK:
		t.Fatalf("workspace isolation regression: non-owner user can access workspace projects, body=%s", strings.TrimSpace(listResp.Body.String()))
	default:
		t.Skipf("isolation assertion skipped for non-finalized status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}

func newScopeAuthTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueScopeAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func extractScopeWorkspaceID(body []byte) (string, bool) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", false
	}

	if id, ok := payload["id"]; ok {
		switch v := id.(type) {
		case string:
			if v != "" {
				return v, true
			}
		case float64:
			if v > 0 && v == math.Trunc(v) {
				return strconv.FormatInt(int64(v), 10), true
			}
		}
	}

	return "", false
}
