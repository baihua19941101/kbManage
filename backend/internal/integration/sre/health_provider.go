package sre

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
)

type HealthOverviewInput struct {
	Policy            *domain.HAPolicy
	ActiveWindow      *domain.MaintenanceWindow
	LatestBaseline    *domain.CapacityBaseline
	LatestEvidence    *domain.ScaleEvidence
	LatestUpgradePlan *domain.SREUpgradePlan
	WorkspaceID       uint64
	ProjectID         uint64
}

type HealthProvider interface {
	BuildOverview(ctx context.Context, input HealthOverviewInput) *domain.PlatformHealthSnapshot
}

type StaticHealthProvider struct{}

func NewStaticHealthProvider() HealthProvider { return StaticHealthProvider{} }

func (StaticHealthProvider) BuildOverview(_ context.Context, input HealthOverviewInput) *domain.PlatformHealthSnapshot {
	overall := domain.PlatformHealthOverallHealthy
	maintenanceStatus := "inactive"
	capacityRisk := "low"
	throttling := "normal"
	recovery := "最近未发生异常恢复动作"
	recommended := []string{"继续观察控制面健康与依赖延迟"}

	if input.ActiveWindow != nil {
		maintenanceStatus = string(input.ActiveWindow.Status)
		overall = domain.PlatformHealthOverallMaintenance
		recommended = append(recommended, "维护窗口期间限制高风险变更")
	}
	if input.LatestBaseline != nil {
		capacityRisk = strings.TrimSpace(string(input.LatestBaseline.Status))
		if capacityRisk == "" {
			capacityRisk = "warning"
		}
		if input.LatestBaseline.Status == domain.CapacityBaselineStatusCritical {
			overall = domain.PlatformHealthOverallCritical
			throttling = "engaged"
			recommended = append(recommended, "当前容量风险较高，建议暂停额外升级")
		}
	}
	if input.Policy != nil {
		recovery = strings.TrimSpace(input.Policy.LastRecoveryResult)
		if recovery == "" {
			recovery = fmt.Sprintf("高可用模式 %s 正常运行", input.Policy.DeploymentMode)
		}
		if input.Policy.Status == domain.HAPolicyStatusDegraded {
			overall = domain.PlatformHealthOverallWarning
			recommended = append(recommended, "检查故障切换门槛与依赖连通性")
		}
	}
	if input.LatestUpgradePlan != nil && input.LatestUpgradePlan.Status == domain.SREUpgradeStatusRolling {
		recommended = append(recommended, "升级进行中，关注滚动阶段健康变化")
	}
	if input.LatestEvidence != nil && strings.EqualFold(input.LatestEvidence.ConfidenceLevel, "low") {
		recommended = append(recommended, "规模化结论可信度不足，需要补充样本")
	}

	return &domain.PlatformHealthSnapshot{
		WorkspaceID:             input.WorkspaceID,
		SnapshotAt:              time.Now(),
		ComponentHealthSummary:  `{"controlPlane":"healthy","scheduler":"healthy"}`,
		DependencyHealthSummary: `{"mysql":"healthy","redis":"healthy"}`,
		TaskBacklogSummary:      `{"pending":2,"running":1,"blocked":0}`,
		CapacityRiskLevel:       capacityRisk,
		ThrottlingStatus:        throttling,
		RecoverySummary:         recovery,
		MaintenanceStatus:       maintenanceStatus,
		OverallStatus:           overall,
		RecommendedActions:      strings.Join(recommended, "；"),
	}
}
