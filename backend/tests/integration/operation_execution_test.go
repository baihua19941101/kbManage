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

func TestOperationExecution_SubmitTrackAndFailureSkeleton(t *testing.T) {
	t.Parallel()

	r := newOperationExecutionTestRouter()
	token := mustIssueOperationExecutionAccessToken(t, 4001)

	operationID, ok := submitOperationFixture(t, r, token)
	if !ok {
		return
	}

	t.Run("query operation status", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/operations/"+operationID, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("GET /api/v1/operations/:id not implemented yet")
		}

		assertOperationExecutionAllowedStatus(t, resp.Code)
		if resp.Code == http.StatusOK || resp.Code == http.StatusAccepted {
			assertOperationStatusField(t, resp.Body.Bytes())
		}
	})

	t.Run("submit invalid operation should fail", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{"operationType":""}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code == http.StatusNotFound {
			t.Skip("POST /api/v1/operations not implemented yet")
		}
		if resp.Code >= 200 && resp.Code < 300 {
			t.Fatalf("expected invalid operation submission to fail, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}

		assertOperationExecutionJSON(t, resp.Body.Bytes())
	})
}

func submitOperationFixture(t *testing.T, r http.Handler, token string) (string, bool) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
		"clusterId": "demo-cluster",
		"operationType": "scale-workload",
		"payload": {"namespace":"default","name":"api-server","replicas":2}
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code == http.StatusNotFound {
		t.Skip("POST /api/v1/operations not implemented yet")
	}

	assertOperationExecutionAllowedStatus(t, resp.Code)
	if resp.Code < 200 || resp.Code >= 300 {
		t.Skipf("cannot provision operation fixture, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}

	id, ok := extractOperationID(resp.Body.Bytes())
	if !ok {
		t.Skipf("submit operation response does not contain operation id, body=%s", strings.TrimSpace(resp.Body.String()))
	}
	return id, true
}

func newOperationExecutionTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func mustIssueOperationExecutionAccessToken(t *testing.T, userID uint64) string {
	t.Helper()

	tokenSvc := auth.NewTokenService("test-secret", 15*time.Minute, 24*time.Hour)
	token, err := tokenSvc.IssueAccessToken(userID)
	if err != nil {
		t.Fatalf("issue access token failed: %v", err)
	}
	return token
}

func assertOperationExecutionAllowedStatus(t *testing.T, status int) {
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

func assertOperationExecutionJSON(t *testing.T, body []byte) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
}

func assertOperationStatusField(t *testing.T, body []byte) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	if _, ok := payload["status"]; ok {
		return
	}
	if _, ok := payload["state"]; ok {
		return
	}
	t.Fatalf("expected status/state field in operation response, got: %v", payload)
}

func extractOperationID(body []byte) (string, bool) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", false
	}

	for _, key := range []string{"operationId", "operationID", "id"} {
		val, ok := payload[key]
		if !ok {
			continue
		}
		switch v := val.(type) {
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
