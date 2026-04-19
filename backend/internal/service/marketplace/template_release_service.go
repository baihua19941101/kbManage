package marketplace

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

func (s *Service) ListTemplateReleases(ctx context.Context, userID, templateID uint64) ([]domain.TemplateReleaseScope, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	items, err := s.releases.List(ctx, repository.TemplateReleaseScopeListFilter{TemplateID: templateID})
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.TemplateReleaseScope, 0, len(items))
	for _, item := range items {
		scopeType, scopeID, parseErr := parseScopeRef(item.ScopeRef)
		if parseErr != nil {
			continue
		}
		if s.scope.CanReadScope(ctx, userID, scopeType, scopeID) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *Service) CreateTemplateRelease(ctx context.Context, userID, templateID uint64, in CreateTemplateReleaseInput) (*domain.TemplateReleaseScope, error) {
	scopeID := in.ScopeID
	if scopeID == 0 {
		if parsed, err := strconv.ParseUint(strings.TrimSpace(in.TargetRef), 10, 64); err == nil {
			scopeID = parsed
		}
	}
	if err := s.scope.EnsureScopePermission(ctx, userID, in.ScopeType, scopeID, "marketplace:publish-template"); err != nil {
		return nil, err
	}
	template, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	var version *domain.TemplateVersion
	if in.VersionID != 0 {
		version, err = s.versions.GetByID(ctx, in.VersionID)
		if err != nil {
			return nil, err
		}
	} else if strings.TrimSpace(in.Version) != "" {
		version, err = s.versions.FindByTemplateAndVersion(ctx, templateID, in.Version)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrMarketplaceInvalid
	}
	if version.TemplateID != template.ID || version.Status != domain.TemplateVersionStatusActive {
		return nil, ErrMarketplaceBlocked
	}
	compatibility, err := s.compatibility.ListByOwner(ctx, domain.CompatibilityOwnerTemplateVersion, strconv.FormatUint(version.ID, 10))
	if err != nil {
		return nil, err
	}
	if len(isBlockedCompatibility(compatibility)) > 0 {
		return nil, ErrMarketplaceBlocked
	}
	reference := scopeRef(in.ScopeType, scopeID)
	item, err := s.releases.FindByTemplateScope(ctx, templateID, in.ScopeType, reference)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	if item == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		item = &domain.TemplateReleaseScope{
			TemplateID:     templateID,
			VersionID:      version.ID,
			ScopeType:      strings.TrimSpace(in.ScopeType),
			ScopeRef:       reference,
			Status:         domain.TemplateReleaseStatusPublished,
			VisibilityMode: firstNonEmptyString(in.VisibilityMode, "scope"),
			PublishedBy:    userID,
			PublishedAt:    now,
		}
		if err := s.releases.Create(ctx, item); err != nil {
			return nil, err
		}
	} else {
		item.VersionID = version.ID
		item.Status = domain.TemplateReleaseStatusPublished
		item.VisibilityMode = firstNonEmptyString(in.VisibilityMode, item.VisibilityMode)
		item.PublishedBy = userID
		item.PublishedAt = now
		item.WithdrawnAt = nil
		if err := s.releases.Update(ctx, item); err != nil {
			return nil, err
		}
	}

	if err := s.upsertInstallationRecord(ctx, template, version, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionTemplateRelease, "template-release", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{
		"templateId": templateID,
		"versionId":  version.ID,
		"scopeType":  in.ScopeType,
		"scopeId":    scopeID,
	})
	return item, nil
}

func (s *Service) ListInstallationRecords(ctx context.Context, userID uint64, filter InstallationListFilter) ([]domain.InstallationRecord, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	items, err := s.installations.List(ctx, repository.InstallationRecordListFilter{
		ScopeType: strings.TrimSpace(filter.ScopeType),
		ScopeRef: func() string {
			if filter.ScopeID == 0 {
				return ""
			}
			return scopeRef(filter.ScopeType, filter.ScopeID)
		}(),
		Status: strings.TrimSpace(filter.Status),
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.InstallationRecord, 0, len(items))
	for _, item := range items {
		scopeType, scopeID, parseErr := parseScopeRef(item.ScopeRef)
		if parseErr != nil {
			continue
		}
		if s.scope.CanReadScope(ctx, userID, scopeType, scopeID) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *Service) upsertInstallationRecord(ctx context.Context, template *domain.ApplicationTemplate, version *domain.TemplateVersion, release *domain.TemplateReleaseScope) error {
	item, err := s.installations.FindByTemplateScope(ctx, template.ID, release.ScopeType, release.ScopeRef)
	now := time.Now()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	changeSummary := firstNonEmptyString(version.ReleaseNotes, template.ReleaseNotesSummary, "模板已发布到目标范围")
	if item == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return s.installations.Create(ctx, &domain.InstallationRecord{
			TemplateID:              template.ID,
			VersionID:               version.ID,
			ScopeType:               release.ScopeType,
			ScopeRef:                release.ScopeRef,
			ReleaseScopeID:          release.ID,
			LifecycleStatus:         domain.InstallationLifecycleInstalled,
			CurrentInstalledVersion: version.Version,
			UpgradeTargetVersion:    version.Version,
			ChangeSummary:           changeSummary,
			InstalledAt:             now,
			LastChangedAt:           now,
		})
	}
	item.VersionID = version.ID
	item.ReleaseScopeID = release.ID
	item.LifecycleStatus = domain.InstallationLifecycleUpgraded
	item.UpgradeTargetVersion = version.Version
	item.CurrentInstalledVersion = version.Version
	item.ChangeSummary = changeSummary
	item.LastChangedAt = now
	return s.installations.Update(ctx, item)
}
