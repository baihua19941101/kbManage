import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { LogFilters, type LogFilterValues } from '@/features/observability/components/LogFilters';
import { LogTable } from '@/features/observability/components/LogTable';
import { queryObservabilityLogs } from '@/services/observability/logs';

const toPositiveNumber = (value: string | null, fallback: number) => {
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
};

export const LogExplorerPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const initialFilters = useMemo<LogFilterValues>(
    () => ({
      clusterId: searchParams.get('clusterId') ?? undefined,
      namespace: searchParams.get('namespace') ?? undefined,
      resourceKind: searchParams.get('resourceKind') ?? undefined,
      resourceName: searchParams.get('resourceName') ?? undefined,
      workload: searchParams.get('workload') ?? searchParams.get('resourceName') ?? undefined,
      pod: searchParams.get('pod') ?? undefined,
      container: searchParams.get('container') ?? undefined,
      keyword: searchParams.get('keyword') ?? undefined,
      startAt: searchParams.get('startAt') ?? undefined,
      endAt: searchParams.get('endAt') ?? undefined,
      limit: toPositiveNumber(searchParams.get('limit'), 100)
    }),
    [searchParams]
  );
  const [filters, setFilters] = useState<LogFilterValues>(initialFilters);

  const logsQuery = useQuery({
    queryKey: ['observability', 'logs', filters],
    queryFn: () => queryObservabilityLogs({ ...filters, limit: filters.limit ?? 100 })
  });

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ marginBottom: 0 }}>
          日志检索
        </Typography.Title>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
          <Button onClick={() => void navigate(`/observability/events${location.search}`)}>事件</Button>
          <Button onClick={() => void navigate(`/observability/metrics${location.search}`)}>指标</Button>
        </Space>
      </Space>
      <Card>
        <LogFilters initialValues={initialFilters} loading={logsQuery.isFetching} onSearch={setFilters} />
      </Card>
      <Card>
        {logsQuery.error ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="日志查询失败"
            description={`${logsQuery.error}`}
          />
        ) : null}
        <LogTable loading={logsQuery.isFetching} items={logsQuery.data?.items ?? []} />
      </Card>
    </Space>
  );
};
