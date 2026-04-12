package alerts

import "context"

// AlertmanagerProvider is a placeholder adapter for Alertmanager-compatible backends.
// Current implementation reuses mock data until real upstream integration is wired.
type AlertmanagerProvider struct {
	baseURL string
}

func NewAlertmanagerProvider(baseURL string) *AlertmanagerProvider {
	return &AlertmanagerProvider{baseURL: baseURL}
}

func (p *AlertmanagerProvider) ListAlerts(ctx context.Context, req AlertQuery) ([]AlertItem, error) {
	_ = p.baseURL
	return NewMockProvider().ListAlerts(ctx, req)
}
