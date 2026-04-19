package marketplace

import "kbmanage/backend/internal/domain"

func (s *Service) BuildInstallationUpgradeState(current *domain.InstallationRecord, targetVersion string) domain.InstallationLifecycleStatus {
	if current == nil {
		return domain.InstallationLifecycleInstalled
	}
	if current.CurrentInstalledVersion != targetVersion {
		return domain.InstallationLifecycleUpgraded
	}
	return current.LifecycleStatus
}
