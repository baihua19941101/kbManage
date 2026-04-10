import { Menu } from 'antd';
import type { MenuProps } from 'antd';
import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';

type MenuRole = 'platform-admin' | 'workspace-admin' | 'developer' | 'viewer';

type MenuItemConfig = {
  key: string;
  label: string;
  requiredRoles?: MenuRole[];
};

const allMenuItems: MenuItemConfig[] = [
  { key: '/', label: '首页' },
  { key: '/clusters', label: '集群' },
  { key: '/resources', label: '资源' },
  { key: '/workspaces', label: '工作空间', requiredRoles: ['platform-admin', 'workspace-admin'] },
  { key: '/projects', label: '项目', requiredRoles: ['platform-admin', 'workspace-admin', 'developer'] }
];

const findBestSelectedKey = (pathname: string, keys: string[]): string => {
  if (keys.includes(pathname)) {
    return pathname;
  }

  return keys
    .filter((key) => pathname.startsWith(key) && key !== '/')
    .sort((a, b) => b.length - a.length)[0] ?? '/';
};

export const AuthorizedMenu = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const user = useAuthStore((state) => state.user) as
    | (ReturnType<typeof useAuthStore.getState>['user'] & {
        roles?: string[];
        role?: string;
      })
    | null;

  const currentRoles = useMemo(() => {
    if (!user) {
      return [] as string[];
    }

    if (Array.isArray(user.roles) && user.roles.length > 0) {
      return user.roles;
    }

    if (typeof user.role === 'string' && user.role.length > 0) {
      return [user.role];
    }

    return [];
  }, [user]);

  const authorizedItems = useMemo(() => {
    const hasRoleData = currentRoles.length > 0;
    return allMenuItems.filter((item) => {
      if (!item.requiredRoles || item.requiredRoles.length === 0) {
        return true;
      }

      if (!hasRoleData) {
        return true;
      }

      return item.requiredRoles.some((role) => currentRoles.includes(role));
    });
  }, [currentRoles]);

  const menuItems = useMemo<NonNullable<MenuProps['items']>>(
    () =>
      authorizedItems.map((item) => ({
        key: item.key,
        label: item.label
      })),
    [authorizedItems]
  );

  const selectedKey = findBestSelectedKey(
    location.pathname,
    authorizedItems.map((item) => item.key)
  );

  return (
    <Menu
      mode="inline"
      selectedKeys={[selectedKey]}
      items={menuItems}
      onClick={({ key }) => {
        void navigate(String(key));
      }}
    />
  );
};
