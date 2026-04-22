import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Space } from 'antd';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { RiskTrendChart } from '@/features/enterprise-polish/components/RiskTrendChart';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { listKeyOperations } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const RiskTrackingPage = () => {
  const permissions = useEnterprisePermissions();
  const operationsQuery = useQuery({ queryKey: ['enterprisePolish', 'keyOperations'], queryFn: listKeyOperations });
  if (!permissions.canRead) return <PermissionDenied description="你暂无风险追踪访问权限。" />;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="高风险访问追踪" description="查看关键操作、高风险访问与长期异常趋势。" />
      {operationsQuery.error ? <Alert type="error" showIcon message={normalizeApiError(operationsQuery.error, '风险追踪加载失败')} /> : null}
      <Card><RiskTrendChart items={operationsQuery.data?.items ?? []} /></Card>
    </Space>
  );
};
