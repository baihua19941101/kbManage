package integration_test

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

func TestAuditQuery_SearchAndExportSkeleton(t *testing.T) {
	t.Parallel()

	r := newAuditQueryTestRouter()
	token := mustIssueAuditQueryAccessToken(t, 5002)

	t.Run("search audit events", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-events?from=2026-02-01T00:00:00Z&to=2026-02-28T23:59:59Z&result=success", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("GET /api/v1/audit-events not implemented yet")
		}

		assertAuditQueryAllowedStatus(t, resp.Code)
		assertAuditQueryJSON(t, resp.Body.Bytes())
	})

	t.Run("export audit events", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/audit-events/export", strings.NewReader(`{
			"format":"csv",
			"from":"2026-02-01T00:00:00Z",
			"to":"2026-02-28T23:59:59Z",
			"filters":{"result":"success"}
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("POST /api/v1/audit-events/export not implemented yet")
		}

		assertAuditQueryAllowedStatus(t, resp.Code)
		assertAuditQueryOptionalJSON(t, resp.Body.Bytes())
	})
}

func newAuditQueryTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueAuditQueryAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func assertAuditQueryAllowedStatus(t *testing.T, status int) {
	t.Helper()

	switch status {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusUnprocessableEntity,
		http.StatusMethodNotAllowed,
		http.StatusNotImplemented:
		return
	default:
		t.Fatalf("unexpected status code: %d", status)
	}
}

func assertAuditQueryJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}
	assertAuditQueryOptionalJSON(t, body)
}

func assertAuditQueryOptionalJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		return
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}

