import { create } from 'zustand';

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
