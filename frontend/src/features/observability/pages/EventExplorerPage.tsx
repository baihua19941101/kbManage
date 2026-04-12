import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Form, Input, InputNumber, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { EventTimeline } from '@/features/observability/components/EventTimeline';
import { listObservabilityEvents, type EventQueryParams } from '@/services/observability/events';

const toNumber = (value: string | null, fallback: number) => {
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
};

export const EventExplorerPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const initialFilters = useMemo<EventQueryParams>(
    () => ({
      clusterId: searchParams.get('clusterId') ?? undefined,
      namespace: searchParams.get('namespace') ?? undefined,
      resourceKind: searchParams.get('resourceKind') ?? undefined,
      resourceName: searchParams.get('resourceName') ?? undefined,
      eventType: searchParams.get('eventType') ?? undefined,
      startAt: searchParams.get('startAt') ?? undefined,
      endAt: searchParams.get('endAt') ?? undefined,
      limit: toNumber(searchParams.get('limit'), 100)
    }),
    [searchParams]
  );
  const [filters, setFilters] = useState<EventQueryParams>(initialFilters);

  const eventsQuery = useQuery({
    queryKey: ['observability', 'events', filters],
    queryFn: () => listObservabilityEvents(filters)
  });

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ marginBottom: 0 }}>
          事件时间线
        </Typography.Title>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
          <Button onClick={() => void navigate(`/observability/logs${location.search}`)}>日志</Button>
          <Button onClick={() => void navigate(`/observability/metrics${location.search}`)}>指标</Button>
        </Space>
      </Space>

      <Card>
        <Form<EventQueryParams>
          layout="vertical"
          initialValues={initialFilters}
          onFinish={(values) => setFilters(values)}
        >
          <Space wrap style={{ width: '100%' }} align="start">
            <Form.Item name="clusterId" label="Cluster ID" style={{ width: 180 }}>
              <Input placeholder="cluster-1" />
            </Form.Item>
            <Form.Item name="namespace" label="Namespace" style={{ width: 180 }}>
              <Input placeholder="default" />
            </Form.Item>
            <Form.Item name="resourceKind" label="Kind" style={{ width: 180 }}>
              <Input placeholder="Deployment" />
            </Form.Item>
            <Form.Item name="resourceName" label="Resource" style={{ width: 180 }}>
              <Input placeholder="mock-app" />
            </Form.Item>
          </Space>
          <Space wrap style={{ width: '100%' }} align="start">
            <Form.Item name="eventType" label="级别" style={{ width: 180 }}>
              <Input placeholder="warning/normal" />
            </Form.Item>
            <Form.Item name="startAt" label="开始时间" style={{ width: 220 }}>
              <Input placeholder="2026-04-11T00:00:00Z" />
            </Form.Item>
            <Form.Item name="endAt" label="结束时间" style={{ width: 220 }}>
              <Input placeholder="2026-04-11T01:00:00Z" />
            </Form.Item>
            <Form.Item name="limit" label="条数" style={{ width: 120 }}>
              <InputNumber min={1} max={500} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item label=" ">
              <Button type="primary" htmlType="submit" loading={eventsQuery.isFetching}>
                查询
              </Button>
            </Form.Item>
          </Space>
        </Form>
      </Card>

      <Card>
        {eventsQuery.error ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="事件查询失败"
            description={`${eventsQuery.error}`}
          />
        ) : null}
        <EventTimeline items={eventsQuery.data?.items ?? []} />
      </Card>
    </Space>
  );
};
