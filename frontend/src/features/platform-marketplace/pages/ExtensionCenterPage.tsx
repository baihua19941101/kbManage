import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { ExtensionPackageDrawer } from '@/features/platform-marketplace/components/ExtensionPackageDrawer';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { useExtensionAction } from '@/features/platform-marketplace/hooks/useExtensionAction';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  extensionQueryScope,
  hasPermissionDeclaration,
  listExtensions,
  platformMarketplaceQueryKeys,
  type ExtensionPackage
} from '@/services/platformMarketplace';

const columns = (
  onEnable: (extensionId: string) => void,
  onDisable: (extensionId: string) => void
): ColumnsType<ExtensionPackage> => [
  { title: '扩展名称', dataIndex: 'name', key: 'name' },
  { title: '版本', dataIndex: 'version', key: 'version', render: (value?: string) => value || '—' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> },
  {
    title: '兼容性',
    dataIndex: 'compatibilityStatus',
    key: 'compatibilityStatus',
    render: (value?: string) => <StatusTag value={value} />
  },
  {
    title: '权限声明',
    dataIndex: 'permissionSummary',
    key: 'permissionSummary',
    render: (value: string | undefined, record) => (hasPermissionDeclaration(record) ? value : '未声明')
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Space size={4}>
        <Button type="link" onClick={() => onEnable(record.id)}>
          启用
        </Button>
        <Button type="link" danger onClick={() => onDisable(record.id)}>
          停用
        </Button>
      </Space>
    )
  }
];

export const ExtensionCenterPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const action = useExtensionAction();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const extensionsQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.extensions(extensionQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listExtensions({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无扩展中心访问权限。" />;
  }

  const items = extensionsQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="扩展中心"
        description="注册平台扩展、查看权限声明和兼容性状态，并控制启停影响范围。"
        actions={
          <Button type="primary" disabled={!permissions.canManageExtension} onClick={() => setDrawerOpen(true)}>
            注册扩展
          </Button>
        }
      />

      {extensionsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="扩展列表加载失败"
          description={normalizeApiError(extensionsQuery.error, '扩展列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Alert
        type="info"
        showIcon
        message="启停影响分析为契约占位"
        description="兼容性阻断提示、停用影响空态和审计页入口需要主线程补全全局接线后连通。"
      />

      <Card size="small" title={`扩展列表（${items.length}）`}>
        <Table<ExtensionPackage>
          rowKey={(record) => record.id}
          columns={columns(
            (extensionId) => action.enableExtensionMutation.mutate(extensionId),
            (extensionId) => action.disableExtensionMutation.mutate(extensionId)
          )}
          dataSource={items}
          loading={extensionsQuery.isLoading || extensionsQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <ExtensionPackageDrawer
        open={drawerOpen}
        submitting={action.createExtensionMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          action.createExtensionMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
