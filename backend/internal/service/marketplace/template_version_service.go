package marketplace

import (
	"context"

	"kbmanage/backend/internal/domain"
)

func (s *Service) ListTemplateVersions(ctx context.Context, templateID uint64) ([]domain.TemplateVersion, error) {
	return s.versions.ListByTemplateID(ctx, templateID)
}
