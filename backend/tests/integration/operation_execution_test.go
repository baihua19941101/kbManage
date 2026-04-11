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

	"kbmanage/backend/tests/testutil"
)

func TestOperationExecution_SubmitTrackAndFailureContract(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "operation-integration-user",
		Password: "Operation@123456",
		Email:    "operation-integration-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	operationID := submitOperationFixture(t, app.Router, token)
	if operationID == "" {
		t.Fatal("submit operation fixture returned empty id")
	}

	t.Run("query operation status", func(t *testing.T) {
		statusRank := map[string]int{
			"pending":   1,
			"running":   2,
			"succeeded": 3,
			"failed":    3,
		}
		observed := make([]string, 0, 8)
		var lastBody string
		var finalPayload map[string]any
		var finalStatus string
		lastRank := 0

		for i := 0; i < 15; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/operations/"+operationID, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp := httptest.NewRecorder()
			app.Router.ServeHTTP(resp, req)
			lastBody = strings.TrimSpace(resp.Body.String())

			if resp.Code != http.StatusOK {
				t.Fatalf("expected query operation status=200, got status=%d body=%s", resp.Code, lastBody)
			}

			finalPayload = mustDecodeOperationExecutionObject(t, resp.Body.Bytes())
			finalStatus = strings.TrimSpace(mustDecodeOperationExecutionStringField(t, finalPayload, "status"))
			rank, ok := statusRank[finalStatus]
			if !ok {
				t.Fatalf("unexpected operation status=%q payload=%v", finalStatus, finalPayload)
			}
			if rank < lastRank {
				t.Fatalf("operation status regressed, previousRank=%d currentStatus=%q observed=%v payload=%v", lastRank, finalStatus, observed, finalPayload)
			}
			lastRank = rank
			observed = append(observed, finalStatus)

			idInResp := extractOperationID(resp.Body.Bytes())
			if idInResp != operationID {
				t.Fatalf("expected queried operation id=%s, got id=%s payload=%v", operationID, idInResp, finalPayload)
			}
			if opType := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, finalPayload, "operationType")); opType != "scale" {
				t.Fatalf("expected queried operationType=scale, got operationType=%q payload=%v", opType, finalPayload)
			}

			targetRef := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, finalPayload, "targetRef"))
			if !strings.Contains(targetRef, "cluster:1") || !strings.Contains(targetRef, "ns:default") || !strings.Contains(targetRef, "name:api-server") {
				t.Fatalf("expected targetRef contains cluster/ns/name, got targetRef=%q payload=%v", targetRef, finalPayload)
			}

			if finalStatus == "succeeded" || finalStatus == "failed" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}

		if len(observed) == 0 {
			t.Fatalf("operation status polling produced no observations, last body=%s", lastBody)
		}
		if finalStatus != "succeeded" {
			t.Fatalf("expected final operation status=succeeded, got status=%q observed=%v payload=%v", finalStatus, observed, finalPayload)
		}

		resultMessage := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, finalPayload, "resultMessage"))
		if !strings.Contains(resultMessage, "operation executed") || !strings.Contains(resultMessage, "name:api-server") {
			t.Fatalf("expected semantic success resultMessage, got resultMessage=%q payload=%v", resultMessage, finalPayload)
		}
	})

	t.Run("submit missing operation type should fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
			"clusterId": 1,
			"resourceKind": "Deployment",
			"operationType": "",
			"namespace": "default",
			"name": "api-server",
			"riskConfirmed": true
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected missing operationType status=400, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeOperationExecutionObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "operationtype") || !strings.Contains(errMsg, "required") {
			t.Fatalf("expected operationType semantic error, got error=%q payload=%v", errMsg, payload)
		}
	})

	t.Run("submit high-risk operation without confirmation should fail", func(t *testing.T) {
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
			t.Fatalf("expected unconfirmed high-risk operation status=400, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeOperationExecutionObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "risk confirmation") || !strings.Contains(errMsg, "required") {
			t.Fatalf("expected risk confirmation semantic error, got error=%q payload=%v", errMsg, payload)
		}
	})
}

func submitOperationFixture(t *testing.T, r http.Handler, token string) string {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/operations", strings.NewReader(`{
		"clusterId": 1,
		"resourceKind": "Deployment",
		"operationType": "scale",
		"namespace": "default",
		"name": "api-server",
		"riskConfirmed": true
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusAccepted {
		t.Fatalf("expected submit operation status=202, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeOperationExecutionObject(t, resp.Body.Bytes())
	if status := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "status")); status != "pending" {
		t.Fatalf("expected fixture submit status=pending, got status=%q payload=%v", status, payload)
	}
	if opType := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "operationType")); opType != "scale" {
		t.Fatalf("expected fixture operationType=scale, got operationType=%q payload=%v", opType, payload)
	}
	if risk := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "riskLevel")); risk != "medium" {
		t.Fatalf("expected fixture riskLevel=medium, got riskLevel=%q payload=%v", risk, payload)
	}
	if targetRef := strings.TrimSpace(mustDecodeOperationExecutionStringField(t, payload, "targetRef")); !strings.Contains(targetRef, "kind:Deployment") {
		t.Fatalf("expected fixture targetRef includes resource kind, got targetRef=%q payload=%v", targetRef, payload)
	}

	return extractOperationID(resp.Body.Bytes())
}

func mustDecodeOperationExecutionObject(t *testing.T, body []byte) map[string]any {
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

func mustDecodeOperationExecutionStringField(t *testing.T, payload map[string]any, field string) string {
	t.Helper()

	val, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	v, _ := val.(string)
	if strings.TrimSpace(v) == "" && field != "resultMessage" {
		t.Fatalf("field %q should not be empty string, payload=%v", field, payload)
	}
	return v
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
