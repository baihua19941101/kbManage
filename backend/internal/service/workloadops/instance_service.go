package workloadops

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) buildInstances(_ context.Context, target WorkloadReference) []WorkloadInstance {
	now := time.Now().Add(-5 * time.Minute)
	return []WorkloadInstance{
		{
			PodName:           fmt.Sprintf("%s-pod-0", target.ResourceName),
			ContainerName:     "app",
			NodeName:          "node-unknown",
			Phase:             "Running",
			Ready:             true,
			RestartCount:      0,
			StartedAt:         &now,
			LastTransitionAt:  &now,
			LogAvailable:      true,
			TerminalAvailable: true,
		},
	}
}
