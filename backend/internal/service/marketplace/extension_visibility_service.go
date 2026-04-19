package marketplace

func (s *Service) ResolveExtensionVisibility(scopeType string, scopeID uint64) string {
	if scopeID == 0 {
		return "platform"
	}
	return scopeRef(scopeType, scopeID)
}
