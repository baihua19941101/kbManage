package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestAuthContract_LoginAndRefresh(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "auth-contract-admin",
		Password:    "Admin@123456",
		DisplayName: "Auth Contract Admin",
		Email:       "auth-contract-admin@example.test",
	})

	t.Run("POST /api/v1/auth/login success", func(t *testing.T) {
		resp := performAuthJSONRequest(t, app.Router, http.MethodPost, "/api/v1/auth/login", `{"username":"`+seeded.User.Username+`","password":"`+seeded.Password+`"}`)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected login success status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}

		accessToken, refreshToken := assertAuthContractTokenPair(t, resp.Body.Bytes())
		if accessToken == "" || refreshToken == "" {
			t.Fatalf("expected non-empty token pair, body=%s", strings.TrimSpace(resp.Body.String()))
		}
	})

	t.Run("POST /api/v1/auth/login wrong password fails", func(t *testing.T) {
		resp := performAuthJSONRequest(t, app.Router, http.MethodPost, "/api/v1/auth/login", `{"username":"`+seeded.User.Username+`","password":"Admin@123456-wrong"}`)

		if resp.Code >= 200 && resp.Code < 300 {
			t.Fatalf("expected wrong password to fail, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}

		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("unexpected status for wrong-password login: %d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
	})

	t.Run("POST /api/v1/auth/refresh success", func(t *testing.T) {
		loginResp := performAuthJSONRequest(t, app.Router, http.MethodPost, "/api/v1/auth/login", `{"username":"`+seeded.User.Username+`","password":"`+seeded.Password+`"}`)
		if loginResp.Code != http.StatusOK {
			t.Fatalf("cannot prepare refresh fixture with login, status=%d body=%s", loginResp.Code, strings.TrimSpace(loginResp.Body.String()))
		}

		_, refreshToken := assertAuthContractTokenPair(t, loginResp.Body.Bytes())
		if refreshToken == "" {
			t.Fatalf("login response missing refresh token, body=%s", strings.TrimSpace(loginResp.Body.String()))
		}

		refreshResp := performAuthJSONRequest(t, app.Router, http.MethodPost, "/api/v1/auth/refresh", `{"refreshToken":"`+refreshToken+`"}`)
		if refreshResp.Code != http.StatusOK {
			t.Fatalf("expected refresh success status=200, got status=%d body=%s", refreshResp.Code, strings.TrimSpace(refreshResp.Body.String()))
		}

		accessToken, newRefreshToken := assertAuthContractTokenPair(t, refreshResp.Body.Bytes())
		if accessToken == "" || newRefreshToken == "" {
			t.Fatalf("expected refresh response token pair, body=%s", strings.TrimSpace(refreshResp.Body.String()))
		}
	})
}

func performAuthJSONRequest(t *testing.T, r http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func assertAuthContractTokenPair(t *testing.T, body []byte) (string, string) {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	accessToken, _ := payload["accessToken"].(string)
	refreshToken, _ := payload["refreshToken"].(string)

	if _, ok := payload["user"]; !ok {
		t.Fatalf("response does not include user field")
	}

	return accessToken, refreshToken
}
