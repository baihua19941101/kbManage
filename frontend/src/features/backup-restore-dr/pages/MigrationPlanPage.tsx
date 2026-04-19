import { useState } from 'react';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { useRestoreAction } from '@/features/backup-restore-dr/hooks/useRestoreAction';
import { MigrationPlanDrawer } from '@/features/backup-restore-dr/components/MigrationPlanDrawer';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { normalizeApiError } from '@/services/api/client';
import { backupRestoreScopeSummary, type MigrationPlan } from '@/services/backupRestore';

const columns: ColumnsType<MigrationPlan> = [
  { title: '计划名称', dataIndex: 'name', key: 'name' },
  { title: '源集群', dataIndex: 'sourceClusterId', key: 'sourceClusterId', render: (value?: string) => value || '—' },
  { title: '目标集群', dataIndex: 'targetClusterId', key: 'targetClusterId', render: (value?: string) => value || '—' },
  {
    title: '迁移范围',
    dataIndex: 'scopeSelection',
    key: 'scopeSelection',
    render: (value?: Record<string, unknown>) => backupRestoreScopeSummary(value)
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <StatusTag value={value} />
  }
];

export const MigrationPlanPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [plans, setPlans] = useState<MigrationPlan[]>([]);
  const permissions = useBackupRestorePermissions();
  const { migrationMutation } = useRestoreAction();

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无迁移计划访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="迁移计划中心"
        description="定义跨集群迁移范围、映射规则和切换步骤。"
        actions={
          <Button type="primary" disabled={!permissions.canMigrate} onClick={() => setDrawerOpen(true)}>
            新建迁移计划
          </Button>
        }
      />

      {!permissions.canMigrate ? (
        <Alert type="info" showIcon message="你当前只有查看权限，迁移动作已被禁用。" />
      ) : null}
      {migrationMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="迁移计划创建失败"
          description={normalizeApiError(
            migrationMutation.error,
            '迁移计划创建失败，请稍后重试。'
          )}
        />
      ) : null}

      <Card size="small" title={`迁移计划（${plans.length}）`}>
        <Table<MigrationPlan>
          rowKey={(record) => record.id}
          dataSource={plans}
          columns={columns}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <MigrationPlanDrawer
        open={drawerOpen}
        submitting={migrationMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          migrationMutation.mutate(payload, {
            onSuccess: (plan) => {
              setPlans((current) => [plan, ...current]);
              setDrawerOpen(false);
            }
          })
        }
      />
    </Space>
  );
};
