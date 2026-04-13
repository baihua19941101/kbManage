package workloadops

import "context"

func (s *Service) buildRevisions(_ context.Context, target WorkloadReference) []ReleaseRevision {
	return []ReleaseRevision{
		{
			Revision:          3,
			SourceKind:        "replicaset",
			SourceName:        target.ResourceName + "-rs-3",
			IsCurrent:         true,
			RollbackAvailable: true,
			Summary:           "current revision",
		},
		{
			Revision:          2,
			SourceKind:        "replicaset",
			SourceName:        target.ResourceName + "-rs-2",
			IsCurrent:         false,
			RollbackAvailable: true,
			Summary:           "previous stable revision",
		},
	}
}
