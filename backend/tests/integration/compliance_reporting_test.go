package integration_test

import (
	"context"
	"testing"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

func TestComplianceReportingIntegration_OverviewAndTrendAvailable(t *testing.T) {
	t.Parallel()
	overviewSvc := complianceSvc.NewOverviewService()
	trendSvc := complianceSvc.NewTrendService()
	overview, err := overviewSvc.GetOverview(context.Background(), complianceSvc.OverviewFilter{})
	if err != nil || overview == nil {
		t.Fatalf("overview failed err=%v", err)
	}
	_, err = trendSvc.RecordSnapshot(context.Background(), complianceSvc.TrendFilter{})
	if err != nil {
		t.Fatalf("record trend snapshot failed: %v", err)
	}
	trends, err := trendSvc.GetTrends(context.Background(), complianceSvc.TrendFilter{})
	if err != nil || trends == nil {
		t.Fatalf("get trends failed err=%v", err)
	}
}
