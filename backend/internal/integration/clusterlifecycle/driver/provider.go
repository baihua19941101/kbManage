package driver

import (
	"context"
	"fmt"
	"strings"
)

type OperationResult struct {
	Summary string `json:"summary"`
	Warning string `json:"warning,omitempty"`
}

type ProvisionRequest struct {
	ClusterName        string
	InfrastructureType string
	DriverKey          string
	DriverVersion      string
	KubernetesVersion  string
	NodePools          []map[string]any
	Parameters         map[string]any
}

type UpgradeRequest struct {
	ClusterName   string
	FromVersion   string
	ToVersion     string
	ImpactSummary string
}

type NodePoolScaleRequest struct {
	ClusterName  string
	NodePoolName string
	DesiredCount int
}

type RetirementRequest struct {
	ClusterName string
	Reason      string
}

type Provider interface {
	ImportCluster(ctx context.Context, apiServer string) (OperationResult, error)
	IssueRegistration(ctx context.Context, clusterName string) (string, error)
	ProvisionCluster(ctx context.Context, req ProvisionRequest) (OperationResult, error)
	UpgradeCluster(ctx context.Context, req UpgradeRequest) (OperationResult, error)
	ScaleNodePool(ctx context.Context, req NodePoolScaleRequest) (OperationResult, error)
	DisableCluster(ctx context.Context, clusterName string) (OperationResult, error)
	RetireCluster(ctx context.Context, req RetirementRequest) (OperationResult, error)
}

type StaticProvider struct{}

func NewStaticProvider() Provider {
	return &StaticProvider{}
}

func (p *StaticProvider) ImportCluster(_ context.Context, apiServer string) (OperationResult, error) {
	if strings.TrimSpace(apiServer) == "" {
		return OperationResult{}, fmt.Errorf("apiServer is required")
	}
	return OperationResult{Summary: "cluster import accepted"}, nil
}

func (p *StaticProvider) IssueRegistration(_ context.Context, clusterName string) (string, error) {
	if strings.TrimSpace(clusterName) == "" {
		return "", fmt.Errorf("clusterName is required")
	}
	return fmt.Sprintf("kubectl apply -f https://kbmanage.local/register/%s.yaml", strings.ToLower(strings.ReplaceAll(clusterName, " ", "-"))), nil
}

func (p *StaticProvider) ProvisionCluster(_ context.Context, req ProvisionRequest) (OperationResult, error) {
	if strings.TrimSpace(req.ClusterName) == "" || strings.TrimSpace(req.DriverKey) == "" {
		return OperationResult{}, fmt.Errorf("clusterName and driverKey are required")
	}
	return OperationResult{Summary: "cluster provisioned via static provider"}, nil
}

func (p *StaticProvider) UpgradeCluster(_ context.Context, req UpgradeRequest) (OperationResult, error) {
	if strings.TrimSpace(req.ToVersion) == "" {
		return OperationResult{}, fmt.Errorf("toVersion is required")
	}
	return OperationResult{Summary: "upgrade completed"}, nil
}

func (p *StaticProvider) ScaleNodePool(_ context.Context, req NodePoolScaleRequest) (OperationResult, error) {
	if strings.TrimSpace(req.NodePoolName) == "" || req.DesiredCount < 0 {
		return OperationResult{}, fmt.Errorf("invalid node pool scale request")
	}
	return OperationResult{Summary: "node pool scaled"}, nil
}

func (p *StaticProvider) DisableCluster(_ context.Context, clusterName string) (OperationResult, error) {
	if strings.TrimSpace(clusterName) == "" {
		return OperationResult{}, fmt.Errorf("clusterName is required")
	}
	return OperationResult{Summary: "cluster disabled"}, nil
}

func (p *StaticProvider) RetireCluster(_ context.Context, req RetirementRequest) (OperationResult, error) {
	if strings.TrimSpace(req.ClusterName) == "" || strings.TrimSpace(req.Reason) == "" {
		return OperationResult{}, fmt.Errorf("clusterName and reason are required")
	}
	return OperationResult{Summary: "cluster retired"}, nil
}
