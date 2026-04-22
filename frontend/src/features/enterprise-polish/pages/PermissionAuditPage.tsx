import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Space } from 'antd';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { PermissionTrailTable } from '@/features/enterprise-polish/components/PermissionTrailTable';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { listPermissionTrails } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const PermissionAuditPage = () => {
  const permissions = useEnterprisePermissions();
  const trailsQuery = useQuery({ queryKey: ['enterprisePolish', 'permissionTrails'], queryFn: listPermissionTrails });
  if (!permissions.canRead) return <PermissionDenied description="你暂无深度审计访问权限。" />;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="权限变更审计" description="查看权限变更链路、授权依据和证据完整度。" />
      {trailsQuery.error ? <Alert type="error" showIcon message={normalizeApiError(trailsQuery.error, '权限审计加载失败')} /> : null}
      <Card><PermissionTrailTable items={trailsQuery.data?.items ?? []} /></Card>
    </Space>
  );
};
