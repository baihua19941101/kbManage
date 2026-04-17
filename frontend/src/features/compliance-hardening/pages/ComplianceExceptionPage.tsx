import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canReadCompliance,
  canReviewComplianceException,
  useAuthStore
} from '@/features/auth/store';
import { ComplianceExceptionReviewDrawer } from '@/features/compliance-hardening/components/ComplianceExceptionReviewDrawer';
import { useComplianceAction } from '@/features/compliance-hardening/hooks/useComplianceAction';
import { exceptionStatusColorMap, formatDateTime } from '@/features/compliance-hardening/utils';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import {
  listComplianceExceptions,
  type ComplianceExceptionRequest,
  type ComplianceExceptionStatus
} from '@/services/compliance';

const statusOptions: Array<{ label: string; value: ComplianceExceptionStatus | '' }> = [
  { label: '全部状态', value: '' },
  { label: '待审批', value: 'pending' },
  { label: '批准', value: 'approved' },
  { label: '拒绝', value: 'rejected' },
  { label: '生效中', value: 'active' },
  { label: '已到期', value: 'expired' },
  { label: '已撤销', value: 'revoked' }
];

const columns = (
  readonly: boolean,
  onReview: (exception: ComplianceExceptionRequest) => void
): ColumnsType<ComplianceExceptionRequest> => [
  {
    title: '申请单',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '原因',
    dataIndex: 'reason',
    key: 'reason',
    render: (value?: string) => value || '—'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: ComplianceExceptionRequest['status']) =>
      value ? <Tag color={exceptionStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '生效窗口',
    key: 'window',
    render: (_, record) => `${formatDateTime(record.startsAt)} ~ ${formatDateTime(record.expiresAt)}`
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" disabled={readonly} onClick={() => onReview(record)}>
        审批
      </Button>
    )
  }
];

export const ComplianceExceptionPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const canReview = canReviewComplianceException(user);
  const [status, setStatus] = useState<ComplianceExceptionStatus | ''>('');
  const [selectedException, setSelectedException] = useState<ComplianceExceptionRequest>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { reviewExceptionMutation } = useComplianceAction();

  const queryInput = useMemo(() => ({ status: status || undefined }), [status]);
  const exceptionsQuery = useQuery({
    queryKey: ['compliance', 'exceptions', queryInput],
    enabled: canRead,
    queryFn: () => listComplianceExceptions(queryInput)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无例外审批访问权限。" />;
  }

  const exceptions = exceptionsQuery.data?.items || [];
  const permissionChanged = isAuthorizationError(exceptionsQuery.error);

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          合规例外审批
        </Typography.Title>
        <Typography.Text type="secondary">
          审批失败项例外申请，查看审批意见和生效时间窗口。
        </Typography.Text>
      </div>

      {!canReview ? (
        <Alert type="info" showIcon message="当前为只读模式" description="你可以查看例外申请，但无法审批。" />
      ) : null}

      {permissionChanged ? (
        <Alert type="warning" showIcon message="权限已变更" description="当前会话已失去例外审批权限，请刷新或重新登录。" />
      ) : null}

      {exceptionsQuery.error && !permissionChanged ? (
        <Alert type="error" showIcon message="例外申请加载失败" description={normalizeApiError(exceptionsQuery.error, '例外申请加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选">
        <Select style={{ width: 180 }} options={statusOptions} value={status} onChange={setStatus} />
      </Card>

      <Card size="small" title={`例外申请（${exceptions.length}）`}>
        <Table<ComplianceExceptionRequest>
          rowKey="id"
          dataSource={exceptions}
          loading={exceptionsQuery.isLoading || exceptionsQuery.isFetching}
          columns={columns(!canReview || permissionChanged, (item) => {
            setSelectedException(item);
            setDrawerOpen(true);
          })}
          pagination={{ pageSize: 8 }}
          scroll={{ x: 900 }}
        />
      </Card>

      <ComplianceExceptionReviewDrawer
        open={drawerOpen}
        exception={selectedException}
        readonly={!canReview || permissionChanged}
        loading={reviewExceptionMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(exceptionId, payload) => reviewExceptionMutation.mutate({ exceptionId, payload })}
      />
    </Space>
  );
};
