import { create } from 'zustand';

export type PlatformRole = 'platform-admin' | 'ops-operator' | 'audit-reader' | 'readonly';

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
