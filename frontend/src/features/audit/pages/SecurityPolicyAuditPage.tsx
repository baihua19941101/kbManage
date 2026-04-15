import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Space, Typography } from 'antd';
import {
  canReadPolicyAudit,
  useAuthStore
} from '@/features/auth/store';
import { AuditEventTable } from '@/features/audit/components/AuditEventTable';
import {
  listSecurityPolicyAuditEvents,
  type SecurityPolicyAuditFilters
} from '@/services/audit';

type FilterFormValues = {
  from?: string;
  to?: string;
  actorUserId?: string;
  action?: string;
  result?: string;
};

const normalizeValue = (value?: string) => {
  const trimmed = value?.trim();
  return trimmed && trimmed.length > 0 ? trimmed : undefined;
};

const mapFormToFilters = (values: FilterFormValues): SecurityPolicyAuditFilters => ({
  from: normalizeValue(values.from),
  to: normalizeValue(values.to),
  actorUserId: normalizeValue(values.actorUserId),
  action: normalizeValue(values.action),
  result: normalizeValue(values.result)
});

export const SecurityPolicyAuditPage = () => {
  const [form] = Form.useForm<FilterFormValues>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadPolicyAudit(user);
  const [filters, setFilters] = useState<SecurityPolicyAuditFilters>({});

  const queryKey = useMemo(() => ['audit-events', 'security-policy', filters], [filters]);

  const { data, isFetching, refetch, error } = useQuery({
    queryKey,
    enabled: canRead,
    queryFn: () => listSecurityPolicyAuditEvents(filters)
  });

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无策略审计访问权限，请联系管理员授予 securitypolicy:read 或审计角色。"
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          安全策略审计
        </Typography.Title>
        <Typography.Text type="secondary">
          按时间、操作者、动作和结果检索策略治理链路记录，支持违规处置复盘。
        </Typography.Text>
      </div>

      <Card>
        <Form<FilterFormValues> form={form} layout="vertical" onFinish={(values) => setFilters(mapFormToFilters(values))}>
          <Space wrap size="middle" align="start" style={{ width: '100%' }}>
            <Form.Item label="开始时间(ISO)" name="from">
              <Input placeholder="例如：2026-04-14T00:00:00Z" style={{ width: 240 }} />
            </Form.Item>
            <Form.Item label="结束时间(ISO)" name="to">
              <Input placeholder="例如：2026-04-14T23:59:59Z" style={{ width: 240 }} />
            </Form.Item>
            <Form.Item label="操作者" name="actorUserId">
              <Input placeholder="例如：admin" style={{ width: 180 }} />
            </Form.Item>
            <Form.Item label="动作" name="action">
              <Input placeholder="例如：securitypolicy.exception.review" style={{ width: 300 }} />
            </Form.Item>
            <Form.Item label="结果" name="result">
              <Input placeholder="例如：success / failed / denied" style={{ width: 220 }} />
            </Form.Item>
          </Space>

          <Space>
            <Button type="primary" htmlType="submit" loading={isFetching}>
              查询
            </Button>
            <Button
              onClick={() => {
                form.resetFields();
                setFilters({});
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
          message="策略审计查询失败"
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
