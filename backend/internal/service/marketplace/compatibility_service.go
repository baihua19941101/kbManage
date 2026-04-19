package marketplace

import "kbmanage/backend/internal/domain"

func (s *Service) blockedCompatibilityReasons(items []domain.CompatibilityStatement) []string {
	return isBlockedCompatibility(items)
}
