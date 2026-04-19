package marketplace

func (s *Service) BuildChangeLogSummary(currentVersion, targetVersion, releaseNotes string) string {
	if targetVersion == "" || targetVersion == currentVersion {
		return releaseNotes
	}
	if releaseNotes == "" {
		return "版本从 " + currentVersion + " 升级到 " + targetVersion
	}
	return releaseNotes
}
