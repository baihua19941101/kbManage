package contract_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/tests/testutil"

	"gorm.io/gorm"
)

func TestAuditContract_QueryAndExportContract(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "audit-contract-user",
		Password:    "Audit@123456",
		DisplayName: "Audit Contract User",
		Email:       "audit-contract-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	startAt := mustParseAuditRFC3339(t, "2026-01-01T00:00:00Z")
	endAt := mustParseAuditRFC3339(t, "2026-01-31T23:59:59Z")
	inRangeAt := mustParseAuditRFC3339(t, "2026-01-10T09:30:00Z")
	outOfRangeAt := mustParseAuditRFC3339(t, "2026-02-10T09:30:00Z")
	otherActorID := seeded.User.ID + 99

	createAuditFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-contract-match",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: inRangeAt,
	})
	createAuditFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-contract-wrong-outcome",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeDenied,
		CreatedAt: inRangeAt,
	})
	createAuditFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-contract-wrong-actor",
		ActorID:   &otherActorID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: inRangeAt,
	})
	createAuditFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-contract-out-of-range",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: outOfRangeAt,
	})

	t.Run("query audit events", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/audits/events?startAt=2026-01-01T00:00:00Z&endAt=2026-01-31T23:59:59Z&actorId="+strconv.FormatUint(seeded.User.ID, 10)+"&action=operation.execute&outcome=success",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected query audits status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeAuditContractObject(t, resp.Body.Bytes())
		items := mustDecodeAuditContractArrayField(t, payload, "items")
		count := mustDecodeAuditContractIntField(t, payload, "count")
		if count != 1 || len(items) != 1 {
			t.Fatalf("expected exactly one filtered audit event, got count=%d items=%d payload=%v", count, len(items), payload)
		}

		item := mustDecodeAuditContractObjectFromAny(t, items[0], "items[0]")
		if got := strings.TrimSpace(mustDecodeAuditContractStringField(t, item, "Action")); got != "operation.execute" {
			t.Fatalf("expected filtered event action=operation.execute, got action=%q item=%v", got, item)
		}
		if got := strings.TrimSpace(mustDecodeAuditContractStringField(t, item, "Outcome")); got != string(domain.AuditOutcomeSuccess) {
			t.Fatalf("expected filtered event outcome=success, got outcome=%q item=%v", got, item)
		}
		if got := mustDecodeAuditContractUint64Field(t, item, "ActorID"); got != seeded.User.ID {
			t.Fatalf("expected filtered event actorId=%d, got actorId=%d item=%v", seeded.User.ID, got, item)
		}
		createdAt := mustDecodeAuditContractTimeField(t, item, "CreatedAt")
		if createdAt.Before(startAt) || createdAt.After(endAt) {
			t.Fatalf("expected filtered event createdAt in [%s, %s], got createdAt=%s item=%v", startAt.Format(time.RFC3339), endAt.Format(time.RFC3339), createdAt.Format(time.RFC3339), item)
		}
	})

	t.Run("export audit events", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/audits/exports", strings.NewReader(`{
			"startAt":"2026-01-01T00:00:00Z",
			"endAt":"2026-01-31T23:59:59Z",
			"actorId":`+strconv.FormatUint(seeded.User.ID, 10)+`,
			"action":"operation.execute",
			"outcome":"success"
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusAccepted {
			t.Fatalf("expected export audits status=202, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}

		payload := mustDecodeAuditContractObject(t, resp.Body.Bytes())
		taskID := strings.TrimSpace(mustDecodeAuditContractStringField(t, payload, "taskId"))
		if !strings.HasPrefix(taskID, "aexp-") {
			t.Fatalf("expected taskId with prefix aexp-, got taskId=%q payload=%v", taskID, payload)
		}
		if status := strings.TrimSpace(mustDecodeAuditContractStringField(t, payload, "status")); status != "pending" {
			t.Fatalf("expected export submit status=pending, got status=%q payload=%v", status, payload)
		}
		if operatorID := mustDecodeAuditContractUint64Field(t, payload, "operatorId"); operatorID != seeded.User.ID {
			t.Fatalf("expected operatorId=%d, got operatorId=%d payload=%v", seeded.User.ID, operatorID, payload)
		}
		if resultTotal := mustDecodeAuditContractIntField(t, payload, "resultTotal"); resultTotal != 0 {
			t.Fatalf("expected initial resultTotal=0, got resultTotal=%d payload=%v", resultTotal, payload)
		}

		var finalPayload map[string]any
		var finalStatus string
		for i := 0; i < 15; i++ {
			statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/audits/exports/"+taskID, nil)
			statusReq.Header.Set("Authorization", "Bearer "+token)

			statusResp := httptest.NewRecorder()
			app.Router.ServeHTTP(statusResp, statusReq)
			if statusResp.Code != http.StatusOK {
				t.Fatalf("expected export status query status=200, got status=%d body=%s", statusResp.Code, strings.TrimSpace(statusResp.Body.String()))
			}

			finalPayload = mustDecodeAuditContractObject(t, statusResp.Body.Bytes())
			finalStatus = strings.TrimSpace(mustDecodeAuditContractStringField(t, finalPayload, "status"))
			if finalStatus == "succeeded" || finalStatus == "failed" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}

		if finalStatus != "succeeded" {
			t.Fatalf("expected export status to reach succeeded, got status=%q payload=%v", finalStatus, finalPayload)
		}
		if resultTotal := mustDecodeAuditContractIntField(t, finalPayload, "resultTotal"); resultTotal != 1 {
			t.Fatalf("expected succeeded export resultTotal=1, got resultTotal=%d payload=%v", resultTotal, finalPayload)
		}
		downloadURL := strings.TrimSpace(mustDecodeAuditContractStringField(t, finalPayload, "downloadUrl"))
		if !strings.Contains(downloadURL, taskID) {
			t.Fatalf("expected downloadUrl to include task id, taskId=%q downloadUrl=%q payload=%v", taskID, downloadURL, finalPayload)
		}
		if errorMessage := strings.TrimSpace(mustDecodeAuditContractStringField(t, finalPayload, "errorMessage")); errorMessage != "" {
			t.Fatalf("expected errorMessage empty for succeeded export, got errorMessage=%q payload=%v", errorMessage, finalPayload)
		}
	})

	t.Run("query with invalid time range should return semantic error", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/audits/events?startAt=2026-01-31T23:59:59Z&endAt=2026-01-01T00:00:00Z",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("expected invalid time range status=500 per current contract, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeAuditContractObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeAuditContractStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "startat") || !strings.Contains(errMsg, "earlier") {
			t.Fatalf("expected semantic time-range error, got error=%q payload=%v", errMsg, payload)
		}
	})
}

func createAuditFixtureEvent(t *testing.T, db *gorm.DB, event domain.AuditEvent) {
	t.Helper()

	if db == nil {
		t.Fatal("audit fixture requires non-nil db")
	}
	if err := db.WithContext(context.Background()).Create(&event).Error; err != nil {
		t.Fatalf("seed audit fixture failed: %v", err)
	}
}

func mustParseAuditRFC3339(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("invalid RFC3339 fixture %q: %v", value, err)
	}
	return parsed
}

func mustDecodeAuditContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()

	if len(body) == 0 {
		t.Fatal("response body is empty")
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	return payload
}

func mustDecodeAuditContractObjectFromAny(t *testing.T, value any, field string) map[string]any {
	t.Helper()

	payload, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s is not a JSON object, got type=%T value=%v", field, value, value)
	}
	return payload
}

func mustDecodeAuditContractArrayField(t *testing.T, payload map[string]any, field string) []any {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q is not JSON array, got type=%T value=%v", field, raw, raw)
	}
	return items
}

func mustDecodeAuditContractStringField(t *testing.T, payload map[string]any, field string) string {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	text, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q is not string, got type=%T value=%v", field, raw, raw)
	}
	return text
}

func mustDecodeAuditContractIntField(t *testing.T, payload map[string]any, field string) int {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	number, ok := raw.(float64)
	if !ok {
		t.Fatalf("field %q is not JSON number, got type=%T value=%v", field, raw, raw)
	}
	if number < 0 {
		t.Fatalf("field %q should not be negative, got value=%v", field, number)
	}
	return int(number)
}

func mustDecodeAuditContractUint64Field(t *testing.T, payload map[string]any, field string) uint64 {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	switch v := raw.(type) {
	case float64:
		if v < 0 {
			t.Fatalf("field %q should not be negative, got value=%v", field, v)
		}
		return uint64(v)
	case string:
		parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
		if err != nil {
			t.Fatalf("field %q should be uint64 string, got value=%q err=%v", field, v, err)
		}
		return parsed
	default:
		t.Fatalf("field %q has unsupported type=%T value=%v", field, raw, raw)
		return 0
	}
}

func mustDecodeAuditContractTimeField(t *testing.T, payload map[string]any, field string) time.Time {
	t.Helper()

	text := strings.TrimSpace(mustDecodeAuditContractStringField(t, payload, field))
	parsed, err := time.Parse(time.RFC3339Nano, text)
	if err != nil {
		t.Fatalf("field %q should be RFC3339 timestamp, got value=%q err=%v", field, text, err)
	}
	return parsed
}
