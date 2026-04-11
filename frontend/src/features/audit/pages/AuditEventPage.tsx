import { useMemo, useState } from 'react';
import type { Dayjs } from 'dayjs';
import { Alert, Button, Card, DatePicker, Form, Input, Select, Space, Steps, Tag, Typography, message } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { AuditEventTable } from '@/features/audit/components/AuditEventTable';
import { AuditExportModal } from '@/features/audit/components/AuditExportModal';
import {
  getAuditExportTask,
  listAuditEvents,
  type AuditEventFilters
} from '@/services/audit';

const { RangePicker } = DatePicker;

type FilterFormValues = {
  timeRange?: [Dayjs, Dayjs];
  actorUserId?: string;
  clusterId?: string;
  eventType?: string;
  result?: string;
};

type ExportProgressStatus = 'pending' | 'running' | 'succeeded' | 'failed';

const exportStepMap: Record<ExportProgressStatus, { current: number; stepStatus?: 'error' | 'process' | 'finish' | 'wait' }> = {
  pending: { current: 0, stepStatus: 'process' },
  running: { current: 1, stepStatus: 'process' },
  succeeded: { current: 2, stepStatus: 'finish' },
  failed: { current: 2, stepStatus: 'error' }
};

const statusMeta: Record<ExportProgressStatus, { color: string; label: string }> = {
  pending: { color: 'default', label: '排队中' },
  running: { color: 'processing', label: '处理中' },
  succeeded: { color: 'success', label: '已完成' },
  failed: { color: 'error', label: '失败' }
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

const toErrorText = (error: unknown): string => {
  if (error instanceof Error && error.message.trim().length > 0) {
    return error.message;
  }
  return '请求失败，请稍后重试。';
};

export const AuditEventPage = () => {
  const [form] = Form.useForm<FilterFormValues>();
  const [filters, setFilters] = useState<AuditEventFilters>({});
  const [exportOpen, setExportOpen] = useState(false);
  const [exportTaskId, setExportTaskId] = useState<string>('');

  const queryKey = useMemo(() => ['audit-events', filters], [filters]);

  const { data, isFetching, refetch } = useQuery({
    queryKey,
    queryFn: () => listAuditEvents(filters)
  });

  const exportStatusQuery = useQuery({
    queryKey: ['audit-export-status', exportTaskId],
    queryFn: () => getAuditExportTask(exportTaskId),
    enabled: exportTaskId.trim().length > 0,
    refetchInterval: (query) => {
      const status = query.state.data?.status;
      return status === 'pending' || status === 'running' ? 2000 : false;
    }
  });

  const exportStatus = exportStatusQuery.data?.status;
  const exportStatusDisplay = exportStatus ? statusMeta[exportStatus] : undefined;
  const isExportPolling = exportStatus === 'pending' || exportStatus === 'running';

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

      {exportTaskId ? (
        <Card
          title={`导出任务：${exportTaskId}`}
          extra={
            <Space>
              {exportStatusDisplay ? (
                <Tag color={exportStatusDisplay.color}>{exportStatusDisplay.label}</Tag>
              ) : null}
              {exportStatusQuery.data?.downloadUrl ? (
                <Button
                  size="small"
                  type="primary"
                  onClick={() => {
                    window.open(exportStatusQuery.data?.downloadUrl, '_blank', 'noopener,noreferrer');
                  }}
                >
                  下载文件
                </Button>
              ) : null}
              <Button
                size="small"
                loading={exportStatusQuery.isFetching}
                onClick={() => {
                  void exportStatusQuery.refetch();
                }}
              >
                刷新状态
              </Button>
            </Space>
          }
        >
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Steps
              size="small"
              current={exportStatus ? exportStepMap[exportStatus].current : 0}
              status={exportStatus ? exportStepMap[exportStatus].stepStatus : 'process'}
              items={[
                { title: '已提交' },
                { title: '处理中' },
                { title: '完成' }
              ]}
            />

            {exportStatusQuery.isError ? (
              <Alert
                type="error"
                showIcon
                message="导出状态查询失败"
                description={toErrorText(exportStatusQuery.error)}
                action={
                  <Button
                    size="small"
                    onClick={() => {
                      void exportStatusQuery.refetch();
                    }}
                  >
                    重试查询
                  </Button>
                }
              />
            ) : null}

            {exportStatusQuery.data?.status === 'failed' ? (
              <Alert
                type="error"
                showIcon
                message="导出失败"
                description={exportStatusQuery.data.errorMessage || '后端未返回具体失败原因。'}
                action={
                  <Button
                    size="small"
                    onClick={() => {
                      setExportOpen(true);
                    }}
                  >
                    重新导出
                  </Button>
                }
              />
            ) : null}

            {exportStatusQuery.data?.status === 'succeeded' ? (
              <Alert
                type="success"
                showIcon
                message="导出已完成"
                description={
                  exportStatusQuery.data.resultTotal !== undefined
                    ? `共 ${exportStatusQuery.data.resultTotal} 条结果`
                    : '导出结果已生成'
                }
              />
            ) : null}

            {!exportStatusQuery.isError &&
            !exportStatusQuery.data &&
            (exportStatusQuery.isFetching || isExportPolling) ? (
              <Typography.Text type="secondary">导出任务已提交，正在轮询任务状态...</Typography.Text>
            ) : null}
          </Space>
        </Card>
      ) : null}

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
        onSubmitted={(taskId) => {
          setExportTaskId(taskId);
          void refetch();
        }}
      />
    </Space>
  );
};
