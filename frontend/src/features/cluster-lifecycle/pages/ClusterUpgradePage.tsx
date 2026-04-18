import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useLocation } from 'react-router-dom';
import {
  canReadClusterLifecycle,
  canUpgradeClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import { UpgradePlanDrawer } from '@/features/cluster-lifecycle/components/UpgradePlanDrawer';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  getClusterLifecycleDetail,
  type UpgradePlan
} from '@/services/clusterLifecycle';

const columns: ColumnsType<UpgradePlan> = [
  { title: '计划 ID', dataIndex: 'id', key: 'id' },
  { title: '当前版本', dataIndex: 'fromVersion', key: 'fromVersion' },
  { title: '目标版本', dataIndex: 'toVersion', key: 'toVersion' },
  {
    title: '预检查',
    dataIndex: 'precheckStatus',
    key: 'precheckStatus',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  }
];

export const ClusterUpgradePage = () => {
  const location = useLocation();
  const clusterId = useMemo(
    () => new URLSearchParams(location.search).get('clusterId') || '',
    [location.search]
  );
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canUpgrade = canUpgradeClusterLifecycle(user);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [planId, setPlanId] = useState('');
  const { upgradePlanMutation, executeUpgradeMutation } = useLifecycleAction();

  const clusterQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.clusterDetail(clusterId),
    enabled: canRead && clusterId.length > 0,
    queryFn: () => getClusterLifecycleDetail(clusterId)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无集群升级访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="集群升级计划"
        description="为指定集群创建升级计划，并执行已审批的版本变更。"
        actions={
          <Button
            type="primary"
            disabled={!canUpgrade || !clusterId}
            onClick={() => setDrawerOpen(true)}
          >
            新建升级计划
          </Button>
        }
      />

      {!clusterId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请通过 query string 提供 clusterId。" />
      ) : null}

      {clusterQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="集群信息加载失败"
          description={normalizeApiError(clusterQuery.error, '集群信息加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="升级目标">
        <Typography.Text>
          目标集群：{clusterQuery.data?.displayName || clusterQuery.data?.name || clusterId || '—'}
        </Typography.Text>
        <Typography.Paragraph type="secondary" style={{ marginBottom: 0, marginTop: 8 }}>
          当前版本：{clusterQuery.data?.kubernetesVersion || '—'}，目标版本：
          {clusterQuery.data?.targetVersion || '待设置'}
        </Typography.Paragraph>
      </Card>

      <Card size="small" title="执行已创建计划">
        <Form
          layout="inline"
          onFinish={() =>
            clusterId && planId && executeUpgradeMutation.mutate({ clusterId, planId })
          }
        >
          <Form.Item>
            <Input
              value={planId}
              onChange={(event) => setPlanId(event.target.value)}
              placeholder="输入计划 ID"
            />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              disabled={!canUpgrade || !clusterId}
              loading={executeUpgradeMutation.isPending}
            >
              执行升级
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card size="small" title="升级计划占位列表">
        <Table<UpgradePlan>
          columns={columns}
          dataSource={upgradePlanMutation.data ? [upgradePlanMutation.data] : []}
          rowKey="id"
          pagination={false}
        />
      </Card>

      <UpgradePlanDrawer
        open={drawerOpen}
        submitting={upgradePlanMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => {
          if (!clusterId) {
            return;
          }
          upgradePlanMutation.mutate({ clusterId, payload }, { onSuccess: () => setDrawerOpen(false) });
        }}
      />
    </Space>
  );
};
