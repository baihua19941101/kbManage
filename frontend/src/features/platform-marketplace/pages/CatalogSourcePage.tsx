import { useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { CatalogSourceDrawer } from '@/features/platform-marketplace/components/CatalogSourceDrawer';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { useTemplateDistributionAction } from '@/features/platform-marketplace/hooks/useTemplateDistributionAction';
import { normalizeApiError } from '@/services/api/client';
import {
  catalogSourceQueryScope,
  createCatalogSource,
  listCatalogSources,
  platformMarketplaceQueryKeys,
  type CatalogSource
} from '@/services/platformMarketplace';

const columns = (onSync: (sourceId: string) => void): ColumnsType<CatalogSource> => [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '类型', dataIndex: 'sourceType', key: 'sourceType', render: (value?: string) => value || '—' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> },
  {
    title: '同步状态',
    dataIndex: 'syncStatus',
    key: 'syncStatus',
    render: (value?: string) => <StatusTag value={value} />
  },
  { title: '模板数', dataIndex: 'templateCount', key: 'templateCount', render: (value?: number) => value ?? 0 },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" onClick={() => onSync(record.id)}>
        触发同步
      </Button>
    )
  }
];

export const CatalogSourcePage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = usePlatformMarketplacePermissions();
  const distributionAction = useTemplateDistributionAction();
  const createMutation = useMutation({ mutationFn: createCatalogSource });
  const sourcesQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.catalogSources(catalogSourceQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listCatalogSources({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无应用目录访问权限。" />;
  }

  const items = sourcesQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="目录来源"
        description="维护应用目录来源、同步状态和可见模板集合，作为模板中心的可信入口。"
        actions={
          <Button type="primary" disabled={!permissions.canManageSource} onClick={() => setDrawerOpen(true)}>
            新增目录来源
          </Button>
        }
      />

      <Alert
        type="info"
        showIcon
        message="目录同步建议走受控镜像源"
        description="当前页面仅按契约展示目录元数据；主线程接线后可衔接后端同步日志和镜像代理配置。"
      />

      {sourcesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="目录来源加载失败"
          description={normalizeApiError(sourcesQuery.error, '目录来源加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`目录来源（${items.length}）`}>
        <Table<CatalogSource>
          rowKey={(record) => record.id}
          columns={columns((sourceId) => distributionAction.syncSourceMutation.mutate(sourceId))}
          dataSource={items}
          loading={sourcesQuery.isLoading || sourcesQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <CatalogSourceDrawer
        open={drawerOpen}
        submitting={createMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createMutation.mutate(payload, {
            onSuccess: () => {
              setDrawerOpen(false);
            }
          })
        }
      />
    </Space>
  );
};
