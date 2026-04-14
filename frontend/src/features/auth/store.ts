import { create } from 'zustand';

export type PlatformRole = 'platform-admin' | 'ops-operator' | 'audit-reader' | 'readonly';

export type GitOpsPermission =
  | 'gitops:read'
  | 'gitops:manage-source'
  | 'gitops:sync'
  | 'gitops:promote'
  | 'gitops:rollback'
  | 'gitops:override';

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
