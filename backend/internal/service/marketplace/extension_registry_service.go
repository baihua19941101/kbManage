package marketplace

func (s *Service) NormalizeExtensionRegistrationScope(scope string) string {
	return firstNonEmptyString(scope, "platform")
}
