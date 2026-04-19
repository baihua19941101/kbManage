package marketplace

import (
	"context"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	marketplaceint "kbmanage/backend/internal/integration/marketplace"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListExtensions(ctx context.Context, userID uint64, filter ExtensionListFilter) ([]domain.ExtensionPackage, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	return s.extensions.List(ctx, repository.ExtensionPackageListFilter(filter))
}

func (s *Service) RegisterExtension(ctx context.Context, userID uint64, in CreateExtensionPackageInput) (*domain.ExtensionPackage, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:manage-extension"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.ExtensionType) == "" || strings.TrimSpace(in.Version) == "" {
		return nil, ErrMarketplaceInvalid
	}
	normalizedPermission, err := s.registry.NormalizePermissionDeclaration(ctx, mustJSON(in.PermissionDeclaration))
	if err != nil {
		return nil, ErrMarketplaceInvalid
	}
	normalizedPolicy, err := s.registry.NormalizeCompatibilityPolicy(ctx, mustJSON(in.Compatibility))
	if err != nil {
		return nil, ErrMarketplaceInvalid
	}
	item := &domain.ExtensionPackage{
		Name:                  strings.TrimSpace(in.Name),
		ExtensionType:         strings.TrimSpace(in.ExtensionType),
		Version:               strings.TrimSpace(in.Version),
		Status:                normalizeExtensionStatus("registered"),
		CompatibilityPolicy:   normalizedPolicy,
		PermissionDeclaration: normalizedPermission,
		VisibilityScope:       firstNonEmptyString(in.VisibilityScope, "platform"),
		EntrySummary:          in.EntrySummary,
		OwnerUserID:           userID,
	}
	if err := s.extensions.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := s.replaceExtensionCompatibility(ctx, item, in.Compatibility); err != nil {
		return nil, err
	}
	if err := s.lifecycle.Create(ctx, &domain.ExtensionLifecycleRecord{
		ExtensionPackageID: item.ID,
		Action:             domain.ExtensionLifecycleActionRegister,
		ScopeType:          "platform",
		ScopeRef:           "platform",
		Outcome:            "succeeded",
		ExecutedBy:         userID,
		ExecutedAt:         time.Now(),
	}); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionExtensionRegister, "extension-package", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"name": item.Name, "version": item.Version})
	return item, nil
}

func (s *Service) GetExtensionCompatibility(ctx context.Context, userID, extensionID uint64) (*ExtensionCompatibilityView, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	item, err := s.extensions.GetByID(ctx, extensionID)
	if err != nil {
		return nil, err
	}
	statements, err := s.compatibility.ListByOwner(ctx, domain.CompatibilityOwnerExtension, marketplaceint.CompatibilityOwnerRefForExtension(item))
	if err != nil {
		return nil, err
	}
	return &ExtensionCompatibilityView{Extension: item, Statements: statements, BlockedReasons: isBlockedCompatibility(statements)}, nil
}

func (s *Service) EnableExtension(ctx context.Context, userID, extensionID uint64, in ExtensionLifecycleInput) (*domain.ExtensionPackage, error) {
	if strings.TrimSpace(in.ScopeType) != "" && in.ScopeID != 0 {
		if err := s.scope.EnsureScopePermission(ctx, userID, in.ScopeType, in.ScopeID, "marketplace:manage-extension"); err != nil {
			return nil, err
		}
	} else if err := s.scope.EnsurePermission(ctx, userID, "marketplace:manage-extension"); err != nil {
		return nil, err
	}
	item, err := s.extensions.GetByID(ctx, extensionID)
	if err != nil {
		return nil, err
	}
	view, err := s.GetExtensionCompatibility(ctx, userID, extensionID)
	if err != nil {
		return nil, err
	}
	if containsUnsupportedPermission(item.PermissionDeclaration) || len(view.BlockedReasons) > 0 {
		s.writeAudit(ctx, userID, ActionExtensionEnable, "extension-package", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeFailed, map[string]any{"reason": "compatibility blocked"})
		return nil, ErrMarketplaceBlocked
	}
	item.Status = domain.ExtensionPackageStatusEnabled
	if err := s.extensions.Update(ctx, item); err != nil {
		return nil, err
	}
	if err := s.lifecycle.Create(ctx, &domain.ExtensionLifecycleRecord{
		ExtensionPackageID: item.ID,
		Action:             domain.ExtensionLifecycleActionEnable,
		ScopeType:          firstNonEmptyString(in.ScopeType, "platform"),
		ScopeRef: func() string {
			if in.ScopeID == 0 {
				return "platform"
			}
			return scopeRef(in.ScopeType, in.ScopeID)
		}(),
		Outcome:    "succeeded",
		Reason:     in.Reason,
		ExecutedBy: userID,
		ExecutedAt: time.Now(),
	}); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionExtensionEnable, "extension-package", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"scopeType": in.ScopeType, "scopeId": in.ScopeID})
	return item, nil
}

func (s *Service) DisableExtension(ctx context.Context, userID, extensionID uint64, in ExtensionLifecycleInput) (*domain.ExtensionPackage, error) {
	if strings.TrimSpace(in.ScopeType) != "" && in.ScopeID != 0 {
		if err := s.scope.EnsureScopePermission(ctx, userID, in.ScopeType, in.ScopeID, "marketplace:manage-extension"); err != nil {
			return nil, err
		}
	} else if err := s.scope.EnsurePermission(ctx, userID, "marketplace:manage-extension"); err != nil {
		return nil, err
	}
	item, err := s.extensions.GetByID(ctx, extensionID)
	if err != nil {
		return nil, err
	}
	item.Status = domain.ExtensionPackageStatusDisabled
	if err := s.extensions.Update(ctx, item); err != nil {
		return nil, err
	}
	if err := s.lifecycle.Create(ctx, &domain.ExtensionLifecycleRecord{
		ExtensionPackageID: item.ID,
		Action:             domain.ExtensionLifecycleActionDisable,
		ScopeType:          firstNonEmptyString(in.ScopeType, "platform"),
		ScopeRef: func() string {
			if in.ScopeID == 0 {
				return "platform"
			}
			return scopeRef(in.ScopeType, in.ScopeID)
		}(),
		Outcome:    "succeeded",
		Reason:     in.Reason,
		ExecutedBy: userID,
		ExecutedAt: time.Now(),
	}); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionExtensionDisable, "extension-package", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"scopeType": in.ScopeType, "scopeId": in.ScopeID})
	return item, nil
}

func (s *Service) replaceExtensionCompatibility(ctx context.Context, item *domain.ExtensionPackage, seeds []marketplaceint.CompatibilitySeed) error {
	ownerRef := marketplaceint.CompatibilityOwnerRefForExtension(item)
	statements := make([]domain.CompatibilityStatement, 0, len(seeds))
	now := time.Now()
	for _, seed := range seeds {
		result := domain.CompatibilityResult(seed.Result)
		if result == "" {
			result = domain.CompatibilityResultCompatible
		}
		statements = append(statements, domain.CompatibilityStatement{
			OwnerType:   domain.CompatibilityOwnerExtension,
			OwnerRef:    ownerRef,
			TargetType:  firstNonEmptyString(seed.TargetType, "platform-version"),
			TargetRef:   firstNonEmptyString(seed.TargetRef, "current"),
			Result:      result,
			Summary:     firstNonEmptyString(seed.Summary, "扩展兼容性已校验"),
			EvaluatedAt: now,
			Evaluator:   "marketplace-extension-registry",
		})
	}
	return s.compatibility.ReplaceForOwner(ctx, domain.CompatibilityOwnerExtension, ownerRef, statements)
}

func containsUnsupportedPermission(declaration string) bool {
	normalized := strings.ToLower(strings.TrimSpace(declaration))
	return strings.Contains(normalized, "\"*\"") || strings.Contains(normalized, "system:") || strings.Contains(normalized, "platform-admin")
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
