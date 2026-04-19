import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { UpgradeEntryDrawer } from '@/features/platform-marketplace/components/UpgradeEntryDrawer';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  hasUpgradeGap,
  installationQueryScope,
  listInstallations,
  platformMarketplaceQueryKeys,
  summarizeTargetScope,
  type InstallationRecord
} from '@/services/platformMarketplace';

const columns = (onOpen: (record: InstallationRecord) => void): ColumnsType<InstallationRecord> => [
  { title: '模板', dataIndex: 'templateName', key: 'templateName', render: (value?: string) => value || '—' },
  {
    title: '目标范围',
    key: 'target',
    render: (_, record) => summarizeTargetScope(record.targetType, record.targetRef)
  },
  {
    title: '当前版本',
    dataIndex: 'currentVersion',
    key: 'currentVersion',
    render: (value?: string) => value || '—'
  },
  {
    title: '最新版本',
    dataIndex: 'latestVersion',
    key: 'latestVersion',
    render: (value?: string) => value || '—'
  },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" onClick={() => onOpen(record)}>
        升级入口
      </Button>
    )
  }
];

export const InstallationRecordPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const [selected, setSelected] = useState<InstallationRecord>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const installationsQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.installations(installationQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listInstallations({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无安装记录访问权限。" />;
  }

  const items = installationsQuery.data?.items || [];
  const upgradeCount = items.filter(hasUpgradeGap).length;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="安装记录"
        description="查看模板安装、版本变化、升级入口和下线状态，保留已发生交付动作的历史。"
      />

      {installationsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="安装记录加载失败"
          description={normalizeApiError(installationsQuery.error, '安装记录加载失败，请稍后重试。')}
        />
      ) : null}

      <Alert
        type="info"
        showIcon
        message={`当前有 ${upgradeCount} 条安装记录存在可升级版本`}
        description="版本变化说明、下线提示和升级受限空态需要主线程补全全局路由和后端细粒度契约后进一步增强。"
      />

      <Card size="small" title={`安装记录（${items.length}）`}>
        <Table<InstallationRecord>
          rowKey={(record) => record.id}
          columns={columns((record) => {
            setSelected(record);
            setDrawerOpen(true);
          })}
          dataSource={items}
          loading={installationsQuery.isLoading || installationsQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <UpgradeEntryDrawer open={drawerOpen} record={selected} onClose={() => setDrawerOpen(false)} />
    </Space>
  );
};
