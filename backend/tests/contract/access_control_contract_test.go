package contract_test

import (
	"encoding/json"
	"kbmanage/backend/tests/testutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestAccessControlContract_WorkspaceProjectRoleBindingRoutes(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "access-contract-admin",
		Password:    "Access@123456",
		DisplayName: "Access Contract Admin",
		Email:       "access-contract-admin@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createWorkspaceResp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodPost, "/api/v1/workspaces", `{
		"name":"team-alpha",
		"description":"contract test workspace"
	}`)
	if createWorkspaceResp.Code != http.StatusCreated {
		t.Fatalf("expected create workspace status=201, got status=%d body=%s", createWorkspaceResp.Code, strings.TrimSpace(createWorkspaceResp.Body.String()))
	}
	workspaceID := extractNumericIDField(t, createWorkspaceResp.Body.Bytes(), "id")

	t.Run("list workspaces", func(t *testing.T) {
		resp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodGet, "/api/v1/workspaces", "")
		if resp.Code != http.StatusOK {
			t.Fatalf("expected list workspaces status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		assertAccessControlJSON(t, resp.Body.Bytes())
		assertHasField(t, resp.Body.Bytes(), "items")
	})

	t.Run("create and list projects", func(t *testing.T) {
		createProjectResp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/projects", `{
			"name":"billing-api",
			"description":"project from contract test"
		}`)
		if createProjectResp.Code != http.StatusCreated {
			t.Fatalf("expected create project status=201, got status=%d body=%s", createProjectResp.Code, strings.TrimSpace(createProjectResp.Body.String()))
		}

		listProjectResp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/projects", "")
		if listProjectResp.Code != http.StatusOK {
			t.Fatalf("expected list projects status=200, got status=%d body=%s", listProjectResp.Code, strings.TrimSpace(listProjectResp.Body.String()))
		}
		assertAccessControlJSON(t, listProjectResp.Body.Bytes())
		assertHasField(t, listProjectResp.Body.Bytes(), "items")
	})

	t.Run("create and list role bindings", func(t *testing.T) {
		createBindingResp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodPost, "/api/v1/role-bindings", `{
			"subjectType":"user",
			"subjectId":1002,
			"scopeType":"workspace",
			"scopeId":`+workspaceID+`,
			"roleKey":"workspace-owner"
		}`)
		if createBindingResp.Code != http.StatusCreated {
			t.Fatalf("expected create role binding status=201, got status=%d body=%s", createBindingResp.Code, strings.TrimSpace(createBindingResp.Body.String()))
		}

		listBindingResp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodGet, "/api/v1/role-bindings?scopeType=workspace&scopeId="+workspaceID, "")
		if listBindingResp.Code != http.StatusOK {
			t.Fatalf("expected list role bindings status=200, got status=%d body=%s", listBindingResp.Code, strings.TrimSpace(listBindingResp.Body.String()))
		}
		assertAccessControlJSON(t, listBindingResp.Body.Bytes())
		assertHasField(t, listBindingResp.Body.Bytes(), "items")
	})

	t.Run("create role binding with unknown role key should fail", func(t *testing.T) {
		resp := performAccessAuthedJSONRequest(t, app.Router, token, http.MethodPost, "/api/v1/role-bindings", `{
			"subjectType":"user",
			"subjectId":1003,
			"scopeType":"workspace",
			"scopeId":`+workspaceID+`,
			"roleKey":"workspace-unknown-role"
		}`)
		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected unknown role key status=400, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		assertAccessControlJSON(t, resp.Body.Bytes())
		assertHasField(t, resp.Body.Bytes(), "error")
	})
}

func performAccessAuthedJSONRequest(t *testing.T, h http.Handler, token, method, target, body string) *httptest.ResponseRecorder {
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

func extractNumericIDField(t *testing.T, body []byte, field string) string {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	val, ok := payload[field]
	if !ok && strings.ToLower(field) == "id" {
		val, ok = payload["ID"]
	}
	if !ok {
		t.Fatalf("response does not include %q field", field)
	}

	switch v := val.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			t.Fatalf("field %q is empty string", field)
		}
		return v
	case float64:
		if v <= 0 || v != math.Trunc(v) {
			t.Fatalf("field %q is invalid number: %v", field, v)
		}
		return strconv.FormatInt(int64(v), 10)
	default:
		t.Fatalf("field %q has unsupported type %T", field, val)
		return ""
	}
}

func assertHasField(t *testing.T, body []byte, field string) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	if _, ok := payload[field]; !ok {
		t.Fatalf("expected field %q in response, got: %v", field, payload)
	}
}
