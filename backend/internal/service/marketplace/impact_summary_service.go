package marketplace

func (s *Service) BuildImpactSummary(reason string) string {
	return firstNonEmptyString(reason, "需要关注扩展停用对现有能力的影响")
}
