import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canManageClusterLifecycleDriver,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { ClusterTemplateFormDrawer } from '@/features/cluster-lifecycle/components/ClusterTemplateFormDrawer';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  listClusterDrivers,
  listClusterTemplates,
  type ClusterTemplate
} from '@/services/clusterLifecycle';

const columns: ColumnsType<ClusterTemplate> = [
  { title: '模板名称', dataIndex: 'name', key: 'name' },
  { title: '基础设施类型', dataIndex: 'infrastructureType', key: 'infrastructureType' },
  { title: '驱动键', dataIndex: 'driverKey', key: 'driverKey' },
  { title: '版本范围', dataIndex: 'driverVersionRange', key: 'driverVersionRange' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  }
];

export const ClusterTemplatePage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canManage = canManageClusterLifecycleDriver(user);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { createTemplateMutation } = useLifecycleAction();

  const templatesQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.templates(),
    enabled: canRead,
    queryFn: () => listClusterTemplates()
  });
  const driversQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.drivers('template-page'),
    enabled: canRead,
    queryFn: () => listClusterDrivers()
  });

  const driverOptions = useMemo(
    () =>
      (driversQuery.data?.items || []).map((item) => ({
        label: item.displayName || item.driverKey,
        value: item.driverKey
      })),
    [driversQuery.data?.items]
  );

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无模板管理访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="模板资产管理"
        description="维护模板与驱动版本范围、能力要求之间的关系。"
        actions={
          <Button type="primary" disabled={!canManage} onClick={() => setDrawerOpen(true)}>
            新建模板
          </Button>
        }
      />

      {templatesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板列表加载失败"
          description={normalizeApiError(templatesQuery.error, '模板列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`模板资产（${templatesQuery.data?.items.length ?? 0}）`}>
        <Table<ClusterTemplate>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={templatesQuery.data?.items || []}
          loading={templatesQuery.isLoading || templatesQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <ClusterTemplateFormDrawer
        open={drawerOpen}
        driverOptions={driverOptions}
        submitting={createTemplateMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createTemplateMutation.mutate(payload, { onSuccess: () => setDrawerOpen(false) })
        }
      />
    </Space>
  );
};
