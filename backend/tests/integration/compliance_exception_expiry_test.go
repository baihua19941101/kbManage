package integration_test

import (
	"context"
	"testing"
	"time"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

func TestComplianceExceptionExpiryIntegration_ExpiresDueRequests(t *testing.T) {
	t.Parallel()
	svc := complianceSvc.NewExceptionService()
	_, err := svc.CreateException(context.Background(), 1, "finding-1", complianceSvc.CreateExceptionInput{Reason: "legacy", StartsAt: time.Now().Add(-2 * time.Hour), ExpiresAt: time.Now().Add(-time.Hour)})
	if err != nil {
		t.Fatalf("create exception failed: %v", err)
	}
	items, err := svc.ListExceptions(context.Background(), complianceSvc.ExceptionFilter{})
	if err != nil || len(items) == 0 {
		t.Fatalf("list exception failed err=%v", err)
	}
	_, err = svc.ReviewException(context.Background(), 1, items[0].ID, complianceSvc.ReviewExceptionInput{Decision: "approve", ReviewComment: "approve for expiry"})
	if err != nil {
		t.Fatalf("approve exception failed: %v", err)
	}
	count, err := svc.ExpireDueExceptions(context.Background(), time.Now())
	if err != nil || count == 0 {
		t.Fatalf("expire due exception failed count=%d err=%v", count, err)
	}
}
