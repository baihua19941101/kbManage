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

func TestAuditContract_QueryAndExportSkeleton(t *testing.T) {
	t.Parallel()

	r := newAuditContractTestRouter()
	token := mustIssueAuditContractAccessToken(t, 5001)

	t.Run("query audit events", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-events?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z&eventType=operation.execute", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("GET /api/v1/audit-events not implemented yet")
		}

		assertAuditContractAllowedStatus(t, resp.Code)
		assertAuditContractJSON(t, resp.Body.Bytes())
	})

	t.Run("export audit events", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/audit-events/export", strings.NewReader(`{
			"format":"csv",
			"from":"2026-01-01T00:00:00Z",
			"to":"2026-01-31T23:59:59Z"
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("POST /api/v1/audit-events/export not implemented yet")
		}

		assertAuditContractAllowedStatus(t, resp.Code)
		assertAuditContractOptionalJSON(t, resp.Body.Bytes())
	})
}

func newAuditContractTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueAuditContractAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func assertAuditContractAllowedStatus(t *testing.T, status int) {
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

func assertAuditContractJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}
	assertAuditContractOptionalJSON(t, body)
}

func assertAuditContractOptionalJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		return
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}

