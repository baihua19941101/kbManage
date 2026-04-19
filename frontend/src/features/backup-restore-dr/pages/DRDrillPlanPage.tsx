import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useState } from 'react';
import { DRDrillPlanDrawer } from '@/features/backup-restore-dr/components/DRDrillPlanDrawer';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { useDrillAction } from '@/features/backup-restore-dr/hooks/useDrillAction';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  backupRestoreQueryKeys,
  listDrillPlans,
  type DRDrillPlan
} from '@/services/backupRestore';

const columns = (
  onRun: (planId: string) => void,
  canDrill: boolean
): ColumnsType<DRDrillPlan> => [
  { title: '计划名称', dataIndex: 'name', key: 'name' },
  {
    title: '目标',
    key: 'targets',
    render: (_, record) => `RPO ${record.rpoTargetMinutes || '—'} 分钟 / RTO ${record.rtoTargetMinutes || '—'} 分钟`
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <StatusTag value={value} />
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button size="small" disabled={!canDrill} onClick={() => onRun(record.id)}>
        发起演练
      </Button>
    )
  }
];

export const DRDrillPlanPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useBackupRestorePermissions();
  const { createPlanMutation, runDrillMutation } = useDrillAction();
  const plansQuery = useQuery({
    queryKey: backupRestoreQueryKeys.drillPlans(),
    enabled: permissions.canRead,
    queryFn: () => listDrillPlans()
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无灾备演练计划访问权限。" />;
  }

  const plans = plansQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="灾备演练计划"
        description="维护演练范围、角色分工、RPO/RTO 目标和切换清单。"
        actions={
          <Button type="primary" disabled={!permissions.canDrill} onClick={() => setDrawerOpen(true)}>
            新建演练计划
          </Button>
        }
      />

      {plansQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="演练计划加载失败"
          description={normalizeApiError(plansQuery.error, '演练计划加载失败，请稍后重试。')}
        />
      ) : null}
      {createPlanMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="演练计划创建失败"
          description={normalizeApiError(createPlanMutation.error, '演练计划创建失败，请稍后重试。')}
        />
      ) : null}
      {runDrillMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="演练发起失败"
          description={normalizeApiError(runDrillMutation.error, '演练发起失败，请稍后重试。')}
        />
      ) : null}
      {!permissions.canDrill ? (
        <Alert type="info" showIcon message="你当前只有查看权限，演练动作已被禁用。" />
      ) : null}

      <Card size="small" title={`演练计划（${plans.length}）`}>
        <Table<DRDrillPlan>
          rowKey={(record) => record.id}
          dataSource={plans}
          columns={columns((planId) => runDrillMutation.mutate(planId), permissions.canDrill)}
          loading={plansQuery.isLoading || plansQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <DRDrillPlanDrawer
        open={drawerOpen}
        submitting={createPlanMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createPlanMutation.mutate(payload, {
            onSuccess: () => {
              setDrawerOpen(false);
            }
          })
        }
      />
    </Space>
  );
};
