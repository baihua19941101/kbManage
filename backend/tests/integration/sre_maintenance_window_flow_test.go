package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformSREIntegration_MaintenanceWindowFlow(t *testing.T) {
	ctx := newSREIntegrationCtx(t, "workspace-owner")
	workspaceID := strconv.FormatUint(ctx.Access.WorkspaceID, 10)
	resp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/sre/maintenance-windows", `{
		"workspaceId":`+workspaceID+`,
		"name":"integration-window",
		"windowType":"maintenance",
		"scope":"platform",
		"startAt":"2026-04-19T10:00:00Z",
		"endAt":"2026-04-19T12:00:00Z"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create maintenance window failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	listResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/maintenance-windows", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list maintenance windows failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
