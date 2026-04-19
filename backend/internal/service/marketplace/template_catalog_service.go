package marketplace

import (
	"context"
	"strconv"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListTemplates(ctx context.Context, userID uint64, filter TemplateListFilter) ([]domain.ApplicationTemplate, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	return s.templates.List(ctx, repository.ApplicationTemplateListFilter(filter))
}

func (s *Service) GetTemplateDetail(ctx context.Context, userID, templateID uint64) (*TemplateDetail, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "marketplace:read"); err != nil {
		return nil, err
	}
	template, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	versions, err := s.versions.ListByTemplateID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	releases, err := s.releases.List(ctx, repository.TemplateReleaseScopeListFilter{TemplateID: templateID})
	if err != nil {
		return nil, err
	}
	var compatibility []domain.CompatibilityStatement
	for _, item := range versions {
		list, err := s.compatibility.ListByOwner(ctx, domain.CompatibilityOwnerTemplateVersion, strconv.FormatUint(item.ID, 10))
		if err != nil {
			return nil, err
		}
		compatibility = append(compatibility, list...)
	}
	return &TemplateDetail{Template: template, Versions: versions, Releases: releases, Compatibility: compatibility}, nil
}
