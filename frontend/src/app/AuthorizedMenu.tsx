import { Menu } from 'antd';
import type { MenuProps } from 'antd';
import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  canDrillBackupRestore,
  canManageBackupPolicy,
  canReadBackupRestore,
  canReadBackupRestoreAudit,
  canManageIdentityOrg,
  canManageEnterpriseAudit,
  canManageEnterpriseDelivery,
  canManageEnterpriseReports,
  canManageIdentityRole,
  canManageIdentitySource,
  canManagePlatformMarketplaceExtension,
  canManagePlatformMarketplaceSource,
  canManageSREHA,
  canManageSREScale,
  canManageSREUpgrade,
  canPublishPlatformMarketplaceTemplate,
  canReadIdentityTenancy,
  canReadEnterpriseAudit,
  canReadEnterprisePolish,
  canReadIdentityTenancyAudit,
  canReadPlatformMarketplace,
  canReadPlatformMarketplaceAudit,
  canReadSREAudit,
  canReadSREScale,
  canBatchWorkloadOps,
  canCreateClusterLifecycle,
  canExportComplianceArchive,
  canManageClusterLifecycleDriver,
  canManageObservability,
  canReadClusterLifecycle,
  canReadClusterLifecycleAudit,
  canReadCompliance,
  canReadComplianceAudit,
  canReadGitOps,
  canReadGitOpsAudit,
  canReadObservability,
  canReadPolicy,
  canReadPolicyAudit,
  canReadWorkloadOps,
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
  {
    key: '/workspaces',
    label: '工作空间',
    visibleWhen: (user) => hasAnyRole(user, ['platform-admin'])
  },
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
    key: '/audit-events/gitops',
    label: 'GitOps审计',
    visibleWhen: canReadGitOpsAudit
  },
  {
    key: '/audit-events/security-policy',
    label: '策略审计',
    visibleWhen: canReadPolicyAudit
  },
  {
    key: '/audit-events/compliance',
    label: '合规审计',
    visibleWhen: canReadComplianceAudit
  },
  {
    key: '/audit-events/cluster-lifecycle',
    label: '生命周期审计',
    visibleWhen: canReadClusterLifecycleAudit
  },
  {
    key: '/audit-events/backup-restore',
    label: '备份恢复审计',
    visibleWhen: canReadBackupRestoreAudit
  },
  {
    key: '/audit-events/identity',
    label: '身份治理审计',
    visibleWhen: canReadIdentityTenancyAudit
  },
  {
    key: '/audit-events/platform-marketplace',
    label: '市场审计',
    visibleWhen: canReadPlatformMarketplaceAudit
  },
  {
    key: '/audit-events/sre',
    label: 'SRE审计',
    visibleWhen: canReadSREAudit
  },
  {
    key: '/audit-events/enterprise',
    label: '企业治理审计',
    visibleWhen: canReadEnterpriseAudit
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
  },
  {
    key: '/workload-ops',
    label: '工作负载运维',
    visibleWhen: canReadWorkloadOps
  },
  {
    key: '/workload-ops/batches',
    label: '批量任务',
    visibleWhen: canBatchWorkloadOps
  },
  {
    key: '/gitops',
    label: 'GitOps 发布',
    visibleWhen: canReadGitOps
  },
  {
    key: '/security-policies',
    label: '安全策略',
    visibleWhen: canReadPolicy
  },
  {
    key: '/security-policies/rollout',
    label: '策略灰度',
    visibleWhen: canReadPolicy
  },
  {
    key: '/security-policies/violations',
    label: '违规中心',
    visibleWhen: canReadPolicy
  },
  {
    key: '/cluster-lifecycle',
    label: '集群生命周期',
    visibleWhen: canReadClusterLifecycle
  },
  {
    key: '/cluster-lifecycle/provision',
    label: '集群创建',
    visibleWhen: canCreateClusterLifecycle
  },
  {
    key: '/cluster-lifecycle/drivers',
    label: '驱动管理',
    visibleWhen: canManageClusterLifecycleDriver
  },
  {
    key: '/cluster-lifecycle/templates',
    label: '模板管理',
    visibleWhen: canReadClusterLifecycle
  },
  {
    key: '/cluster-lifecycle/capabilities',
    label: '能力矩阵',
    visibleWhen: canReadClusterLifecycle
  },
  {
    key: '/backup-restore',
    label: '备份恢复中心',
    visibleWhen: canReadBackupRestore
  },
  {
    key: '/backup-restore/policies',
    label: '备份策略',
    visibleWhen: canManageBackupPolicy
  },
  {
    key: '/backup-restore/restore-jobs',
    label: '恢复迁移',
    visibleWhen: canReadBackupRestore
  },
  {
    key: '/backup-restore/drills',
    label: '灾备演练',
    visibleWhen: canDrillBackupRestore
  },
  {
    key: '/identity-tenancy',
    label: '身份租户治理',
    visibleWhen: canReadIdentityTenancy
  },
  {
    key: '/identity-tenancy/organizations',
    label: '组织模型',
    visibleWhen: canManageIdentityOrg
  },
  {
    key: '/identity-tenancy/roles',
    label: '角色授权',
    visibleWhen: canManageIdentityRole
  },
  {
    key: '/identity-tenancy/sources',
    label: '身份源',
    visibleWhen: canManageIdentitySource
  },
  {
    key: '/platform-marketplace',
    label: '应用目录市场',
    visibleWhen: canReadPlatformMarketplace
  },
  {
    key: '/platform-marketplace/catalog-sources',
    label: '目录来源',
    visibleWhen: canManagePlatformMarketplaceSource
  },
  {
    key: '/platform-marketplace/distribution',
    label: '模板分发',
    visibleWhen: canPublishPlatformMarketplaceTemplate
  },
  {
    key: '/platform-marketplace/extensions',
    label: '扩展中心',
    visibleWhen: canManagePlatformMarketplaceExtension
  },
  {
    key: '/sre-scale',
    label: '平台SRE',
    visibleWhen: canReadSREScale
  },
  {
    key: '/sre-scale/ha',
    label: '高可用治理',
    visibleWhen: canManageSREHA
  },
  {
    key: '/enterprise-polish',
    label: '企业治理收尾',
    visibleWhen: canReadEnterprisePolish
  },
  {
    key: '/enterprise-polish/reports',
    label: '治理报表',
    visibleWhen: canManageEnterpriseReports
  },
  {
    key: '/enterprise-polish/delivery',
    label: '交付就绪',
    visibleWhen: canManageEnterpriseDelivery
  },
  {
    key: '/enterprise-polish/audit',
    label: '深度审计',
    visibleWhen: canManageEnterpriseAudit
  },
  {
    key: '/sre-scale/upgrades',
    label: '升级治理',
    visibleWhen: canManageSREUpgrade
  },
  {
    key: '/sre-scale/capacity',
    label: '容量性能',
    visibleWhen: canManageSREScale
  },
  {
    key: '/compliance-hardening/baselines',
    label: '合规基线',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/scans',
    label: '扫描中心',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/remediation',
    label: '整改工作台',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/exceptions',
    label: '例外审批',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/rechecks',
    label: '复检中心',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/overview',
    label: '合规总览',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/trends',
    label: '趋势复盘',
    visibleWhen: canReadCompliance
  },
  {
    key: '/compliance-hardening/archive',
    label: '归档导出',
    visibleWhen: (user) => canReadCompliance(user) || canExportComplianceArchive(user)
  }
];

const findBestSelectedKey = (pathname: string, keys: string[]): string => {
  if (keys.includes(pathname)) {
    return pathname;
  }

  return (
    keys
      .filter((key) => pathname.startsWith(key) && key !== '/')
      .sort((a, b) => b.length - a.length)[0] ?? '/'
  );
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
