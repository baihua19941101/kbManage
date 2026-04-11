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
	"time"
)

func TestOperationsContract_SubmitQueryAndFailureContract(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "operations-contract-user",
		Password:    "Operations@123",
		DisplayName: "Operations Contract User",
		Email:       "operations-contract-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	t.Run("submit operation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
			"clusterId": 1,
			"resourceKind": "Deployment",
			"operationType": "restart",
			"namespace": "default",
			"name": "api-server",
			"riskConfirmed": true
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusAccepted {
			t.Fatalf("expected submit operation status=202, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeOperationsContractObject(t, resp.Body.Bytes())
		operationID := extractOperationID(resp.Body.Bytes())
		if operationID == "" {
			t.Fatalf("expected submit operation response contains numeric id, payload=%v", payload)
		}
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "status")); got != "pending" {
			t.Fatalf("expected submit operation status=pending, got status=%q payload=%v", got, payload)
		}
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "operationType")); got != "restart" {
			t.Fatalf("expected operationType=restart, got operationType=%q payload=%v", got, payload)
		}
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "riskLevel")); got != "medium" {
			t.Fatalf("expected riskLevel=medium for restart, got riskLevel=%q payload=%v", got, payload)
		}
		expectedTargetRef := "cluster:1/ns:default/kind:Deployment/name:api-server"
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "targetRef")); got != expectedTargetRef {
			t.Fatalf("expected targetRef=%q, got targetRef=%q payload=%v", expectedTargetRef, got, payload)
		}
		if resultMessage := strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "resultMessage")); resultMessage != "" {
			t.Fatalf("expected initial resultMessage to be empty, got resultMessage=%q payload=%v", resultMessage, payload)
		}
	})

	t.Run("query operation status", func(t *testing.T) {
		createReq := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
			"clusterId": 1,
			"resourceKind": "Deployment",
			"operationType": "scale",
			"namespace": "default",
			"name": "api-server",
			"riskConfirmed": true
		}`))
		createReq.Header.Set("Authorization", "Bearer "+token)
		createReq.Header.Set("Content-Type", "application/json")

		createResp := httptest.NewRecorder()
		app.Router.ServeHTTP(createResp, createReq)
		if createResp.Code != http.StatusAccepted {
			t.Fatalf("failed to create operation fixture: status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
		}
		operationID := extractOperationID(createResp.Body.Bytes())
		if operationID == "" {
			t.Fatalf("operation fixture missing id, body=%s", strings.TrimSpace(createResp.Body.String()))
		}

		var finalPayload map[string]any
		var finalStatus string
		for i := 0; i < 15; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/operations/"+operationID, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp := httptest.NewRecorder()
			app.Router.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Fatalf("expected query operation status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}

			finalPayload = mustDecodeOperationsContractObject(t, resp.Body.Bytes())
			finalStatus = strings.TrimSpace(mustDecodeOperationsContractStringField(t, finalPayload, "status"))
			if finalStatus == "succeeded" || finalStatus == "failed" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if finalStatus != "succeeded" {
			t.Fatalf("expected operation status to reach succeeded, got status=%q payload=%v", finalStatus, finalPayload)
		}
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, finalPayload, "operationType")); got != "scale" {
			t.Fatalf("expected queried operationType=scale, got operationType=%q payload=%v", got, finalPayload)
		}
		if got := strings.TrimSpace(mustDecodeOperationsContractStringField(t, finalPayload, "riskLevel")); got != "medium" {
			t.Fatalf("expected queried riskLevel=medium, got riskLevel=%q payload=%v", got, finalPayload)
		}
		targetRef := strings.TrimSpace(mustDecodeOperationsContractStringField(t, finalPayload, "targetRef"))
		if !strings.Contains(targetRef, "cluster:1") || !strings.Contains(targetRef, "kind:Deployment") || !strings.Contains(targetRef, "name:api-server") {
			t.Fatalf("expected queried targetRef to include cluster/kind/name, got targetRef=%q payload=%v", targetRef, finalPayload)
		}
		resultMessage := strings.TrimSpace(mustDecodeOperationsContractStringField(t, finalPayload, "resultMessage"))
		if !strings.Contains(resultMessage, "operation executed") || !strings.Contains(resultMessage, "name:api-server") {
			t.Fatalf("expected semantic success resultMessage, got resultMessage=%q payload=%v", resultMessage, finalPayload)
		}
	})

	t.Run("submit invalid payload should fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{"clusterId":""}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code >= 200 && resp.Code < 300 {
			t.Fatalf("expected invalid payload to fail, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid payload status=400, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeOperationsContractObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "clusterid") || !strings.Contains(errMsg, "required") {
			t.Fatalf("expected invalid clusterId semantic error, got error=%q payload=%v", errMsg, payload)
		}
	})

	t.Run("submit high risk operation without confirmation should fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
			"clusterId": 1,
			"resourceKind": "Deployment",
			"operationType": "delete",
			"namespace": "default",
			"name": "api-server",
			"riskConfirmed": false
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected high-risk unconfirmed submission status=400, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeOperationsContractObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeOperationsContractStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "risk confirmation") || !strings.Contains(errMsg, "required") {
			t.Fatalf("expected risk-confirmation semantic error, got error=%q payload=%v", errMsg, payload)
		}
	})
}

func mustDecodeOperationsContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()

	if len(body) == 0 {
		t.Fatalf("response body is empty")
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	return payload
}

func mustDecodeOperationsContractStringField(t *testing.T, payload map[string]any, field string) string {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	text, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q should be string, got type=%T value=%v", field, raw, raw)
	}
	return text
}

func extractOperationID(body []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	val, ok := payload["id"]
	if !ok {
		return ""
	}
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		if v <= 0 || v != math.Trunc(v) {
			return ""
		}
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}
