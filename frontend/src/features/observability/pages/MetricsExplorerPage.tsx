import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Form, Input, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { MetricsTrendChart } from '@/features/observability/components/MetricsTrendChart';
import { queryMetricSeries, type MetricSeriesParams } from '@/services/observability/metrics';

type MetricsFormValues = {
  clusterId?: string;
  namespace?: string;
  subjectType?: MetricSeriesParams['subjectType'];
  subjectRef?: string;
  metricKey?: string;
  startAt?: string;
  endAt?: string;
  step?: string;
};

const DEFAULT_METRICS_FILTERS: MetricSeriesParams = {
  subjectType: 'workload',
  subjectRef: 'mock-app',
  metricKey: 'cpu_usage'
};

const parseSubjectType = (value: string | null): MetricSeriesParams['subjectType'] | undefined => {
  switch (value) {
    case 'cluster':
    case 'node':
    case 'namespace':
    case 'workload':
    case 'pod':
      return value;
    default:
      return undefined;
  }
};

export const MetricsExplorerPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const initialFilters = useMemo<MetricSeriesParams>(
    () => ({
      clusterId: searchParams.get('clusterId') ?? undefined,
      namespace: searchParams.get('namespace') ?? undefined,
      subjectType: parseSubjectType(searchParams.get('subjectType')) ?? DEFAULT_METRICS_FILTERS.subjectType,
      subjectRef: searchParams.get('subjectRef') ?? searchParams.get('resourceName') ?? DEFAULT_METRICS_FILTERS.subjectRef,
      metricKey: searchParams.get('metricKey') ?? DEFAULT_METRICS_FILTERS.metricKey,
      startAt: searchParams.get('startAt') ?? undefined,
      endAt: searchParams.get('endAt') ?? undefined,
      step: searchParams.get('step') ?? undefined
    }),
    [searchParams]
  );
  const [filters, setFilters] = useState<MetricSeriesParams>(initialFilters);

  const metricsQuery = useQuery({
    queryKey: ['observability', 'metrics', filters],
    queryFn: () => queryMetricSeries(filters)
  });

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ marginBottom: 0 }}>
          指标趋势
        </Typography.Title>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
          <Button onClick={() => void navigate(`/observability/logs${location.search}`)}>日志</Button>
          <Button onClick={() => void navigate(`/observability/events${location.search}`)}>事件</Button>
        </Space>
      </Space>

      <Card>
        <Form<MetricsFormValues>
          layout="vertical"
          initialValues={initialFilters}
          onFinish={(values) => {
            setFilters({
              clusterId: values.clusterId,
              namespace: values.namespace,
              subjectType: values.subjectType ?? DEFAULT_METRICS_FILTERS.subjectType,
              subjectRef: values.subjectRef || DEFAULT_METRICS_FILTERS.subjectRef,
              metricKey: values.metricKey || DEFAULT_METRICS_FILTERS.metricKey,
              startAt: values.startAt,
              endAt: values.endAt,
              step: values.step
            });
          }}
        >
          <Space wrap style={{ width: '100%' }} align="start">
            <Form.Item name="clusterId" label="Cluster ID" style={{ width: 180 }}>
              <Input placeholder="cluster-1" />
            </Form.Item>
            <Form.Item name="namespace" label="Namespace" style={{ width: 180 }}>
              <Input placeholder="default" />
            </Form.Item>
            <Form.Item name="subjectType" label="主体类型" style={{ width: 180 }}>
              <Input placeholder="cluster/node/namespace/workload/pod" />
            </Form.Item>
            <Form.Item name="subjectRef" label="主体标识" style={{ width: 180 }}>
              <Input placeholder="mock-app" />
            </Form.Item>
          </Space>
          <Space wrap style={{ width: '100%' }} align="start">
            <Form.Item name="metricKey" label="指标键" style={{ width: 180 }}>
              <Input placeholder="cpu_usage" />
            </Form.Item>
            <Form.Item name="startAt" label="开始时间" style={{ width: 220 }}>
              <Input placeholder="2026-04-11T00:00:00Z" />
            </Form.Item>
            <Form.Item name="endAt" label="结束时间" style={{ width: 220 }}>
              <Input placeholder="2026-04-11T01:00:00Z" />
            </Form.Item>
            <Form.Item name="step" label="步长" style={{ width: 120 }}>
              <Input placeholder="1m" />
            </Form.Item>
            <Form.Item label=" ">
              <Button type="primary" htmlType="submit" loading={metricsQuery.isFetching}>
                查询
              </Button>
            </Form.Item>
          </Space>
        </Form>
      </Card>

      <Card title={`${filters.metricKey} · ${filters.subjectRef}`}>
        {metricsQuery.error ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="指标查询失败"
            description={`${metricsQuery.error}`}
          />
        ) : null}
        {metricsQuery.data ? <MetricsTrendChart series={metricsQuery.data} /> : null}
      </Card>
    </Space>
  );
};
