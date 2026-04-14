import { useMemo, useState } from 'react';
import type { Dayjs } from 'dayjs';
import { Alert, Button, Card, DatePicker, Form, Input, Select, Space, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { canReadGitOps, useAuthStore } from '@/features/auth/store';
import { AuditEventTable } from '@/features/audit/components/AuditEventTable';
import { listAuditEvents, type AuditEventFilters } from '@/services/audit';

const { RangePicker } = DatePicker;

type FilterFormValues = {
  timeRange?: [Dayjs, Dayjs];
  actorUserId?: string;
  workspaceId?: string;
  projectId?: string;
  action?: string;
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
  workspaceId: normalizeValue(values.workspaceId),
  projectId: normalizeValue(values.projectId),
  eventType: normalizeValue(values.action),
  result: normalizeValue(values.result),
  resource: 'gitops',
  actionPrefix: 'gitops.'
});

export const GitOpsAuditPage = () => {
  const [form] = Form.useForm<FilterFormValues>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadGitOps(user);
  const [filters, setFilters] = useState<AuditEventFilters>({
    resource: 'gitops',
    actionPrefix: 'gitops.'
  });

  const queryKey = useMemo(() => ['audit-events', 'gitops', filters], [filters]);

  const { data, isFetching, refetch, error } = useQuery({
    queryKey,
    enabled: canRead,
    queryFn: () => listAuditEvents(filters)
  });

  if (!canRead) {
    return (
      <Alert
        type="warning"
        showIcon
        message="你暂无 GitOps 审计访问权限"
        description="请联系管理员授予 GitOps 读取范围后再访问。"
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          GitOps 审计查询
        </Typography.Title>
        <Typography.Text type="secondary">检索 GitOps 来源校验、同步、推进、回滚与发布动作审计记录。</Typography.Text>
      </div>

      <Card>
        <Form<FilterFormValues> form={form} layout="vertical" onFinish={(values) => setFilters(mapFormToFilters(values))}>
          <Space wrap size="middle" align="start" style={{ width: '100%' }}>
            <Form.Item label="时间范围" name="timeRange">
              <RangePicker showTime allowClear style={{ width: 340 }} />
            </Form.Item>
            <Form.Item label="操作者" name="actorUserId">
              <Input placeholder="例如：admin" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="工作空间 ID" name="workspaceId">
              <Input placeholder="例如：1001" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="项目 ID" name="projectId">
              <Input placeholder="例如：2001" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="动作（可选）" name="action">
              <Input placeholder="例如：gitops.rollback.submit" style={{ width: 240 }} />
            </Form.Item>
            <Form.Item label="结果" name="result">
              <Select
                allowClear
                placeholder="全部"
                style={{ width: 150 }}
                options={[
                  { value: 'success', label: 'success' },
                  { value: 'failed', label: 'failed' },
                  { value: 'denied', label: 'denied' }
                ]}
              />
            </Form.Item>
          </Space>

          <Space>
            <Button type="primary" htmlType="submit" loading={isFetching}>
              查询
            </Button>
            <Button
              onClick={() => {
                form.resetFields();
                setFilters({ resource: 'gitops', actionPrefix: 'gitops.' });
              }}
            >
              重置
            </Button>
          </Space>
        </Form>
      </Card>

      {error ? (
        <Alert
          type="error"
          showIcon
          message="GitOps 审计查询失败"
          description={error instanceof Error ? error.message : '请求失败，请稍后重试。'}
        />
      ) : null}

      <AuditEventTable
        data={data?.items || []}
        loading={isFetching}
        onRefresh={() => {
          void refetch();
        }}
      />
    </Space>
  );
};
