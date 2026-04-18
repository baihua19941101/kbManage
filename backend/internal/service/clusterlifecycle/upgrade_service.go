package clusterlifecycle

import "context"

func (s *Service) ScheduleUpgrade(ctx context.Context, userID, clusterID uint64, input CreateUpgradePlanInput) (*UpgradeExecutionBundle, error) {
	plan, err := s.CreateUpgradePlan(ctx, userID, clusterID, input)
	if err != nil {
		return nil, err
	}
	return &UpgradeExecutionBundle{Plan: plan}, nil
}

type UpgradeExecutionBundle struct {
	Plan any `json:"plan"`
}
