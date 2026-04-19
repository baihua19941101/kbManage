import { useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { BackupPolicyDrawer } from '@/features/backup-restore-dr/components/BackupPolicyDrawer';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { useRestoreAction } from '@/features/backup-restore-dr/hooks/useRestoreAction';
import { normalizeApiError } from '@/services/api/client';
import {
  backupRestoreQueryKeys,
  createBackupPolicy,
  listBackupPolicies,
  policyQueryScope,
  type BackupPolicy
} from '@/services/backupRestore';

const columns = (
  onRun: (policyId: string) => void,
  canRunBackup: boolean
): ColumnsType<BackupPolicy> => [
  { title: '策略名称', dataIndex: 'name', key: 'name' },
  { title: '范围类型', dataIndex: 'scopeType', key: 'scopeType', render: (value?: string) => value || '—' },
  { title: '执行方式', dataIndex: 'executionMode', key: 'executionMode', render: (value?: string) => value || '—' },
  { title: '保留规则', dataIndex: 'retentionRule', key: 'retentionRule', render: (value?: string) => value || '—' },
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
      <Button size="small" disabled={!canRunBackup} onClick={() => onRun(record.id)}>
        手动备份
      </Button>
    )
  }
];

export const BackupPolicyPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useBackupRestorePermissions();
  const { runBackupMutation } = useRestoreAction();
  const createPolicyMutation = useMutation({ mutationFn: createBackupPolicy });
  const policyQuery = useQuery({
    queryKey: backupRestoreQueryKeys.policies(policyQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listBackupPolicies({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无备份策略访问权限。" />;
  }

  const policies = policyQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="备份策略中心"
        description="定义平台对象保护范围、保留规则和手动备份入口。"
        actions={
          <Button type="primary" disabled={!permissions.canManagePolicy} onClick={() => setDrawerOpen(true)}>
            新建策略
          </Button>
        }
      />

      {policyQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="备份策略加载失败"
          description={normalizeApiError(policyQuery.error, '备份策略加载失败，请稍后重试。')}
        />
      ) : null}

      {createPolicyMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="策略创建失败"
          description={normalizeApiError(createPolicyMutation.error, '策略创建失败，请稍后重试。')}
        />
      ) : null}

      {runBackupMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="手动备份失败"
          description={normalizeApiError(runBackupMutation.error, '手动备份失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`策略列表（${policies.length}）`}>
        <Table<BackupPolicy>
          rowKey={(record) => record.id}
          dataSource={policies}
          columns={columns((policyId) => runBackupMutation.mutate(policyId), permissions.canRunBackup)}
          loading={policyQuery.isLoading || policyQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <BackupPolicyDrawer
        open={drawerOpen}
        submitting={createPolicyMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createPolicyMutation.mutate(payload, {
            onSuccess: () => {
              setDrawerOpen(false);
            }
          })
        }
      />
    </Space>
  );
};
