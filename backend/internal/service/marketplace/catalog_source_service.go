package marketplace

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	marketplaceint "kbmanage/backend/internal/integration/marketplace"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

func (s *Service) ListCatalogSources(ctx context.Context, userID uint64, filter CatalogSourceListFilter) ([]domain.CatalogSource, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	return s.sources.List(ctx, repository.CatalogSourceListFilter(filter))
}

func (s *Service) CreateCatalogSource(ctx context.Context, userID uint64, in CreateCatalogSourceInput) (*domain.CatalogSource, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:manage-source"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.SourceType) == "" || strings.TrimSpace(in.EndpointRef) == "" {
		return nil, ErrMarketplaceInvalid
	}
	envelope := marketplaceint.CatalogSeedEnvelope{Templates: in.TemplateSeeds}
	item := &domain.CatalogSource{
		Name:            strings.TrimSpace(in.Name),
		SourceType:      domain.CatalogSourceType(strings.TrimSpace(in.SourceType)),
		EndpointRef:     strings.TrimSpace(in.EndpointRef),
		Status:          normalizeSourceStatus(in.Status),
		SyncState:       domain.CatalogSyncStateIdle,
		OwnerUserID:     userID,
		VisibilityScope: firstNonEmptyString(in.VisibilityScope, "platform"),
		ConfigSummary:   mustJSON(envelope),
	}
	if err := s.sources.Create(ctx, item); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrMarketplaceConflict
		}
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionCatalogSourceCreate, "catalog-source", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{
		"name":       item.Name,
		"sourceType": item.SourceType,
	})
	return item, nil
}

func (s *Service) SyncCatalogSource(ctx context.Context, userID, sourceID uint64) (*domain.CatalogSource, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:manage-source"); err != nil {
		return nil, err
	}
	source, err := s.sources.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	source.SyncState = domain.CatalogSyncStateSyncing
	_ = s.sources.Update(ctx, source)

	result, err := s.catalog.Sync(ctx, *source)
	if err != nil {
		source.SyncState = domain.CatalogSyncStateFailed
		source.Status = domain.CatalogSourceStatusDegraded
		source.LastError = err.Error()
		_ = s.sources.Update(ctx, source)
		s.writeAudit(ctx, userID, ActionCatalogSourceSync, "catalog-source", strconv.FormatUint(source.ID, 10), domain.AuditOutcomeFailed, map[string]any{"error": err.Error()})
		return nil, err
	}

	for _, seed := range result.Templates {
		template, err := s.upsertTemplateFromSeed(ctx, source, seed)
		if err != nil {
			return nil, err
		}
		if _, err := s.upsertTemplateVersionsFromSeed(ctx, template, seed); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	source.SyncState = domain.CatalogSyncStateSucceeded
	source.Status = domain.CatalogSourceStatusActive
	source.LastSyncedAt = &now
	source.LastError = ""
	if err := s.sources.Update(ctx, source); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionCatalogSourceSync, "catalog-source", strconv.FormatUint(source.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"templateCount": len(result.Templates)})
	return source, nil
}

func (s *Service) upsertTemplateFromSeed(ctx context.Context, source *domain.CatalogSource, seed marketplaceint.TemplateSeed) (*domain.ApplicationTemplate, error) {
	item, err := s.templates.FindBySourceAndSlug(ctx, source.ID, seed.Slug)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		item = &domain.ApplicationTemplate{
			CatalogSourceID: source.ID,
			Name:            strings.TrimSpace(seed.Name),
			Slug:            strings.TrimSpace(seed.Slug),
			Category:        firstNonEmptyString(seed.Category, "general"),
		}
	}
	item.Name = strings.TrimSpace(seed.Name)
	item.Category = firstNonEmptyString(seed.Category, "general")
	item.Summary = seed.Summary
	item.PublishStatus = normalizeTemplatePublishStatus(seed.PublishStatus)
	item.SupportedScopes = mustJSON(seed.SupportedScopes)
	item.ReleaseNotesSummary = seed.ReleaseNotesSummary
	if item.ID == 0 {
		return item, s.templates.Create(ctx, item)
	}
	return item, s.templates.Update(ctx, item)
}

func (s *Service) upsertTemplateVersionsFromSeed(ctx context.Context, template *domain.ApplicationTemplate, seed marketplaceint.TemplateSeed) ([]domain.TemplateVersion, error) {
	versionItems, err := s.versions.ListByTemplateID(ctx, template.ID)
	if err != nil {
		return nil, err
	}
	for _, versionSeed := range seed.Versions {
		item, err := s.versions.FindByTemplateAndVersion(ctx, template.ID, versionSeed.Version)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			item = &domain.TemplateVersion{TemplateID: template.ID, Version: strings.TrimSpace(versionSeed.Version)}
		}
		item.Status = normalizeTemplateVersionStatus(versionSeed.Status)
		item.DependencySnapshot = mustJSON(versionSeed.Dependencies)
		item.ParameterSchemaSummary = versionSeed.ParameterSchemaSummary
		item.DeploymentConstraintSummary = versionSeed.DeploymentConstraintSummary
		item.ReleaseNotes = versionSeed.ReleaseNotes
		if versionSeed.IsUpgradeable != nil {
			item.IsUpgradeable = *versionSeed.IsUpgradeable
		} else {
			item.IsUpgradeable = true
		}
		item.SupersedesVersionID = supersededVersionID(versionItems, versionSeed.SupersedesVersion)
		if item.ID == 0 {
			if err := s.versions.Create(ctx, item); err != nil {
				return nil, err
			}
			versionItems = append(versionItems, *item)
		} else if err := s.versions.Update(ctx, item); err != nil {
			return nil, err
		}
		if template.DefaultVersionID == nil && item.Status == domain.TemplateVersionStatusActive {
			template.DefaultVersionID = &item.ID
			if err := s.templates.Update(ctx, template); err != nil {
				return nil, err
			}
		}
		if item.Status == domain.TemplateVersionStatusActive {
			ownerRef := strconv.FormatUint(item.ID, 10)
			compatibility := []domain.CompatibilityStatement{{
				OwnerType:   domain.CompatibilityOwnerTemplateVersion,
				OwnerRef:    ownerRef,
				TargetType:  "scope",
				TargetRef:   "default",
				Result:      domain.CompatibilityResultCompatible,
				Summary:     "模板版本可分发",
				EvaluatedAt: time.Now(),
				Evaluator:   "marketplace-sync",
			}}
			if err := s.compatibility.ReplaceForOwner(ctx, domain.CompatibilityOwnerTemplateVersion, ownerRef, compatibility); err != nil {
				return nil, err
			}
		}
	}
	return s.versions.ListByTemplateID(ctx, template.ID)
}
