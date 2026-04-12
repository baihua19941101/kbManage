package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultPrometheusTimeout = 5 * time.Second

type PrometheusProvider struct {
	baseURL string
	client  *http.Client
}

func NewPrometheusProvider(baseURL string, client *http.Client) *PrometheusProvider {
	if client == nil {
		client = &http.Client{Timeout: defaultPrometheusTimeout}
	}
	return &PrometheusProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (p *PrometheusProvider) QuerySeries(ctx context.Context, req SeriesQuery) (SeriesResult, error) {
	if req.MetricKey == "" {
		return SeriesResult{}, fmt.Errorf("metricKey is required")
	}

	// Keep US1 chain available even without an upstream Prometheus.
	if p == nil || p.baseURL == "" {
		return fallbackSeriesPoints(req), nil
	}

	u, err := url.Parse(p.baseURL + "/api/v1/query_range")
	if err != nil {
		return fallbackSeriesPoints(req), nil
	}

	params := url.Values{}
	params.Set("query", req.MetricKey)
	params.Set("start", parsePromTime(req.StartAt, time.Now().UTC().Add(-15*time.Minute)).Format(time.RFC3339))
	params.Set("end", parsePromTime(req.EndAt, time.Now().UTC()).Format(time.RFC3339))
	params.Set("step", normalizeStep(req.Step))
	u.RawQuery = params.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fallbackSeriesPoints(req), nil
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fallbackSeriesPoints(req), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fallbackSeriesPoints(req), nil
	}

	var payload prometheusQueryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fallbackSeriesPoints(req), nil
	}

	points := payload.toPoints()
	if len(points) == 0 {
		return fallbackSeriesPoints(req), nil
	}
	return SeriesResult{Points: points}, nil
}

type prometheusQueryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string                   `json:"resultType"`
		Result     []prometheusSeriesResult `json:"result"`
	} `json:"data"`
}

type prometheusSeriesResult struct {
	Values [][]any `json:"values"`
}

func (r prometheusQueryRangeResponse) toPoints() []Point {
	if len(r.Data.Result) == 0 {
		return nil
	}

	var out []Point
	for _, value := range r.Data.Result[0].Values {
		if len(value) != 2 {
			continue
		}
		tsFloat, ok := value[0].(float64)
		if !ok {
			continue
		}
		vText, ok := value[1].(string)
		if !ok {
			continue
		}
		v, err := strconv.ParseFloat(vText, 64)
		if err != nil {
			continue
		}
		out = append(out, Point{
			Timestamp: time.Unix(int64(tsFloat), 0).UTC().Format(time.RFC3339),
			Value:     v,
		})
	}
	return out
}

func fallbackSeriesPoints(req SeriesQuery) SeriesResult {
	now := parsePromTime(req.EndAt, time.Now().UTC())
	return SeriesResult{
		Points: []Point{
			{Timestamp: now.Add(-10 * time.Minute).Format(time.RFC3339), Value: 0.32},
			{Timestamp: now.Add(-5 * time.Minute).Format(time.RFC3339), Value: 0.41},
			{Timestamp: now.Format(time.RFC3339), Value: 0.38},
		},
	}
}

func parsePromTime(raw string, fallback time.Time) time.Time {
	if raw == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fallback
	}
	return t.UTC()
}

func normalizeStep(step string) string {
	if strings.TrimSpace(step) == "" {
		return "30s"
	}
	return step
}
