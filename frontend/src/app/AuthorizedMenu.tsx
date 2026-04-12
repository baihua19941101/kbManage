import { Menu } from 'antd';
import type { MenuProps } from 'antd';
import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  canManageObservability,
  canReadObservability,
  hasAnyRole,
  useAuthStore
} from '@/features/auth/store';

type MenuItemConfig = {
  key: string;
  label: string;
  visibleWhen?: (rolesUser: ReturnType<typeof useAuthStore.getState>['user']) => boolean;
};

const allMenuItems: MenuItemConfig[] = [
  { key: '/', label: '首页' },
  {
    key: '/clusters',
    label: '集群',
    visibleWhen: (user) => hasAnyRole(user, ['platform-admin', 'ops-operator', 'readonly'])
  },
  {
    key: '/resources',
    label: '资源',
    visibleWhen: (user) => hasAnyRole(user, ['platform-admin', 'ops-operator', 'readonly'])
  },
  { key: '/workspaces', label: '工作空间', visibleWhen: (user) => hasAnyRole(user, ['platform-admin']) },
  {
    key: '/projects',
    label: '项目',
    visibleWhen: (user) => hasAnyRole(user, ['platform-admin', 'ops-operator'])
  },
  {
    key: '/audit-events',
    label: '审计',
    visibleWhen: (user) => hasAnyRole(user, ['platform-admin', 'audit-reader'])
  },
  {
    key: '/observability',
    label: '可观测',
    visibleWhen: canReadObservability
  },
  {
    key: '/observability/alerts',
    label: '告警中心',
    visibleWhen: canReadObservability
  },
  {
    key: '/observability/alert-rules',
    label: '告警规则',
    visibleWhen: canReadObservability
  },
  {
    key: '/observability/silences',
    label: '静默窗口',
    visibleWhen: canManageObservability
  }
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
  const user = useAuthStore((state) => state.user);

  const authorizedItems = useMemo(() => {
    return allMenuItems.filter((item) => {
      if (!item.visibleWhen) {
        return true;
      }
      return item.visibleWhen(user);
    });
  }, [user]);

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
