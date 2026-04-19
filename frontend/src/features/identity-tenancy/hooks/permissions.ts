import { hasAnyRole, useAuthStore } from '@/features/auth/store';

const IDENTITY_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'identity:read'
];
const IDENTITY_MANAGE_SOURCE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'identity:manage-source'
];
const IDENTITY_MANAGE_ORG_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'identity:manage-org'
];
const IDENTITY_MANAGE_ROLE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'identity:manage-role'
];
const IDENTITY_DELEGATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'identity:delegate'
];
const IDENTITY_SESSION_GOVERN_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:session-govern'
];

export const useIdentityTenancyPermissions = () => {
  const user = useAuthStore((state) => state.user);

  return {
    canRead: hasAnyRole(user, IDENTITY_READ_IDENTIFIERS),
    canManageSource: hasAnyRole(user, IDENTITY_MANAGE_SOURCE_IDENTIFIERS),
    canManageOrg: hasAnyRole(user, IDENTITY_MANAGE_ORG_IDENTIFIERS),
    canManageRole: hasAnyRole(user, IDENTITY_MANAGE_ROLE_IDENTIFIERS),
    canDelegate: hasAnyRole(user, IDENTITY_DELEGATE_IDENTIFIERS),
    canGovernSession: hasAnyRole(user, IDENTITY_SESSION_GOVERN_IDENTIFIERS),
    canReadAudit: hasAnyRole(user, ['platform-admin', 'audit-reader', 'identity:read'])
  };
};
