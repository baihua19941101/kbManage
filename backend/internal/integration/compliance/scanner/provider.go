package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"kbmanage/backend/internal/domain"
)

type Request struct {
	ExecutionID uint64
	Profile     domain.ScanProfile
	Baseline    domain.ComplianceBaselineSnapshot
}

type Evidence struct {
	EvidenceType    domain.ComplianceEvidenceType
	SourceRef       string
	Confidence      domain.ComplianceEvidenceConfidence
	Summary         string
	ArtifactRef     string
	RedactionStatus domain.ComplianceEvidenceRedactionStatus
	Payload         map[string]any
}

type Finding struct {
	ControlID    string
	ControlTitle string
	Result       domain.ComplianceFindingResult
	RiskLevel    domain.ComplianceRiskLevel
	ClusterID    *uint64
	NodeName     string
	Namespace    string
	ResourceKind string
	ResourceName string
	ResourceUID  string
	Summary      string
	Evidences    []Evidence
}

type PartialFailure struct {
	ScopeRef string
	Reason   string
}

type Result struct {
	Status          domain.ComplianceScanStatus
	CoverageStatus  domain.ComplianceCoverageStatus
	StartedAt       time.Time
	CompletedAt     time.Time
	Score           float64
	PassCount       int
	FailCount       int
	WarningCount    int
	ErrorSummary    string
	Findings        []Finding
	PartialFailures []PartialFailure
}

type Provider interface {
	Execute(ctx context.Context, req Request) (Result, error)
}

type MockProvider struct{}

func NewMockProvider() Provider {
	return MockProvider{}
}

func (MockProvider) Execute(_ context.Context, req Request) (Result, error) {
	startedAt := time.Now().UTC()
	completedAt := startedAt.Add(2 * time.Second)
	findings := make([]Finding, 0, 4)
	partials := make([]PartialFailure, 0)
	clusterIDs := extractClusterRefs(req.Profile)
	if len(clusterIDs) == 0 {
		clusterIDs = []uint64{1}
	}
	for idx, clusterID := range clusterIDs {
		cid := clusterID
		findings = append(findings,
			Finding{
				ControlID:    fmt.Sprintf("%s-%02d", req.Baseline.Version, idx+1),
				ControlTitle: "API Server audit logging enabled",
				Result:       domain.ComplianceFindingResultFail,
				RiskLevel:    domain.ComplianceRiskLevelHigh,
				ClusterID:    &cid,
				NodeName:     firstNode(req.Profile),
				Namespace:    firstNamespace(req.Profile),
				ResourceKind: firstResourceKind(req.Profile),
				ResourceName: fmt.Sprintf("target-%d", clusterID),
				ResourceUID:  fmt.Sprintf("cluster-%d-control-audit", clusterID),
				Summary:      "audit-log-path or policy is not aligned with selected baseline",
				Evidences: []Evidence{
					{
						EvidenceType:    domain.ComplianceEvidenceTypeConfiguration,
						SourceRef:       fmt.Sprintf("cluster/%d/apiserver", clusterID),
						Confidence:      domain.ComplianceEvidenceConfidenceHigh,
						Summary:         "captured audit policy and startup flags",
						ArtifactRef:     fmt.Sprintf("memory://scan/%d/cluster/%d/audit-policy", req.ExecutionID, clusterID),
						RedactionStatus: domain.ComplianceEvidenceRedactionMasked,
						Payload:         map[string]any{"flag": "--audit-policy-file", "present": false},
					},
				},
			},
			Finding{
				ControlID:    fmt.Sprintf("%s-%02d-pass", req.Baseline.Version, idx+1),
				ControlTitle: "Anonymous auth disabled",
				Result:       domain.ComplianceFindingResultPass,
				RiskLevel:    domain.ComplianceRiskLevelLow,
				ClusterID:    &cid,
				NodeName:     firstNode(req.Profile),
				Namespace:    firstNamespace(req.Profile),
				ResourceKind: firstResourceKind(req.Profile),
				ResourceName: fmt.Sprintf("target-%d", clusterID),
				ResourceUID:  fmt.Sprintf("cluster-%d-control-anon", clusterID),
				Summary:      "anonymous-auth is disabled",
			},
		)
	}
	status := domain.ComplianceScanStatusSucceeded
	coverage := domain.ComplianceCoverageStatusFull
	errorSummary := ""
	if req.Profile.ScopeType == domain.ComplianceScopeTypeNode && firstNode(req.Profile) == "" {
		status = domain.ComplianceScanStatusPartiallySucceeded
		coverage = domain.ComplianceCoverageStatusPartial
		partials = append(partials, PartialFailure{ScopeRef: "node-selector", Reason: "node selector resolved empty target set"})
		errorSummary = "some requested node targets were not resolved"
	}
	return Result{
		Status:          status,
		CoverageStatus:  coverage,
		StartedAt:       startedAt,
		CompletedAt:     completedAt,
		Score:           66.7,
		PassCount:       len(clusterIDs),
		FailCount:       len(clusterIDs),
		WarningCount:    0,
		ErrorSummary:    errorSummary,
		Findings:        findings,
		PartialFailures: partials,
	}, nil
}

func extractClusterRefs(profile domain.ScanProfile) []uint64 {
	var refs []uint64
	_ = jsonUnmarshalStringArrayAsUint(profile.ClusterRefsJSON, &refs)
	return refs
}

func firstNode(profile domain.ScanProfile) string {
	var items []map[string]string
	_ = jsonUnmarshal(profile.NodeSelectorsJSON, &items)
	if len(items) == 0 {
		return ""
	}
	if value := items[0]["nodeName"]; value != "" {
		return value
	}
	if value := items[0]["hostname"]; value != "" {
		return value
	}
	return ""
}

func firstNamespace(profile domain.ScanProfile) string {
	var refs []string
	_ = jsonUnmarshal(profile.NamespaceRefsJSON, &refs)
	if len(refs) == 0 {
		return ""
	}
	return refs[0]
}

func firstResourceKind(profile domain.ScanProfile) string {
	var refs []string
	_ = jsonUnmarshal(profile.ResourceKindsJSON, &refs)
	if len(refs) == 0 {
		return ""
	}
	return refs[0]
}

func jsonUnmarshal(raw string, out any) error {
	if raw == "" {
		return nil
	}
	return json.Unmarshal([]byte(raw), out)
}

func jsonUnmarshalStringArrayAsUint(raw string, out *[]uint64) error {
	if raw == "" {
		return nil
	}
	var stringsOut []uint64
	return json.Unmarshal([]byte(raw), &stringsOut)
}
