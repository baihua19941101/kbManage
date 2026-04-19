package marketplace

import "kbmanage/backend/internal/domain"

func (s *Service) BuildUpgradeAdvice(currentVersion string, versions []domain.TemplateVersion) string {
	for _, version := range versions {
		if version.Version != currentVersion && version.IsUpgradeable {
			return version.Version
		}
	}
	return ""
}
