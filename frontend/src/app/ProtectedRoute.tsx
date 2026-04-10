import { Spin } from 'antd';
import { useEffect, useState } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { refreshSession } from '@/services/auth';

export const ProtectedRoute = () => {
  const { isAuthenticated, refreshToken, clearSession, setSession, user } =
    useAuthStore();
  const [refreshChecked, setRefreshChecked] = useState(false);

  useEffect(() => {
    let cancelled = false;

    const runRefresh = async () => {
      if (!isAuthenticated || !refreshToken) {
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
  }, [clearSession, isAuthenticated, refreshToken, setSession, user]);

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

  return <Outlet />;
};
