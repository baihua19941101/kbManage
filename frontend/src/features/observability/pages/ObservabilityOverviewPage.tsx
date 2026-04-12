import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Empty, Space, Typography } from 'antd';
import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { canReadObservability, useAuthStore } from '@/features/auth/store';
import { OverviewCards } from '@/features/observability/components/OverviewCards';
import { MetricsTrendChart } from '@/features/observability/components/MetricsTrendChart';
import { ApiError, normalizeApiError } from '@/services/api/client';
import { getObservabilityOverview } from '@/services/observability/overview';
import { queryMetricSeries } from '@/services/observability/metrics';

const isAuthorizationError = (error: unknown): boolean =>
  error instanceof ApiError && (error.status === 401 || error.status === 403);

export const ObservabilityOverviewPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadObservability(user);
  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const clusterId = searchParams.get('clusterId') ?? undefined;
  const namespace = searchParams.get('namespace') ?? undefined;
  const resourceName = searchParams.get('resourceName') ?? undefined;
  const startAt = searchParams.get('startAt') ?? undefined;
  const endAt = searchParams.get('endAt') ?? undefined;

  const metricSubjectType = resourceName ? 'workload' : 'cluster';
  const metricSubjectRef = resourceName || clusterId || 'all-clusters';

  const overviewQuery = useQuery({
    queryKey: ['observability', 'overview', { clusterId, startAt, endAt }],
    queryFn: () => getObservabilityOverview({ clusterId, startAt, endAt }),
    enabled: canRead
  });
  const seriesQuery = useQuery({
    queryKey: ['observability', 'metrics', metricSubjectType, metricSubjectRef, startAt, endAt],
    queryFn: () =>
      queryMetricSeries({
        clusterId,
        namespace,
        subjectType: metricSubjectType,
        subjectRef: metricSubjectRef,
        metricKey: 'cpu_usage',
        startAt,
        endAt
      }),
    enabled: canRead && Boolean(metricSubjectRef)
  });

  const overviewAuthError = isAuthorizationError(overviewQuery.error);
  const seriesAuthError = isAuthorizationError(seriesQuery.error);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无可观测访问权限，请联系管理员授予工作空间/项目范围。"
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            可观测中心
          </Typography.Title>
          <Typography.Text type="secondary">
            统一查看日志、事件、指标和告警的入口。
          </Typography.Text>
          {clusterId ? (
            <Typography.Text type="secondary" style={{ display: 'block' }}>
              当前集群：{clusterId}
            </Typography.Text>
          ) : null}
        </div>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability/logs${location.search}`)}>日志</Button>
          <Button onClick={() => void navigate(`/observability/events${location.search}`)}>事件</Button>
          <Button onClick={() => void navigate(`/observability/metrics${location.search}`)}>指标</Button>
          <Button onClick={() => void navigate(`/observability/context${location.search}`)}>资源上下文</Button>
          <Button onClick={() => void navigate('/observability/alerts')}>告警中心</Button>
          <Button onClick={() => void navigate('/observability/alert-rules')}>规则治理</Button>
          <Button onClick={() => void navigate('/observability/silences')}>静默窗口</Button>
        </Space>
      </Space>

      {overviewQuery.error && !overviewAuthError ? (
        <Alert
          type="error"
          showIcon
          message="概览加载失败"
          description={normalizeApiError(overviewQuery.error, '概览加载失败，请稍后重试。')}
        />
      ) : null}
      {overviewAuthError ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description={normalizeApiError(
            overviewQuery.error,
            '当前账号的可观测权限可能已被回收，请刷新页面或重新登录后重试。'
          )}
        />
      ) : null}

      <OverviewCards cards={overviewQuery.data?.cards ?? []} />

      <div>
        <Typography.Title level={5} style={{ marginBottom: 12 }}>
          CPU 趋势（{metricSubjectRef}）
        </Typography.Title>
        {seriesQuery.error && !seriesAuthError ? (
          <Alert
            type="warning"
            showIcon
            message="指标趋势加载失败"
            description={normalizeApiError(seriesQuery.error, '指标趋势查询失败，请稍后重试。')}
            style={{ marginBottom: 16 }}
          />
        ) : null}
        {seriesAuthError ? (
          <Alert
            type="warning"
            showIcon
            message="权限已变更"
            description={normalizeApiError(
              seriesQuery.error,
              '当前账号的指标访问权限可能已被回收，请刷新页面后重试。'
            )}
            style={{ marginBottom: 16 }}
          />
        ) : null}
        {seriesQuery.data ? <MetricsTrendChart series={seriesQuery.data} /> : null}
      </div>
    </Space>
  );
};
