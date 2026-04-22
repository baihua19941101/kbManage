package integration_test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestEnterprisePolishIntegration_ReportGenerationFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	seedEnterpriseIntegrationData(t, ctx)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/enterprise/reports", `{"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,"reportType":"management","title":"管理汇报","audienceType":"leadership","timeRange":"7d"}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create report status=%d body=%s", resp.Code, resp.Body.String())
	}
}
