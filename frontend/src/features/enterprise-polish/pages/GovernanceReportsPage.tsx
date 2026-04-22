import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space } from 'antd';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { ReportBuilderDrawer } from '@/features/enterprise-polish/components/ReportBuilderDrawer';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { useReportActions } from '@/features/enterprise-polish/hooks/useReportActions';
import { listGovernanceReports } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const GovernanceReportsPage = () => {
  const permissions = useEnterprisePermissions();
  const reportsQuery = useQuery({ queryKey: ['enterprisePolish', 'reports'], queryFn: listGovernanceReports });
  const { createReport } = useReportActions();
  if (!permissions.canRead) return <PermissionDenied description="你暂无治理报表访问权限。" />;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="治理报表中心" description="生成管理汇报、审计复核和客户交付三类报表。" actions={permissions.canManageReports ? <ReportBuilderDrawer onSubmit={async (payload) => { await createReport.mutateAsync(payload); }} /> : null} />
      {reportsQuery.error ? <Alert type="error" showIcon message={normalizeApiError(reportsQuery.error, '治理报表加载失败')} /> : null}
      <Card>
        <List dataSource={reportsQuery.data?.items ?? []} renderItem={(item) => <List.Item>{item.title} / {item.reportType || '未知类型'} / {item.audienceType || '未知对象'}</List.Item>} />
      </Card>
    </Space>
  );
};
