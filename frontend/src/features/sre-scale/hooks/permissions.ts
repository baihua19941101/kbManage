import { hasAnyRole, useAuthStore } from '@/features/auth/store';

const READ_IDENTIFIERS: readonly string[] = ['platform-admin', 'ops-operator', 'audit-reader', 'readonly', 'sre:read'];
const HA_IDENTIFIERS: readonly string[] = ['platform-admin', 'ops-operator', 'sre:manage-ha'];
const UPGRADE_IDENTIFIERS: readonly string[] = ['platform-admin', 'ops-operator', 'sre:manage-upgrade'];
const SCALE_IDENTIFIERS: readonly string[] = ['platform-admin', 'ops-operator', 'sre:manage-scale'];

export const useSREPermissions = () => {
  const user = useAuthStore((state) => state.user);
  return {
    canRead: hasAnyRole(user, READ_IDENTIFIERS),
    canManageHA: hasAnyRole(user, HA_IDENTIFIERS),
    canManageUpgrade: hasAnyRole(user, UPGRADE_IDENTIFIERS),
    canManageScale: hasAnyRole(user, SCALE_IDENTIFIERS)
  };
};
