package marketplace

import (
	"context"
	"encoding/json"
)

func (s *Service) ListTemplateDependencies(ctx context.Context, templateID uint64) (map[uint64][]string, error) {
	versions, err := s.versions.ListByTemplateID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	out := make(map[uint64][]string, len(versions))
	for _, version := range versions {
		var dependencies []string
		if version.DependencySnapshot != "" {
			_ = json.Unmarshal([]byte(version.DependencySnapshot), &dependencies)
		}
		out[version.ID] = dependencies
	}
	return out, nil
}
