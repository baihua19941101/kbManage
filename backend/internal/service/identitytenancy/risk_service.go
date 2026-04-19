package identitytenancy

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListAccessRisks(ctx context.Context, userID uint64, filter AccessRiskListFilter) ([]domain.AccessRiskSnapshot, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.rebuildDerivedRisks(ctx); err != nil {
		return nil, err
	}
	items, err := s.risks.List(ctx, repository.AccessRiskListFilter{
		SubjectType: filter.SubjectType,
		Severity:    filter.Severity,
		Status:      string(domain.AccessRiskStatusOpen),
	})
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionAccessRiskQuery, ResourceTypeAccessRisk, "list", domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"count": len(items),
	})
	return items, nil
}

func (s *Service) recordRiskForAssignment(ctx context.Context, role *domain.RoleDefinition, assignment *domain.RoleAssignment) error {
	level := riskLevelForAssignment(role, assignment)
	if level == domain.IdentityRiskLevelLow {
		return nil
	}
	return s.risks.Create(ctx, &domain.AccessRiskSnapshot{
		SubjectType:       assignment.SubjectType,
		SubjectRef:        assignment.SubjectRef,
		RiskType:          "assignment-expansion",
		Severity:          level,
		Summary:           fmt.Sprintf("%s 在 %s 范围存在 %s 授权风险", assignment.SubjectRef, assignment.ScopeType, assignment.SourceType),
		RecommendedAction: "复核继承边界并确认是否需要回收",
		Status:            domain.AccessRiskStatusOpen,
		GeneratedAt:       time.Now(),
	})
}

func (s *Service) rebuildDerivedRisks(ctx context.Context) error {
	assignments, err := s.assignments.List(ctx, repository.RoleAssignmentListFilter{})
	if err != nil {
		return err
	}
	for _, assignment := range assignments {
		if assignment.Status == domain.RoleAssignmentStatusExpired {
			_ = s.risks.Create(ctx, &domain.AccessRiskSnapshot{
				SubjectType:       assignment.SubjectType,
				SubjectRef:        assignment.SubjectRef,
				RiskType:          "residual-access",
				Severity:          domain.IdentityRiskLevelHigh,
				Summary:           "已过期授权需要确认会话残留访问",
				RecommendedAction: "刷新会话并执行访问回收",
				Status:            domain.AccessRiskStatusOpen,
				GeneratedAt:       time.Now(),
			})
		}
		if assignment.SourceType == domain.RoleAssignmentSourceDelegated || assignment.SourceType == domain.RoleAssignmentSourceTemporary {
			_ = s.risks.Create(ctx, &domain.AccessRiskSnapshot{
				SubjectType:       assignment.SubjectType,
				SubjectRef:        assignment.SubjectRef,
				RiskType:          "delegated-or-temporary",
				Severity:          domain.IdentityRiskLevelHigh,
				Summary:           "委派或临时授权需要持续关注到期与回收状态",
				RecommendedAction: "确认委派链路与到期策略",
				Status:            domain.AccessRiskStatusOpen,
				GeneratedAt:       time.Now(),
			})
		}
	}
	sessions, err := s.sessions.List(ctx, repository.SessionRecordListFilter{})
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if session.Status == domain.IdentitySessionStatusRiskBlocked || session.RiskLevel == domain.IdentityRiskLevelCritical {
			_ = s.risks.Create(ctx, &domain.AccessRiskSnapshot{
				SubjectType:       "user",
				SubjectRef:        strconv.FormatUint(session.UserID, 10),
				RiskType:          "session-risk",
				Severity:          session.RiskLevel,
				Summary:           "存在高风险会话，需要审查权限变更影响",
				RecommendedAction: "终止会话并重新登录",
				Status:            domain.AccessRiskStatusOpen,
				GeneratedAt:       time.Now(),
			})
		}
	}
	return nil
}
