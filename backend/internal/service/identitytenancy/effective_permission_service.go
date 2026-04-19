package identitytenancy

import (
	"context"
	"strings"
	"time"
)

func (s *Service) EffectivePermissionSummary(ctx context.Context, subjectType, subjectRef string) (map[string]any, error) {
	assignments, err := s.assignments.ListActiveBySubject(ctx, subjectType, subjectRef, time.Now())
	if err != nil {
		return nil, err
	}
	levels := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		levels = append(levels, strings.TrimSpace(assignment.ScopeType))
	}
	return map[string]any{
		"subjectType":       subjectType,
		"subjectRef":        subjectRef,
		"activeAssignments": len(assignments),
		"scopeLevels":       levels,
	}, nil
}
