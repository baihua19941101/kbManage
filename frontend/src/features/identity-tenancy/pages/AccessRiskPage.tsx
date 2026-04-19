import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Space, Statistic, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  accessRiskQueryScope,
  identityTenancyQueryKeys,
  listAccessRisks,
  type AccessRiskSnapshot
} from '@/services/identityTenancy';

const columns: ColumnsType<AccessRiskSnapshot> = [
  { title: '主体', key: 'subject', render: (_, record) => `${record.subjectType || '—'} / ${record.subjectRef || '—'}` },
  { title: '风险类型', dataIndex: 'riskType', key: 'riskType', render: (value?: string) => value || '—' },
  { title: '严重等级', dataIndex: 'severity', key: 'severity', render: (value?: string) => <StatusTag value={value} /> },
  { title: '摘要', dataIndex: 'summary', key: 'summary', render: (value?: string) => value || '—' },
  { title: '建议动作', dataIndex: 'recommendedAction', key: 'recommendedAction', render: (value?: string) => value || '—' }
];

export const AccessRiskPage = () => {
  const permissions = useIdentityTenancyPermissions();
  const risksQuery = useQuery({
    queryKey: identityTenancyQueryKeys.risks(accessRiskQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listAccessRisks({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无访问风险视图访问权限。" />;
  }

  const items = risksQuery.data?.items || [];
  const criticalCount = items.filter((item) => item.severity === 'critical').length;
  const highCount = items.filter((item) => item.severity === 'high').length;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="访问风险视图"
        description="统一查看越权授权、边界扩散、到期未回收和高风险会话带来的访问风险。"
      />

      {risksQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="访问风险加载失败"
          description={normalizeApiError(risksQuery.error, '访问风险加载失败，请稍后重试。')}
        />
      ) : null}

      <Space size={16} wrap>
        <Card size="small">
          <Statistic title="严重风险" value={criticalCount} />
        </Card>
        <Card size="small">
          <Statistic title="高风险" value={highCount} />
        </Card>
        <Card size="small">
          <Statistic title="总风险数" value={items.length} />
        </Card>
      </Space>

      <Card size="small" title="风险摘要">
        <Table<AccessRiskSnapshot>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={items}
          loading={risksQuery.isLoading || risksQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>
    </Space>
  );
};
