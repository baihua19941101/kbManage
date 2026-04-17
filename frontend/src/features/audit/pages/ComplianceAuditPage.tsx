import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Space, Typography } from 'antd';
import {
  canReadComplianceAudit,
  useAuthStore
} from '@/features/auth/store';
import { AuditEventTable } from '@/features/audit/components/AuditEventTable';
import type { AuditEvent } from '@/services/audit';
import { listComplianceAuditEvents, type ComplianceAuditQuery } from '@/services/compliance';

type FilterFormValues = ComplianceAuditQuery;

const normalizeValue = (value?: string) => {
  const trimmed = value?.trim();
  return trimmed && trimmed.length > 0 ? trimmed : undefined;
};

const mapComplianceAuditEvent = (item: {
  action?: string;
  operatorId?: string;
  outcome?: string;
  occurredAt?: string;
  details?: Record<string, unknown>;
}): AuditEvent => {
  const details = item.details || {};
  return {
    id: `${item.occurredAt || 'unknown'}:${item.action || 'compliance'}`,
    eventType: item.action,
    action: item.action,
    actorUserId: item.operatorId,
    clusterId: typeof details.clusterId === 'string' ? details.clusterId : undefined,
    scopeId: typeof details.scopeId === 'string' ? details.scopeId : undefined,
    resourceKind: typeof details.resourceKind === 'string' ? details.resourceKind : undefined,
    resourceNamespace:
      typeof details.resourceNamespace === 'string' ? details.resourceNamespace : undefined,
    resourceName: typeof details.resourceName === 'string' ? details.resourceName : undefined,
    result: item.outcome,
    outcome: item.outcome,
    occurredAt: item.occurredAt
  };
};

export const ComplianceAuditPage = () => {
  const [form] = Form.useForm<FilterFormValues>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadComplianceAudit(user);
  const [filters, setFilters] = useState<ComplianceAuditQuery>({});

  const queryKey = useMemo(() => ['audit-events', 'compliance', filters], [filters]);

  const { data, isFetching, refetch, error } = useQuery({
    queryKey,
    enabled: canRead,
    queryFn: () => listComplianceAuditEvents(filters)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无合规审计访问权限。" />;
  }

  const tableData: AuditEvent[] = (data?.items || []).map(mapComplianceAuditEvent);

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          合规审计
        </Typography.Title>
        <Typography.Text type="secondary">
          按时间、动作、结果和基线检索合规治理链路记录，用于审计复盘与汇报。
        </Typography.Text>
      </div>

      <Card>
        <Form<FilterFormValues>
          form={form}
          layout="vertical"
          onFinish={(values) =>
            setFilters({
              timeFrom: normalizeValue(values.timeFrom),
              timeTo: normalizeValue(values.timeTo),
              baselineId: normalizeValue(values.baselineId),
              action: normalizeValue(values.action),
              outcome: normalizeValue(values.outcome)
            })
          }
        >
          <Space wrap size="middle" align="start" style={{ width: '100%' }}>
            <Form.Item label="开始时间(ISO)" name="timeFrom">
              <Input placeholder="例如：2026-04-14T00:00:00Z" style={{ width: 240 }} />
            </Form.Item>
            <Form.Item label="结束时间(ISO)" name="timeTo">
              <Input placeholder="例如：2026-04-14T23:59:59Z" style={{ width: 240 }} />
            </Form.Item>
            <Form.Item label="基线 ID" name="baselineId">
              <Input placeholder="例如：baseline-001" style={{ width: 220 }} />
            </Form.Item>
            <Form.Item label="动作" name="action">
              <Input placeholder="例如：compliance.exception.review" style={{ width: 280 }} />
            </Form.Item>
            <Form.Item label="结果" name="outcome">
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
          message="合规审计查询失败"
          description={error instanceof Error ? error.message : '请求失败，请稍后重试。'}
        />
      ) : null}

      <AuditEventTable
        data={tableData}
        loading={isFetching}
        onRefresh={() => {
          void refetch();
        }}
      />
    </Space>
  );
};
