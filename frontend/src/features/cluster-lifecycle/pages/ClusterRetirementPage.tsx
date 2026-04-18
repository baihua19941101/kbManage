import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Typography } from 'antd';
import { useLocation } from 'react-router-dom';
import {
  canReadClusterLifecycle,
  canRetireClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { RetireClusterDrawer } from '@/features/cluster-lifecycle/components/RetireClusterDrawer';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  extractClusterOperationState,
  getClusterLifecycleDetail,
  type DisableClusterRequest,
  type RetireClusterRequest
} from '@/services/clusterLifecycle';

export const ClusterRetirementPage = () => {
  const location = useLocation();
  const clusterId = useMemo(
    () => new URLSearchParams(location.search).get('clusterId') || '',
    [location.search]
  );
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canRetire = canRetireClusterLifecycle(user);
  const [mode, setMode] = useState<'disable' | 'retire'>('retire');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { disableMutation, retireMutation } = useLifecycleAction();

  const clusterQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.clusterDetail(clusterId),
    enabled: canRead && clusterId.length > 0,
    queryFn: () => getClusterLifecycleDetail(clusterId)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无停用/退役访问权限。" />;
  }

  const state = extractClusterOperationState(clusterQuery.data);

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="停用与退役"
        description="在完成确认后发起停用或退役流程，并保留审计复盘信息。"
      />

      {!clusterId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请通过 query string 提供 clusterId。" />
      ) : null}

      {clusterQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="集群状态加载失败"
          description={normalizeApiError(clusterQuery.error, '集群状态加载失败，请稍后重试。')}
        />
      ) : null}

      {!canRetire ? (
        <Alert type="warning" showIcon message="当前账号没有停用/退役动作权限。" />
      ) : null}

      {state !== 'idle' ? (
        <Alert
          type={state === 'running' ? 'warning' : 'info'}
          showIcon
          message={state === 'running' ? '当前集群存在进行中的高风险动作。' : '当前集群处于不可执行状态。'}
          description="前端已预留冲突动作空态，等待后端锁语义进一步接入。"
        />
      ) : null}

      <Card size="small" title="当前集群">
        <Typography.Text>
          {clusterQuery.data?.displayName || clusterQuery.data?.name || clusterId || '—'}
        </Typography.Text>
        <Typography.Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0 }}>
          当前状态：{clusterQuery.data?.status || '—'}，退役原因：
          {clusterQuery.data?.retirementReason || '未记录'}
        </Typography.Paragraph>
      </Card>

      <Space>
        <Button
          disabled={!canRetire || state !== 'idle'}
          onClick={() => {
            setMode('disable');
            setDrawerOpen(true);
          }}
        >
          发起停用
        </Button>
        <Button
          type="primary"
          danger
          disabled={!canRetire || state !== 'idle'}
          onClick={() => {
            setMode('retire');
            setDrawerOpen(true);
          }}
        >
          发起退役
        </Button>
      </Space>

      <RetireClusterDrawer
        open={drawerOpen}
        mode={mode}
        submitting={disableMutation.isPending || retireMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => {
          if (!clusterId) {
            return;
          }
          if (mode === 'disable') {
            disableMutation.mutate(
              { clusterId, payload: payload as DisableClusterRequest },
              { onSuccess: () => setDrawerOpen(false) }
            );
            return;
          }
          retireMutation.mutate(
            { clusterId, payload: payload as RetireClusterRequest },
            { onSuccess: () => setDrawerOpen(false) }
          );
        }}
      />
    </Space>
  );
};
