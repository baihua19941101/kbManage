package marketplace

func (s *Service) HasUnsupportedPermissionDeclaration(declaration string) bool {
	return containsUnsupportedPermission(declaration)
}
