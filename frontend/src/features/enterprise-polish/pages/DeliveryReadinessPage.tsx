import { useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Space } from 'antd';
import { DeliveryChecklistBoard } from '@/features/enterprise-polish/components/DeliveryChecklistBoard';
import { DeliveryScopeNotice } from '@/features/enterprise-polish/components/DeliveryScopeNotice';
import { PageHeader } from '@/features/enterprise-polish/components/PageHeader';
import { PermissionDenied } from '@/features/enterprise-polish/components/PermissionDenied';
import { ReadinessSummaryCard } from '@/features/enterprise-polish/components/ReadinessSummaryCard';
import { useDeliveryBundleChecklist } from '@/features/enterprise-polish/hooks/useDeliveryBundleActions';
import { useEnterprisePermissions } from '@/features/enterprise-polish/hooks/permissions';
import { listDeliveryBundles } from '@/services/enterprisePolish';
import { normalizeApiError } from '@/services/api/client';

export const DeliveryReadinessPage = () => {
  const permissions = useEnterprisePermissions();
  const bundlesQuery = useQuery({ queryKey: ['enterprisePolish', 'bundles'], queryFn: listDeliveryBundles });
  const firstBundleId = bundlesQuery.data?.items?.[0]?.id;
  const checklistQuery = useDeliveryBundleChecklist(firstBundleId);

  useEffect(() => {
    // keep react compiler-friendly side effect slot for future refresh logic
  }, [firstBundleId]);

  if (!permissions.canRead) return <PermissionDenied description="你暂无交付就绪访问权限。" />;
  const error = bundlesQuery.error || checklistQuery.error;
  const bundle = bundlesQuery.data?.items?.[0];
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="交付就绪检查" description="查看交付包、缺失项和交付清单完成状态。" />
      {error ? <Alert type="error" showIcon message={normalizeApiError(error, '交付就绪加载失败')} /> : null}
      <DeliveryScopeNotice text="交付材料需结合适用版本、环境和客户范围使用。" />
      <ReadinessSummaryCard summary={bundle?.artifactSummary} conclusion={bundle?.readinessConclusion} />
      <Card title="交付检查清单"><DeliveryChecklistBoard items={checklistQuery.data?.items ?? []} /></Card>
    </Space>
  );
};
