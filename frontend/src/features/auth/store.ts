import { create } from 'zustand';

export type PlatformRole = 'platform-admin' | 'ops-operator' | 'audit-reader' | 'readonly';

export type GitOpsPermission =
  | 'gitops:read'
  | 'gitops:manage-source'
  | 'gitops:sync'
  | 'gitops:promote'
  | 'gitops:rollback'
  | 'gitops:override';

export type PolicyPermission =
  | 'policy:read'
  | 'policy:manage'
  | 'policy:assign'
  | 'policy:approve-exception'
  | 'policy:remediate'
  | 'securitypolicy:read'
  | 'securitypolicy:manage'
  | 'securitypolicy:enforce';

export type CompliancePermission =
  | 'compliance:read'
  | 'compliance:manage-baseline'
  | 'compliance:execute-scan'
  | 'compliance:manage-remediation'
  | 'compliance:review-exception'
  | 'compliance:export-archive';

export type ClusterLifecyclePermission =
  | 'clusterlifecycle:read'
  | 'clusterlifecycle:import'
  | 'clusterlifecycle:create'
  | 'clusterlifecycle:upgrade'
  | 'clusterlifecycle:manage-nodepool'
  | 'clusterlifecycle:retire'
  | 'clusterlifecycle:manage-driver';

export type BackupRestorePermission =
  | 'backuprestore:read'
  | 'backuprestore:manage-policy'
  | 'backuprestore:backup'
  | 'backuprestore:restore'
  | 'backuprestore:migrate'
  | 'backuprestore:drill';

export type IdentityTenancyPermission =
  | 'identity:read'
  | 'identity:manage-source'
  | 'identity:manage-org'
  | 'identity:manage-role'
  | 'identity:delegate'
  | 'identity:session-govern';

export type PlatformMarketplacePermission =
  | 'marketplace:read'
  | 'marketplace:manage-source'
  | 'marketplace:publish-template'
  | 'marketplace:manage-extension';

export type SREScalePermission =
  | 'sre:read'
  | 'sre:manage-ha'
  | 'sre:manage-upgrade'
  | 'sre:manage-scale';

export type EnterprisePolishPermission =
  | 'enterprise:read'
  | 'enterprise:manage-audit'
  | 'enterprise:manage-reports'
  | 'enterprise:manage-delivery';

type AuthUser = {
  id: string;
  username: string;
  displayName?: string;
  platformRoles?: string[];
};

type AuthState = {
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  setSession: (session: {
    accessToken: string;
    refreshToken: string;
    user: AuthUser;
  }) => void;
  clearSession: () => void;
};

const SESSION_STORAGE_KEY = 'kbm-auth-session';

const canUseStorage =
  typeof window !== 'undefined' && typeof window.sessionStorage !== 'undefined';

const restoreSession = (): Pick<
  AuthState,
  'accessToken' | 'refreshToken' | 'user' | 'isAuthenticated'
> => {
  if (!canUseStorage) {
    return {
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false
    };
  }

  const raw = window.sessionStorage.getItem(SESSION_STORAGE_KEY);
  if (!raw) {
    return {
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false
    };
  }

  try {
    const parsed = JSON.parse(raw) as {
      accessToken: string;
      refreshToken: string;
      user: AuthUser;
    };

    return {
      accessToken: parsed.accessToken,
      refreshToken: parsed.refreshToken,
      user: parsed.user,
      isAuthenticated: true
    };
  } catch {
    window.sessionStorage.removeItem(SESSION_STORAGE_KEY);
    return {
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false
    };
  }
};

const initialState = restoreSession();

const OBSERVABILITY_READ_ROLES: PlatformRole[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly'
];
const OBSERVABILITY_MANAGE_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator'];
const WORKLOAD_OPS_READ_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator', 'readonly'];
const WORKLOAD_OPS_EXECUTE_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator'];
const WORKLOAD_OPS_TERMINAL_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator'];
const WORKLOAD_OPS_ROLLBACK_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator'];
const WORKLOAD_OPS_BATCH_ROLES: PlatformRole[] = ['platform-admin', 'ops-operator'];

const GITOPS_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'gitops:read'
];
const GITOPS_MANAGE_SOURCE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'gitops:manage-source'
];
const GITOPS_SYNC_IDENTIFIERS: readonly string[] = ['platform-admin', 'ops-operator', 'gitops:sync'];
const GITOPS_PROMOTE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'gitops:promote'
];
const GITOPS_ROLLBACK_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'gitops:rollback'
];
const GITOPS_OVERRIDE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'gitops:override'
];
const GITOPS_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'gitops:read'
];
const POLICY_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'policy:read',
  'securitypolicy:read'
];
const POLICY_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'policy:read',
  'securitypolicy:read'
];
const POLICY_MANAGE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'policy:manage',
  'securitypolicy:manage'
];
const POLICY_ASSIGN_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'policy:assign',
  'securitypolicy:enforce'
];
const POLICY_APPROVE_EXCEPTION_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'policy:approve-exception',
  'securitypolicy:manage',
  'securitypolicy:enforce'
];
const COMPLIANCE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'compliance:read'
];
const COMPLIANCE_MANAGE_BASELINE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'compliance:manage-baseline'
];
const COMPLIANCE_EXECUTE_SCAN_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'compliance:execute-scan'
];
const COMPLIANCE_MANAGE_REMEDIATION_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'compliance:manage-remediation'
];
const COMPLIANCE_REVIEW_EXCEPTION_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'compliance:review-exception'
];
const COMPLIANCE_EXPORT_ARCHIVE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'compliance:export-archive'
];
const COMPLIANCE_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'compliance:read',
  'compliance:export-archive'
];
const CLUSTER_LIFECYCLE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'clusterlifecycle:read'
];
const CLUSTER_LIFECYCLE_IMPORT_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'clusterlifecycle:import'
];
const CLUSTER_LIFECYCLE_CREATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'clusterlifecycle:create'
];
const CLUSTER_LIFECYCLE_UPGRADE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'clusterlifecycle:upgrade'
];
const CLUSTER_LIFECYCLE_NODEPOOL_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'clusterlifecycle:manage-nodepool'
];
const CLUSTER_LIFECYCLE_RETIRE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'clusterlifecycle:retire'
];
const CLUSTER_LIFECYCLE_DRIVER_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'clusterlifecycle:manage-driver'
];
const CLUSTER_LIFECYCLE_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'clusterlifecycle:read'
];
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
const BACKUP_RESTORE_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'backuprestore:read'
];
const IDENTITY_TENANCY_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'identity:read'
];
const IDENTITY_TENANCY_MANAGE_SOURCE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:manage-source'
];
const IDENTITY_TENANCY_MANAGE_ORG_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:manage-org'
];
const IDENTITY_TENANCY_MANAGE_ROLE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:manage-role'
];
const IDENTITY_TENANCY_DELEGATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:delegate'
];
const IDENTITY_TENANCY_SESSION_GOVERN_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'identity:session-govern'
];
const PLATFORM_MARKETPLACE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'marketplace:read'
];
const PLATFORM_MARKETPLACE_MANAGE_SOURCE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'marketplace:manage-source'
];
const PLATFORM_MARKETPLACE_PUBLISH_TEMPLATE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'marketplace:publish-template'
];
const PLATFORM_MARKETPLACE_MANAGE_EXTENSION_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'marketplace:manage-extension'
];
const PLATFORM_MARKETPLACE_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'marketplace:read'
];
const SRE_SCALE_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'sre:read'
];
const SRE_SCALE_MANAGE_HA_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'sre:manage-ha'
];
const SRE_SCALE_MANAGE_UPGRADE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'sre:manage-upgrade'
];
const SRE_SCALE_MANAGE_SCALE_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'sre:manage-scale'
];
const SRE_SCALE_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'sre:read'
];
const ENTERPRISE_POLISH_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'audit-reader',
  'readonly',
  'enterprise:read'
];
const ENTERPRISE_POLISH_MANAGE_AUDIT_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'enterprise:manage-audit'
];
const ENTERPRISE_POLISH_MANAGE_REPORTS_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'enterprise:manage-reports'
];
const ENTERPRISE_POLISH_MANAGE_DELIVERY_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'ops-operator',
  'enterprise:manage-delivery'
];
const ENTERPRISE_POLISH_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'enterprise:read'
];
const IDENTITY_TENANCY_AUDIT_READ_IDENTIFIERS: readonly string[] = [
  'platform-admin',
  'audit-reader',
  'identity:read'
];

const getUserRoles = (user: AuthUser | null | undefined): string[] => {
  if (!user || !Array.isArray(user.platformRoles)) {
    return [];
  }
  return user.platformRoles;
};

export const hasAnyRole = (
  user: AuthUser | null | undefined,
  expectedRoles: readonly string[]
): boolean => {
  if (!user) {
    return false;
  }
  const roles = getUserRoles(user);
  if (roles.length === 0) {
    return false;
  }
  return expectedRoles.some((role) => roles.includes(role));
};

export const canReadObservability = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, OBSERVABILITY_READ_ROLES);

export const canManageObservability = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, OBSERVABILITY_MANAGE_ROLES);

export const canReadWorkloadOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, WORKLOAD_OPS_READ_ROLES);

export const canExecuteWorkloadOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, WORKLOAD_OPS_EXECUTE_ROLES);

export const canAccessWorkloadOpsTerminal = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, WORKLOAD_OPS_TERMINAL_ROLES);

export const canRollbackWorkloadOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, WORKLOAD_OPS_ROLLBACK_ROLES);

export const canBatchWorkloadOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, WORKLOAD_OPS_BATCH_ROLES);

export const canReadGitOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_READ_IDENTIFIERS);

export const canManageGitOpsSource = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_MANAGE_SOURCE_IDENTIFIERS);

export const canSyncGitOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_SYNC_IDENTIFIERS);

export const canPromoteGitOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_PROMOTE_IDENTIFIERS);

export const canRollbackGitOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_ROLLBACK_IDENTIFIERS);

export const canOverrideGitOps = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_OVERRIDE_IDENTIFIERS);

export const canReadGitOpsAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, GITOPS_AUDIT_READ_IDENTIFIERS);

export const canReadPolicy = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, POLICY_READ_IDENTIFIERS);

export const canReadPolicyAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, POLICY_AUDIT_READ_IDENTIFIERS);

export const canManagePolicy = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, POLICY_MANAGE_IDENTIFIERS);

export const canAssignPolicy = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, POLICY_ASSIGN_IDENTIFIERS);

export const canApprovePolicyException = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, POLICY_APPROVE_EXCEPTION_IDENTIFIERS);

export const canReadCompliance = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_READ_IDENTIFIERS);

export const canManageComplianceBaseline = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_MANAGE_BASELINE_IDENTIFIERS);

export const canExecuteComplianceScan = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_EXECUTE_SCAN_IDENTIFIERS);

export const canManageComplianceRemediation = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_MANAGE_REMEDIATION_IDENTIFIERS);

export const canReviewComplianceException = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_REVIEW_EXCEPTION_IDENTIFIERS);

export const canExportComplianceArchive = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_EXPORT_ARCHIVE_IDENTIFIERS);

export const canReadComplianceAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, COMPLIANCE_AUDIT_READ_IDENTIFIERS);

export const canReadClusterLifecycle = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_READ_IDENTIFIERS);

export const canImportClusterLifecycle = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_IMPORT_IDENTIFIERS);

export const canCreateClusterLifecycle = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_CREATE_IDENTIFIERS);

export const canUpgradeClusterLifecycle = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_UPGRADE_IDENTIFIERS);

export const canManageClusterLifecycleNodePool = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_NODEPOOL_IDENTIFIERS);

export const canRetireClusterLifecycle = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_RETIRE_IDENTIFIERS);

export const canManageClusterLifecycleDriver = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_DRIVER_IDENTIFIERS);

export const canReadClusterLifecycleAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, CLUSTER_LIFECYCLE_AUDIT_READ_IDENTIFIERS);

export const canReadBackupRestore = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_READ_IDENTIFIERS);

export const canManageBackupPolicy = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_MANAGE_POLICY_IDENTIFIERS);

export const canRunBackup = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_BACKUP_IDENTIFIERS);

export const canExecuteBackupRestore = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_RESTORE_IDENTIFIERS);

export const canMigrateBackupRestore = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_MIGRATE_IDENTIFIERS);

export const canDrillBackupRestore = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_DRILL_IDENTIFIERS);

export const canReadBackupRestoreAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, BACKUP_RESTORE_AUDIT_READ_IDENTIFIERS);

export const canReadIdentityTenancy = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_READ_IDENTIFIERS);

export const canManageIdentitySource = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_MANAGE_SOURCE_IDENTIFIERS);

export const canManageIdentityOrg = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_MANAGE_ORG_IDENTIFIERS);

export const canManageIdentityRole = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_MANAGE_ROLE_IDENTIFIERS);

export const canDelegateIdentity = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_DELEGATE_IDENTIFIERS);

export const canGovernIdentitySession = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_SESSION_GOVERN_IDENTIFIERS);

export const canReadIdentityTenancyAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, IDENTITY_TENANCY_AUDIT_READ_IDENTIFIERS);

export const canReadPlatformMarketplace = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, PLATFORM_MARKETPLACE_READ_IDENTIFIERS);

export const canManagePlatformMarketplaceSource = (
  user: AuthUser | null | undefined
): boolean => hasAnyRole(user, PLATFORM_MARKETPLACE_MANAGE_SOURCE_IDENTIFIERS);

export const canPublishPlatformMarketplaceTemplate = (
  user: AuthUser | null | undefined
): boolean => hasAnyRole(user, PLATFORM_MARKETPLACE_PUBLISH_TEMPLATE_IDENTIFIERS);

export const canManagePlatformMarketplaceExtension = (
  user: AuthUser | null | undefined
): boolean => hasAnyRole(user, PLATFORM_MARKETPLACE_MANAGE_EXTENSION_IDENTIFIERS);

export const canReadPlatformMarketplaceAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, PLATFORM_MARKETPLACE_AUDIT_READ_IDENTIFIERS);

export const canReadSREScale = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, SRE_SCALE_READ_IDENTIFIERS);

export const canManageSREHA = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, SRE_SCALE_MANAGE_HA_IDENTIFIERS);

export const canManageSREUpgrade = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, SRE_SCALE_MANAGE_UPGRADE_IDENTIFIERS);

export const canManageSREScale = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, SRE_SCALE_MANAGE_SCALE_IDENTIFIERS);

export const canReadSREAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, SRE_SCALE_AUDIT_READ_IDENTIFIERS);

export const canReadEnterprisePolish = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, ENTERPRISE_POLISH_READ_IDENTIFIERS);

export const canManageEnterpriseAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, ENTERPRISE_POLISH_MANAGE_AUDIT_IDENTIFIERS);

export const canManageEnterpriseReports = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, ENTERPRISE_POLISH_MANAGE_REPORTS_IDENTIFIERS);

export const canManageEnterpriseDelivery = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, ENTERPRISE_POLISH_MANAGE_DELIVERY_IDENTIFIERS);

export const canReadEnterpriseAudit = (user: AuthUser | null | undefined): boolean =>
  hasAnyRole(user, ENTERPRISE_POLISH_AUDIT_READ_IDENTIFIERS);

export const useAuthStore = create<AuthState>((set) => ({
  ...initialState,
  setSession: ({ accessToken, refreshToken, user }) => {
    if (canUseStorage) {
      window.sessionStorage.setItem(
        SESSION_STORAGE_KEY,
        JSON.stringify({ accessToken, refreshToken, user })
      );
    }

    set({
      accessToken,
      refreshToken,
      user,
      isAuthenticated: true
    });
  },
  clearSession: () => {
    if (canUseStorage) {
      window.sessionStorage.removeItem(SESSION_STORAGE_KEY);
    }

    set({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false
    });
  }
}));

export const getAccessToken = (): string | null => useAuthStore.getState().accessToken;
