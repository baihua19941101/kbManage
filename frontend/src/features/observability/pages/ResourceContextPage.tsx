import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { EventTimeline } from '@/features/observability/components/EventTimeline';
import { LogTable } from '@/features/observability/components/LogTable';
import { MetricsTrendChart } from '@/features/observability/components/MetricsTrendChart';
import { ResourceContextPanel } from '@/features/observability/components/ResourceContextPanel';
import type { ObservabilityOverviewDTO } from '@/services/api/types';
import { listObservabilityEvents } from '@/services/observability/events';
import { queryObservabilityLogs } from '@/services/observability/logs';
import { queryMetricSeries } from '@/services/observability/metrics';
import { getObservabilityOverview } from '@/services/observability/overview';

const useQueryParams = () => {
  const { search } = useLocation();
  return useMemo(() => new URLSearchParams(search), [search]);
};

export const ResourceContextPage = () => {
  const navigate = useNavigate();
  const params = useQueryParams();
  const clusterId = params.get('clusterId') ?? '1';
  const namespace = params.get('namespace') ?? 'default';
  const resourceKind = params.get('resourceKind') ?? 'Deployment';
  const resourceName = params.get('resourceName') ?? 'mock-app';
  const startAt = params.get('startAt') ?? undefined;
  const endAt = params.get('endAt') ?? undefined;
  const queryString = params.toString();

  const overviewQuery = useQuery<ObservabilityOverviewDTO>({
    queryKey: ['observability', 'overview', clusterId, startAt, endAt],
    queryFn: () => getObservabilityOverview({ clusterId, startAt, endAt })
  });
  const eventsQuery = useQuery({
    queryKey: [
      'observability',
      'events',
      clusterId,
      namespace,
      resourceKind,
      resourceName,
      startAt,
      endAt
    ],
    queryFn: () =>
      listObservabilityEvents({
        clusterId,
        namespace,
        resourceKind,
        resourceName,
        startAt,
        endAt,
        limit: 100
      })
  });
  const logsQuery = useQuery({
    queryKey: ['observability', 'logs', clusterId, namespace, resourceKind, resourceName, startAt, endAt],
    queryFn: () =>
      queryObservabilityLogs({
        clusterId,
        namespace,
        resourceKind,
        resourceName,
        workload: resourceName,
        startAt,
        endAt,
        limit: 50
      })
  });
  const metricSeriesQuery = useQuery({
    queryKey: ['observability', 'metrics', clusterId, namespace, resourceKind, resourceName, startAt, endAt],
    queryFn: () =>
      queryMetricSeries({
        clusterId,
        namespace,
        subjectType: resourceKind.toLowerCase() === 'pod' ? 'pod' : 'workload',
        subjectRef: resourceName,
        metricKey: 'cpu_usage',
        startAt,
        endAt
      })
  });

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ marginBottom: 0 }}>
          资源可观测上下文
        </Typography.Title>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability${queryString ? `?${queryString}` : ''}`)}>
            总览
          </Button>
          <Button
            onClick={() => void navigate(`/observability/logs${queryString ? `?${queryString}` : ''}`)}
          >
            日志
          </Button>
          <Button
            onClick={() => void navigate(`/observability/events${queryString ? `?${queryString}` : ''}`)}
          >
            事件
          </Button>
          <Button
            onClick={() => void navigate(`/observability/metrics${queryString ? `?${queryString}` : ''}`)}
          >
            指标
          </Button>
        </Space>
      </Space>
      <ResourceContextPanel
        clusterId={clusterId}
        namespace={namespace}
        resourceKind={resourceKind}
        resourceName={resourceName}
        overview={overviewQuery.data}
      />
      <Card title="资源日志">
        {logsQuery.error ? (
          <Alert type="warning" showIcon message="日志加载失败" description={`${logsQuery.error}`} />
        ) : null}
        <LogTable loading={logsQuery.isFetching} items={logsQuery.data?.items ?? []} />
      </Card>
      <Card title="事件时间线">
        {eventsQuery.error ? (
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
            message="事件加载失败"
            description={`${eventsQuery.error}`}
          />
        ) : null}
        <EventTimeline items={eventsQuery.data?.items ?? []} />
      </Card>
      <Card title="指标趋势">
        {metricSeriesQuery.error ? (
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 16 }}
            message="指标加载失败"
            description={`${metricSeriesQuery.error}`}
          />
        ) : null}
        {metricSeriesQuery.data ? <MetricsTrendChart series={metricSeriesQuery.data} /> : null}
      </Card>
    </Space>
  );
};
