import { hasAnyRole, useAuthStore } from '@/features/auth/store';

const BACKUP_RESTORE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'backuprestore:read'
];
const BACKUP_RESTORE_MANAGE_POLICY_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'backuprestore:manage-policy'
];
const BACKUP_RESTORE_BACKUP_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'backuprestore:backup'
];
const BACKUP_RESTORE_RESTORE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'backuprestore:restore'
];
const BACKUP_RESTORE_MIGRATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'backuprestore:migrate'
];
const BACKUP_RESTORE_DRILL_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'backuprestore:drill'
];

export const useBackupRestorePermissions = () => {
  const user = useAuthStore((state) => state.user);
  return {
    canRead: hasAnyRole(user, BACKUP_RESTORE_READ_IDENTIFIERS),
    canManagePolicy: hasAnyRole(user, BACKUP_RESTORE_MANAGE_POLICY_IDENTIFIERS),
    canRunBackup: hasAnyRole(user, BACKUP_RESTORE_BACKUP_IDENTIFIERS),
    canRestore: hasAnyRole(user, BACKUP_RESTORE_RESTORE_IDENTIFIERS),
    canMigrate: hasAnyRole(user, BACKUP_RESTORE_MIGRATE_IDENTIFIERS),
    canDrill: hasAnyRole(user, BACKUP_RESTORE_DRILL_IDENTIFIERS)
  };
};
