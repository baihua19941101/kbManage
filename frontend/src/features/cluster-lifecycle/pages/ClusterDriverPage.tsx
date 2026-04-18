import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canManageClusterLifecycleDriver,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { DriverVersionDrawer } from '@/features/cluster-lifecycle/components/DriverVersionDrawer';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  listClusterDrivers,
  type ClusterDriverVersion
} from '@/services/clusterLifecycle';

const columns: ColumnsType<ClusterDriverVersion> = [
  { title: '展示名称', key: 'displayName', render: (_, record) => record.displayName || record.driverKey },
  { title: '驱动键', dataIndex: 'driverKey', key: 'driverKey' },
  { title: '版本', dataIndex: 'version', key: 'version' },
  { title: '基础设施类型', dataIndex: 'providerType', key: 'providerType' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  }
];

export const ClusterDriverPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canManage = canManageClusterLifecycleDriver(user);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { createDriverMutation } = useLifecycleAction();

  const driversQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.drivers(),
    enabled: canRead,
    queryFn: () => listClusterDrivers()
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无驱动管理访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="驱动版本管理"
        description="维护集群驱动版本、基础设施类型和后续能力矩阵入口。"
        actions={
          <Button type="primary" disabled={!canManage} onClick={() => setDrawerOpen(true)}>
            登记驱动版本
          </Button>
        }
      />

      {driversQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="驱动版本加载失败"
          description={normalizeApiError(driversQuery.error, '驱动版本加载失败，请稍后重试。')}
        />
      ) : null}

      {!canManage ? <Alert type="warning" showIcon message="当前账号没有驱动管理权限。" /> : null}

      <Card size="small" title={`驱动版本（${driversQuery.data?.items.length ?? 0}）`}>
        <Table<ClusterDriverVersion>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={driversQuery.data?.items || []}
          loading={driversQuery.isLoading || driversQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <DriverVersionDrawer
        open={drawerOpen}
        submitting={createDriverMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createDriverMutation.mutate(payload, { onSuccess: () => setDrawerOpen(false) })
        }
      />
    </Space>
  );
};
