import { useQuery } from '@tanstack/react-query';
import { Alert, Space } from 'antd';
import { CapacityTrendChart } from '@/features/sre-scale/components/CapacityTrendChart';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { SelfDiagnosisCard } from '@/features/sre-scale/components/SelfDiagnosisCard';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import { listCapacityBaselines, listScaleEvidence } from '@/services/sreScale';

export const CapacityGovernancePage = () => {
  const permissions = useSREPermissions();
  const baselineQuery = useQuery({ queryKey: ['sreScale', 'capacity'], queryFn: () => listCapacityBaselines() });
  const evidenceQuery = useQuery({ queryKey: ['sreScale', 'scaleEvidence'], queryFn: () => listScaleEvidence({}) });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无容量性能治理访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="容量与性能治理" description="查看容量基线、趋势、瓶颈摘要和可信度说明。" />
      {baselineQuery.error ? <Alert type="error" showIcon message={normalizeApiError(baselineQuery.error, '容量基线加载失败')} /> : null}
      {evidenceQuery.error ? <Alert type="error" showIcon message={normalizeApiError(evidenceQuery.error, '规模化证据加载失败')} /> : null}
      <CapacityTrendChart baselines={baselineQuery.data?.items || []} evidence={evidenceQuery.data?.items || []} />
      <SelfDiagnosisCard evidence={evidenceQuery.data?.items?.[0]} />
    </Space>
  );
};
