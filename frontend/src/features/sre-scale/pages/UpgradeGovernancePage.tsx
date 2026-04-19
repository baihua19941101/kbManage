import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { StatusTag } from '@/features/sre-scale/components/status';
import { UpgradePlanDrawer } from '@/features/sre-scale/components/UpgradePlanDrawer';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { useUpgradeAction } from '@/features/sre-scale/hooks/useUpgradeAction';
import { normalizeApiError } from '@/services/api/client';
import { listUpgradePlans, type SREUpgradePlan } from '@/services/sreScale';

const columns: ColumnsType<SREUpgradePlan> = [
  { title: '计划', dataIndex: 'name', key: 'name' },
  { title: '当前版本', dataIndex: 'currentVersion', key: 'currentVersion' },
  { title: '目标版本', dataIndex: 'targetVersion', key: 'targetVersion' },
  { title: '阶段', dataIndex: 'executionStage', key: 'executionStage' },
  { title: '进度', dataIndex: 'executionProgress', key: 'executionProgress' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const UpgradeGovernancePage = () => {
  const permissions = useSREPermissions();
  const action = useUpgradeAction();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const plansQuery = useQuery({ queryKey: ['sreScale', 'upgrades'], queryFn: () => listUpgradePlans() });
  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无升级治理访问权限。" />;
  }
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="升级治理"
        description="执行升级前检查、创建升级计划并跟踪滚动阶段。"
        actions={
          <Button type="primary" disabled={!permissions.canManageUpgrade} onClick={() => setDrawerOpen(true)}>
            创建升级计划
          </Button>
        }
      />
      {plansQuery.error ? <Alert type="error" showIcon message={normalizeApiError(plansQuery.error, '升级计划加载失败')} /> : null}
      <Card size="small" title="升级计划">
        <Table rowKey="id" columns={columns} dataSource={plansQuery.data?.items || []} pagination={{ pageSize: 6 }} />
      </Card>
      <UpgradePlanDrawer
        open={drawerOpen}
        submitting={action.createUpgradeMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          action.createUpgradeMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
