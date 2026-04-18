import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useLocation } from 'react-router-dom';
import {
  canManageClusterLifecycleNodePool,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { NodePoolScaleDrawer } from '@/features/cluster-lifecycle/components/NodePoolScaleDrawer';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  listNodePools,
  normalizeNodePoolTarget,
  type NodePoolProfile
} from '@/services/clusterLifecycle';

export const NodePoolPage = () => {
  const location = useLocation();
  const clusterId = useMemo(
    () => new URLSearchParams(location.search).get('clusterId') || '',
    [location.search]
  );
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canManage = canManageClusterLifecycleNodePool(user);
  const [selectedNodePool, setSelectedNodePool] = useState<NodePoolProfile | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { scaleNodePoolMutation } = useLifecycleAction();

  const nodePoolsQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.nodePools(clusterId),
    enabled: canRead && clusterId.length > 0,
    queryFn: () => listNodePools(clusterId)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无节点池管理访问权限。" />;
  }

  const columns: ColumnsType<NodePoolProfile> = [
    { title: '节点池', dataIndex: 'name', key: 'name' },
    { title: '角色', dataIndex: 'role', key: 'role' },
    { title: '当前节点数', dataIndex: 'currentCount', key: 'currentCount' },
    { title: '目标节点数', dataIndex: 'desiredCount', key: 'desiredCount' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (value?: string) => <LifecycleStatusTag value={value} />
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Button
          type="link"
          disabled={!canManage}
          onClick={() => {
            setSelectedNodePool(record);
            setDrawerOpen(true);
          }}
        >
          调整容量
        </Button>
      )
    }
  ];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="节点池管理" description="查看节点池版本、容量和伸缩状态。" />

      {!clusterId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请通过 query string 提供 clusterId。" />
      ) : null}

      {!canManage ? (
        <Alert type="warning" showIcon message="当前账号没有节点池扩缩动作权限。" />
      ) : null}

      {nodePoolsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="节点池列表加载失败"
          description={normalizeApiError(nodePoolsQuery.error, '节点池列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`节点池列表（${nodePoolsQuery.data?.items.length ?? 0}）`}>
        <Table<NodePoolProfile>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={nodePoolsQuery.data?.items || []}
          loading={nodePoolsQuery.isLoading || nodePoolsQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <Typography.Text type="secondary">
        扩缩动作会在 clusterId 级别上走后端互斥控制；当前页面先展示前端入口和空态。
      </Typography.Text>

      <NodePoolScaleDrawer
        open={drawerOpen}
        initialDesiredCount={normalizeNodePoolTarget(selectedNodePool)}
        submitting={scaleNodePoolMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => {
          if (!clusterId || !selectedNodePool) {
            return;
          }
          scaleNodePoolMutation.mutate(
            { clusterId, nodePoolId: selectedNodePool.id, payload },
            { onSuccess: () => setDrawerOpen(false) }
          );
        }}
      />
    </Space>
  );
};
