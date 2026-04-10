import { useMemo, useState } from 'react';
import type { Dayjs } from 'dayjs';
import { Button, Card, DatePicker, Form, Input, Select, Space, Typography, message } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { AuditEventTable } from '@/features/audit/components/AuditEventTable';
import { AuditExportModal } from '@/features/audit/components/AuditExportModal';
import { listAuditEvents, type AuditEventFilters } from '@/services/audit';

const { RangePicker } = DatePicker;

type FilterFormValues = {
  timeRange?: [Dayjs, Dayjs];
  actorUserId?: string;
  clusterId?: string;
  eventType?: string;
  result?: string;
};

const normalizeValue = (value?: string) => {
  const trimmed = value?.trim();
  return trimmed && trimmed.length > 0 ? trimmed : undefined;
};

const mapFormToFilters = (values: FilterFormValues): AuditEventFilters => ({
  from: values.timeRange?.[0]?.toISOString(),
  to: values.timeRange?.[1]?.toISOString(),
  actorUserId: normalizeValue(values.actorUserId),
  clusterId: normalizeValue(values.clusterId),
  eventType: normalizeValue(values.eventType),
  result: normalizeValue(values.result)
});

export const AuditEventPage = () => {
  const [form] = Form.useForm<FilterFormValues>();
  const [filters, setFilters] = useState<AuditEventFilters>({});
  const [exportOpen, setExportOpen] = useState(false);

  const queryKey = useMemo(() => ['audit-events', filters], [filters]);

  const { data, isFetching, refetch } = useQuery({
    queryKey,
    queryFn: () => listAuditEvents(filters)
  });

  const onSearch = (values: FilterFormValues) => {
    setFilters(mapFormToFilters(values));
  };

  const onReset = () => {
    form.resetFields();
    setFilters({});
  };

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          审计事件查询
        </Typography.Title>
        <Typography.Text type="secondary">
          支持按时间范围、操作者、事件类型与执行结果筛选，并可提交导出任务。
        </Typography.Text>
      </div>

      <Card>
        <Form<FilterFormValues> form={form} layout="vertical" onFinish={onSearch}>
          <Space wrap size="middle" align="start" style={{ width: '100%' }}>
            <Form.Item label="时间范围" name="timeRange">
              <RangePicker showTime allowClear style={{ width: 340 }} />
            </Form.Item>
            <Form.Item label="操作者" name="actorUserId">
              <Input placeholder="例如：admin" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="集群" name="clusterId">
              <Input placeholder="例如：prod-cn" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="事件类型" name="eventType">
              <Input placeholder="例如：operation.restart" style={{ width: 220 }} />
            </Form.Item>
            <Form.Item label="结果" name="result">
              <Select
                allowClear
                placeholder="全部"
                style={{ width: 150 }}
                options={[
                  { value: 'success', label: 'success' },
                  { value: 'failed', label: 'failed' },
                  { value: 'denied', label: 'denied' },
                  { value: 'pending', label: 'pending' }
                ]}
              />
            </Form.Item>
          </Space>

          <Space>
            <Button type="primary" htmlType="submit" loading={isFetching}>
              查询
            </Button>
            <Button onClick={onReset}>重置</Button>
            <Button
              onClick={() => {
                if (!filters.from || !filters.to) {
                  message.warning('导出前请先设置时间范围');
                  return;
                }
                setExportOpen(true);
              }}
            >
              导出
            </Button>
          </Space>
        </Form>
      </Card>

      <AuditEventTable
        data={data?.items || []}
        loading={isFetching}
        onRefresh={() => {
          void refetch();
        }}
      />

      <AuditExportModal
        open={exportOpen}
        filters={filters}
        onCancel={() => setExportOpen(false)}
        onSubmitted={() => {
          void refetch();
        }}
      />
    </Space>
  );
};
