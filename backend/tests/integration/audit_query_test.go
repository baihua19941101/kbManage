package integration_test

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

func TestAuditQuery_SearchAndExportContract(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "audit-integration-user",
		Password: "Audit@123456",
		Email:    "audit-integration-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)
	otherActorID := seeded.User.ID + 42

	inRangeEarly := mustParseAuditQueryRFC3339(t, "2026-02-10T09:00:00Z")
	inRangeLate := mustParseAuditQueryRFC3339(t, "2026-02-10T09:05:00Z")
	outRange := mustParseAuditQueryRFC3339(t, "2026-03-01T00:00:00Z")

	createAuditQueryFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-query-match-early",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: inRangeEarly,
	})
	createAuditQueryFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-query-match-late",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: inRangeLate,
	})
	createAuditQueryFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-query-wrong-outcome",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeFailed,
		CreatedAt: inRangeLate,
	})
	createAuditQueryFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-query-wrong-actor",
		ActorID:   &otherActorID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: inRangeLate,
	})
	createAuditQueryFixtureEvent(t, app.DB, domain.AuditEvent{
		RequestID: "audit-query-out-range",
		ActorID:   &seeded.User.ID,
		Action:    "operation.execute",
		Outcome:   domain.AuditOutcomeSuccess,
		CreatedAt: outRange,
	})

	t.Run("search audit events", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/audits/events?startAt=2026-02-01T00:00:00Z&endAt=2026-02-28T23:59:59Z&actorId="+strconv.FormatUint(seeded.User.ID, 10)+"&action=operation.execute&outcome=success&limit=2",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected search audit events status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeAuditQueryObject(t, resp.Body.Bytes())
		items := mustDecodeAuditQueryArrayField(t, payload, "items")
		count := mustDecodeAuditQueryIntField(t, payload, "count")
		if count != 2 || len(items) != 2 {
			t.Fatalf("expected exactly two filtered audit events, got count=%d items=%d payload=%v", count, len(items), payload)
		}

		first := mustDecodeAuditQueryObjectFromAny(t, items[0], "items[0]")
		second := mustDecodeAuditQueryObjectFromAny(t, items[1], "items[1]")
		if strings.TrimSpace(mustDecodeAuditQueryStringField(t, first, "Outcome")) != string(domain.AuditOutcomeSuccess) {
			t.Fatalf("first item outcome should be success, item=%v", first)
		}
		if strings.TrimSpace(mustDecodeAuditQueryStringField(t, second, "Outcome")) != string(domain.AuditOutcomeSuccess) {
			t.Fatalf("second item outcome should be success, item=%v", second)
		}
		if strings.TrimSpace(mustDecodeAuditQueryStringField(t, first, "Action")) != "operation.execute" {
			t.Fatalf("first item action should be operation.execute, item=%v", first)
		}
		if strings.TrimSpace(mustDecodeAuditQueryStringField(t, second, "Action")) != "operation.execute" {
			t.Fatalf("second item action should be operation.execute, item=%v", second)
		}
		if mustDecodeAuditQueryUint64Field(t, first, "ActorID") != seeded.User.ID || mustDecodeAuditQueryUint64Field(t, second, "ActorID") != seeded.User.ID {
			t.Fatalf("all returned items should belong to actorId=%d, items=%v", seeded.User.ID, items)
		}
		firstAt := mustDecodeAuditQueryTimeField(t, first, "CreatedAt")
		secondAt := mustDecodeAuditQueryTimeField(t, second, "CreatedAt")
		if firstAt.Before(secondAt) {
			t.Fatalf("expected audits sorted by createdAt desc, first=%s second=%s", firstAt.Format(time.RFC3339), secondAt.Format(time.RFC3339))
		}
	})

	t.Run("export audit events", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/audits/exports", strings.NewReader(`{
			"startAt":"2026-02-01T00:00:00Z",
			"endAt":"2026-02-28T23:59:59Z",
			"actorId":`+strconv.FormatUint(seeded.User.ID, 10)+`,
			"action":"operation.execute",
			"outcome":"success"
		}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusAccepted {
			t.Fatalf("expected export audit events status=202, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		submitPayload := mustDecodeAuditQueryObject(t, resp.Body.Bytes())
		taskID := strings.TrimSpace(mustDecodeAuditQueryStringField(t, submitPayload, "taskId"))
		if !strings.HasPrefix(taskID, "aexp-") {
			t.Fatalf("expected export taskId with aexp- prefix, got taskId=%q payload=%v", taskID, submitPayload)
		}
		if strings.TrimSpace(mustDecodeAuditQueryStringField(t, submitPayload, "status")) != "pending" {
			t.Fatalf("expected submit export status=pending, payload=%v", submitPayload)
		}
		if mustDecodeAuditQueryUint64Field(t, submitPayload, "operatorId") != seeded.User.ID {
			t.Fatalf("expected submit export operatorId=%d payload=%v", seeded.User.ID, submitPayload)
		}

		var finalPayload map[string]any
		var lastStatus string
		for i := 0; i < 15; i++ {
			statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/audits/exports/"+taskID, nil)
			statusReq.Header.Set("Authorization", "Bearer "+token)

			statusResp := httptest.NewRecorder()
			app.Router.ServeHTTP(statusResp, statusReq)
			if statusResp.Code != http.StatusOK {
				t.Fatalf("expected get export status=200, got status=%d body=%s", statusResp.Code, strings.TrimSpace(statusResp.Body.String()))
			}

			finalPayload = mustDecodeAuditQueryObject(t, statusResp.Body.Bytes())
			lastStatus = strings.TrimSpace(mustDecodeAuditQueryStringField(t, finalPayload, "status"))
			if lastStatus == "succeeded" || lastStatus == "failed" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}

		if lastStatus != "succeeded" {
			t.Fatalf("expected audit export final status=succeeded, got status=%q payload=%v", lastStatus, finalPayload)
		}
		if total := mustDecodeAuditQueryIntField(t, finalPayload, "resultTotal"); total != 2 {
			t.Fatalf("expected export resultTotal=2, got total=%d payload=%v", total, finalPayload)
		}
		downloadURL := strings.TrimSpace(mustDecodeAuditQueryStringField(t, finalPayload, "downloadUrl"))
		if !strings.Contains(downloadURL, taskID) {
			t.Fatalf("expected downloadUrl to contain task id, taskId=%q downloadUrl=%q payload=%v", taskID, downloadURL, finalPayload)
		}
		if errMsg := strings.TrimSpace(mustDecodeAuditQueryStringField(t, finalPayload, "errorMessage")); errMsg != "" {
			t.Fatalf("expected successful export errorMessage empty, got errorMessage=%q payload=%v", errMsg, finalPayload)
		}
		createdAt := mustDecodeAuditQueryTimeField(t, finalPayload, "createdAt")
		updatedAt := mustDecodeAuditQueryTimeField(t, finalPayload, "updatedAt")
		if updatedAt.Before(createdAt) {
			t.Fatalf("expected updatedAt >= createdAt, createdAt=%s updatedAt=%s payload=%v", createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339), finalPayload)
		}
		completedAt := mustDecodeAuditQueryTimeField(t, finalPayload, "completedAt")
		if completedAt.Before(createdAt) {
			t.Fatalf("expected completedAt >= createdAt, createdAt=%s completedAt=%s payload=%v", createdAt.Format(time.RFC3339), completedAt.Format(time.RFC3339), finalPayload)
		}
	})

	t.Run("query unknown export task should return 404 with error message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/audits/exports/aexp-not-exist", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		app.Router.ServeHTTP(resp, req)

		if resp.Code != http.StatusNotFound {
			t.Fatalf("expected unknown export task status=404, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		payload := mustDecodeAuditQueryObject(t, resp.Body.Bytes())
		errMsg := strings.ToLower(strings.TrimSpace(mustDecodeAuditQueryStringField(t, payload, "error")))
		if !strings.Contains(errMsg, "not found") && !strings.Contains(errMsg, "record") {
			t.Fatalf("expected semantic not-found error, got error=%q payload=%v", errMsg, payload)
		}
	})
}

func createAuditQueryFixtureEvent(t *testing.T, db *gorm.DB, event domain.AuditEvent) {
	t.Helper()

	if db == nil {
		t.Fatal("audit fixture requires non-nil db")
	}
	if err := db.WithContext(context.Background()).Create(&event).Error; err != nil {
		t.Fatalf("seed audit fixture failed: %v", err)
	}
}

func mustParseAuditQueryRFC3339(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("invalid RFC3339 fixture %q: %v", value, err)
	}
	return parsed
}

func mustDecodeAuditQueryObject(t *testing.T, body []byte) map[string]any {
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

func mustDecodeAuditQueryObjectFromAny(t *testing.T, value any, field string) map[string]any {
	t.Helper()

	payload, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s is not JSON object, got type=%T value=%v", field, value, value)
	}
	return payload
}

func mustDecodeAuditQueryArrayField(t *testing.T, payload map[string]any, field string) []any {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q is not array, got type=%T value=%v", field, raw, raw)
	}
	return items
}

func mustDecodeAuditQueryStringField(t *testing.T, payload map[string]any, field string) string {
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

func mustDecodeAuditQueryIntField(t *testing.T, payload map[string]any, field string) int {
	t.Helper()

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response does not include %q field: %v", field, payload)
	}
	number, ok := raw.(float64)
	if !ok {
		t.Fatalf("field %q is not number, got type=%T value=%v", field, raw, raw)
	}
	if number < 0 {
		t.Fatalf("field %q should not be negative, got value=%v", field, number)
	}
	return int(number)
}

func mustDecodeAuditQueryUint64Field(t *testing.T, payload map[string]any, field string) uint64 {
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

func mustDecodeAuditQueryTimeField(t *testing.T, payload map[string]any, field string) time.Time {
	t.Helper()

	text := strings.TrimSpace(mustDecodeAuditQueryStringField(t, payload, field))
	parsed, err := time.Parse(time.RFC3339Nano, text)
	if err != nil {
		t.Fatalf("field %q should be RFC3339 timestamp, got value=%q err=%v", field, text, err)
	}
	return parsed
}
