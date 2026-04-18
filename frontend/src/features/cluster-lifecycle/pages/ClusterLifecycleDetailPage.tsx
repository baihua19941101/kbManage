import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Descriptions, Empty, Space, Timeline, Typography } from 'antd';
import { useLocation, useParams } from 'react-router-dom';
import { canReadClusterLifecycle, useAuthStore } from '@/features/auth/store';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import {
  HealthStatusTag,
  LifecycleStatusTag
} from '@/features/cluster-lifecycle/components/status';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  getClusterLifecycleDetail,
  inferClusterDisplayName
} from '@/services/clusterLifecycle';

const useResolvedClusterId = () => {
  const params = useParams<{ clusterId?: string }>();
  const location = useLocation();
  const search = useMemo(() => new URLSearchParams(location.search), [location.search]);
  return params.clusterId || search.get('clusterId') || '';
};

export const ClusterLifecycleDetailPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const clusterId = useResolvedClusterId();
  const detailQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.clusterDetail(clusterId),
    enabled: canRead && clusterId.length > 0,
    queryFn: () => getClusterLifecycleDetail(clusterId)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无生命周期详情访问权限。" />;
  }

  if (!clusterId) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="缺少 clusterId，无法加载详情。" />;
  }

  const cluster = detailQuery.data;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title={inferClusterDisplayName(cluster)}
        description="查看接入模式、驱动版本、健康状态、最近校验和退役信息。"
      />

      {detailQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="生命周期详情加载失败"
          description={normalizeApiError(detailQuery.error, '生命周期详情加载失败，请稍后重试。')}
        />
      ) : null}

      <Card loading={detailQuery.isLoading}>
        <Descriptions column={2} size="small">
          <Descriptions.Item label="集群 ID">{cluster?.id || clusterId}</Descriptions.Item>
          <Descriptions.Item label="生命周期模式">{cluster?.lifecycleMode || '—'}</Descriptions.Item>
          <Descriptions.Item label="基础设施">{cluster?.infrastructureType || '—'}</Descriptions.Item>
          <Descriptions.Item label="驱动版本">{cluster?.driverVersion || '—'}</Descriptions.Item>
          <Descriptions.Item label="生命周期状态">
            <LifecycleStatusTag value={cluster?.status} />
          </Descriptions.Item>
          <Descriptions.Item label="健康状态">
            <HealthStatusTag value={cluster?.healthStatus} />
          </Descriptions.Item>
          <Descriptions.Item label="Kubernetes 版本">{cluster?.kubernetesVersion || '—'}</Descriptions.Item>
          <Descriptions.Item label="目标版本">{cluster?.targetVersion || '—'}</Descriptions.Item>
          <Descriptions.Item label="最近校验">{cluster?.lastValidationStatus || '—'}</Descriptions.Item>
          <Descriptions.Item label="最近校验时间">{cluster?.lastValidationAt || '—'}</Descriptions.Item>
          <Descriptions.Item label="节点池摘要">{cluster?.nodePoolSummary || '—'}</Descriptions.Item>
          <Descriptions.Item label="退役原因">{cluster?.retirementReason || '—'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card size="small" title="生命周期时间线">
        <Timeline
          items={[
            { color: 'blue', children: `创建时间：${cluster?.createdAt || '未知'}` },
            { color: 'green', children: `最近更新时间：${cluster?.updatedAt || '未知'}` },
            {
              color: cluster?.lastValidationStatus === 'failed' ? 'red' : 'gray',
              children: `最近校验：${cluster?.lastValidationStatus || '未知'}`
            }
          ]}
        />
        <Typography.Text type="secondary">
          页面已预留升级、节点池和退役详情联动区域，待后端动作查询接口补齐后可直接扩展。
        </Typography.Text>
      </Card>
    </Space>
  );
};
