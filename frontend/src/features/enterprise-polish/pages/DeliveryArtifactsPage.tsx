import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Space } from 'antd';
import { DeliveryArtifactCatalog } from '@/features/enterprise-polish/components/DeliveryArtifactCatalog';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { listDeliveryArtifacts } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const DeliveryArtifactsPage = () => {
  const permissions = useEnterprisePermissions();
  const artifactsQuery = useQuery({ queryKey: ['enterprisePolish', 'artifacts'], queryFn: listDeliveryArtifacts });
  if (!permissions.canRead) return <PermissionDenied description="你暂无交付材料访问权限。" />;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="交付材料目录" description="查看安装、升级、运维和初始化配置模板目录。" />
      {artifactsQuery.error ? <Alert type="error" showIcon message={normalizeApiError(artifactsQuery.error, '交付材料加载失败')} /> : null}
      <Card><DeliveryArtifactCatalog items={artifactsQuery.data?.items ?? []} /></Card>
    </Space>
  );
};
