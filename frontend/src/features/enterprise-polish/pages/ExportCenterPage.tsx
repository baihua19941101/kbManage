import { useQueries } from '@tanstack/react-query';
import { Alert, Card, Space } from 'antd';
import { ActionItemList } from '@/features/enterprise-polish/components/ActionItemList';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { listGovernanceActionItems, listGovernanceCoverage } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const ExportCenterPage = () => {
  const permissions = useEnterprisePermissions();
  const [coverageQuery, actionQuery] = useQueries({
    queries: [
      { queryKey: ['enterprisePolish', 'coverage'], queryFn: listGovernanceCoverage },
      { queryKey: ['enterprisePolish', 'actionItems'], queryFn: listGovernanceActionItems }
    ]
  });
  if (!permissions.canRead) return <PermissionDenied description="你暂无导出中心访问权限。" />;
  const error = coverageQuery.error || actionQuery.error;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="导出与治理待办" description="查看治理覆盖率和统一待办，跟踪报表导出背景数据。" />
      {error ? <Alert type="error" showIcon message={normalizeApiError(error, '导出中心加载失败')} /> : null}
      <Card title="统一治理待办"><ActionItemList items={actionQuery.data?.items ?? []} /></Card>
    </Space>
  );
};
