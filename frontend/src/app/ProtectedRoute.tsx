import { Spin } from 'antd';
import { useEffect, useState } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import {
  canReadGitOps,
  canReadWorkloadOps,
  canReadObservability,
  hasAnyRole,
  useAuthStore
} from '@/features/auth/store';
import { refreshSession } from '@/services/auth';

type RouteGuard = {
  pathPrefix: string;
  canAccess: (user: ReturnType<typeof useAuthStore.getState>['user']) => boolean;
};

const routeGuards: RouteGuard[] = [
  { pathPrefix: '/workspaces', canAccess: (user) => hasAnyRole(user, ['platform-admin']) },
  { pathPrefix: '/projects', canAccess: (user) => hasAnyRole(user, ['platform-admin', 'ops-operator']) },
  { pathPrefix: '/audit-events', canAccess: (user) => hasAnyRole(user, ['platform-admin', 'audit-reader']) },
  {
    pathPrefix: '/observability',
    canAccess: canReadObservability
  },
  {
    pathPrefix: '/workload-ops',
    canAccess: canReadWorkloadOps
  },
  {
    pathPrefix: '/gitops',
    canAccess: canReadGitOps
  }
];

export const ProtectedRoute = () => {
  const { isAuthenticated, accessToken, refreshToken, clearSession, setSession, user } =
    useAuthStore();
  const location = useLocation();
  const [refreshChecked, setRefreshChecked] = useState(false);

  useEffect(() => {
    let cancelled = false;

    const runRefresh = async () => {
      // Avoid eager refresh when access token is already present.
      // In React StrictMode (dev), effects can run twice and aggressive refresh
      // may revoke the just-issued refresh token, forcing users to login again.
      if (!isAuthenticated || !refreshToken || accessToken) {
        setRefreshChecked(true);
        return;
      }

      try {
        const session = await refreshSession(refreshToken);
        if (!cancelled) {
          setSession({
            accessToken: session.accessToken,
            refreshToken: session.refreshToken,
            user: session.user ?? user!
          });
        }
      } catch {
        if (!cancelled) {
          clearSession();
        }
      } finally {
        if (!cancelled) {
          setRefreshChecked(true);
        }
      }
    };

    void runRefresh();

    return () => {
      cancelled = true;
    };
  }, [accessToken, clearSession, isAuthenticated, refreshToken, setSession, user]);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!refreshChecked) {
    return (
      <div
        style={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}
      >
        <Spin size="large" tip="正在恢复会话..." />
      </div>
    );
  }

  const matchedGuard = routeGuards.find((guard) => location.pathname.startsWith(guard.pathPrefix));
  if (matchedGuard && !matchedGuard.canAccess(user)) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
};
