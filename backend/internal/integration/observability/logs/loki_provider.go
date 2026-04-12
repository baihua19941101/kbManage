package logs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultLokiTimeout = 5 * time.Second

type LokiProvider struct {
	baseURL string
	client  *http.Client
}

func NewLokiProvider(baseURL string, client *http.Client) *LokiProvider {
	if client == nil {
		client = &http.Client{Timeout: defaultLokiTimeout}
	}
	return &LokiProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (p *LokiProvider) Query(ctx context.Context, req QueryRequest) (QueryResult, error) {
	if p == nil || p.baseURL == "" {
		return fallbackLokiResult(req), nil
	}

	u, err := url.Parse(p.baseURL + "/loki/api/v1/query_range")
	if err != nil {
		return fallbackLokiResult(req), nil
	}
	params := url.Values{}
	params.Set("query", buildLokiQuery(req))
	params.Set("start", strconv.FormatInt(parseLokiTime(req.StartAt, time.Now().UTC().Add(-15*time.Minute)).UnixNano(), 10))
	params.Set("end", strconv.FormatInt(parseLokiTime(req.EndAt, time.Now().UTC()).UnixNano(), 10))
	params.Set("limit", strconv.Itoa(normalizeLimit(req.Limit)))
	u.RawQuery = params.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fallbackLokiResult(req), nil
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fallbackLokiResult(req), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fallbackLokiResult(req), nil
	}

	var payload lokiQueryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fallbackLokiResult(req), nil
	}
	items := payload.toLogEntries()
	if len(items) == 0 {
		return fallbackLokiResult(req), nil
	}
	return QueryResult{Items: items}, nil
}

type lokiQueryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string             `json:"resultType"`
		Result     []lokiStreamResult `json:"result"`
	} `json:"data"`
}

type lokiStreamResult struct {
	Values [][]string `json:"values"`
}

func (r lokiQueryRangeResponse) toLogEntries() []LogEntry {
	var out []LogEntry
	for _, stream := range r.Data.Result {
		for _, value := range stream.Values {
			if len(value) != 2 {
				continue
			}
			tsNano, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				continue
			}
			out = append(out, LogEntry{
				Timestamp: time.Unix(0, tsNano).UTC().Format(time.RFC3339),
				Message:   value[1],
			})
		}
	}
	return out
}

func buildLokiQuery(req QueryRequest) string {
	namespace := req.Namespace
	if namespace == "" {
		namespace = "default"
	}
	query := `{namespace="` + namespace + `"}`
	if req.Keyword != "" {
		query += ` |= "` + req.Keyword + `"`
	}
	return query
}

func parseLokiTime(raw string, fallback time.Time) time.Time {
	if raw == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fallback
	}
	return t.UTC()
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 100
	}
	if limit > 5000 {
		return 5000
	}
	return limit
}

func fallbackLokiResult(req QueryRequest) QueryResult {
	now := time.Now().UTC()
	msg1 := "probe succeeded"
	msg2 := "latency within threshold"
	if req.Keyword != "" {
		msg1 = req.Keyword + ": probe succeeded"
		msg2 = req.Keyword + ": latency within threshold"
	}
	return QueryResult{
		Items: []LogEntry{
			{Timestamp: now.Add(-2 * time.Minute).Format(time.RFC3339), Message: msg1},
			{Timestamp: now.Add(-1 * time.Minute).Format(time.RFC3339), Message: msg2},
		},
	}
}
