import { hasAnyRole, useAuthStore } from '@/features/auth/store';

const MARKETPLACE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'marketplace:read'
];
const MARKETPLACE_MANAGE_SOURCE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'marketplace:manage-source'
];
const MARKETPLACE_PUBLISH_TEMPLATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'marketplace:publish-template'
];
const MARKETPLACE_MANAGE_EXTENSION_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'marketplace:manage-extension'
];

export const usePlatformMarketplacePermissions = () => {
  const user = useAuthStore((state) => state.user);
  return {
    canRead: hasAnyRole(user, MARKETPLACE_READ_IDENTIFIERS),
    canManageSource: hasAnyRole(user, MARKETPLACE_MANAGE_SOURCE_IDENTIFIERS),
    canPublishTemplate: hasAnyRole(user, MARKETPLACE_PUBLISH_TEMPLATE_IDENTIFIERS),
    canManageExtension: hasAnyRole(user, MARKETPLACE_MANAGE_EXTENSION_IDENTIFIERS)
  };
};
