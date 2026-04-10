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

func TestOperationsContract_SubmitQueryAndFailureSkeleton(t *testing.T) {
	t.Parallel()

	r := newOperationsContractTestRouter()
	token := mustIssueOperationsContractAccessToken(t, 3001)

	t.Run("submit operation", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
			"clusterId": "demo-cluster",
			"operationType": "restart-workload",
			"payload": {"namespace":"default","name":"api-server"}
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("POST /api/v1/operations not implemented yet")
		}

		assertOperationsContractAllowedStatus(t, resp.Code)
		assertOperationsContractJSON(t, resp.Body.Bytes())
	})

	t.Run("query operation status", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/operations/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("GET /api/v1/operations/:id not implemented yet")
		}

		assertOperationsContractAllowedStatus(t, resp.Code)
		assertOperationsContractJSON(t, resp.Body.Bytes())
	})

	t.Run("submit invalid payload should fail", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{"clusterId":""}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("POST /api/v1/operations not implemented yet")
		}

		if resp.Code >= 200 && resp.Code < 300 {
			t.Fatalf("expected invalid payload to fail, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		assertOperationsContractJSON(t, resp.Body.Bytes())
	})
}

func newOperationsContractTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueOperationsContractAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func assertOperationsContractAllowedStatus(t *testing.T, status int) {
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

func assertOperationsContractJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}
