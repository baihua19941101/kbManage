import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space } from 'antd';
import { HealthStatusCard } from '@/features/sre-scale/components/HealthStatusCard';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import { getHealthOverview } from '@/services/sreScale';

export const HealthOverviewPage = () => {
  const permissions = useSREPermissions();
  const overviewQuery = useQuery({ queryKey: ['sreScale', 'health'], queryFn: () => getHealthOverview({ workspaceId: 1 }) });
  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无平台健康总览访问权限。" />;
  }
  const overview = overviewQuery.data;
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="平台健康总览" description="集中查看控制面组件健康、依赖状态、任务积压和恢复摘要。" />
      {overviewQuery.error ? <Alert type="error" showIcon message={normalizeApiError(overviewQuery.error, '平台健康总览加载失败')} /> : null}
      <HealthStatusCard overview={overview} />
      <Card size="small" title="建议动作">
        <List dataSource={overview?.recommendedActions || []} renderItem={(item) => <List.Item>{item}</List.Item>} />
      </Card>
    </Space>
  );
};
