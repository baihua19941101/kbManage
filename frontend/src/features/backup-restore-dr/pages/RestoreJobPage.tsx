import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, List, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { RestoreJobDrawer } from '@/features/backup-restore-dr/components/RestoreJobDrawer';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { useRestoreAction } from '@/features/backup-restore-dr/hooks/useRestoreAction';
import { normalizeApiError } from '@/services/api/client';
import {
  backupRestoreQueryKeys,
  listRestoreJobs,
  listRestorePoints,
  restoreJobQueryScope,
  type PrecheckResult,
  type RestoreJob
} from '@/services/backupRestore';

const columns = (
  onValidate: (jobId: string) => void,
  canRestore: boolean
): ColumnsType<RestoreJob> => [
  { title: '任务 ID', dataIndex: 'id', key: 'id' },
  { title: '任务类型', dataIndex: 'jobType', key: 'jobType', render: (value?: string) => value || '—' },
  { title: '目标环境', dataIndex: 'targetEnvironment', key: 'targetEnvironment', render: (value?: string) => value || '—' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <StatusTag value={value} />
  },
  { title: '一致性说明', dataIndex: 'consistencyNotice', key: 'consistencyNotice', render: (value?: string) => value || '—' },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button size="small" disabled={!canRestore} onClick={() => onValidate(record.id)}>
        执行前校验
      </Button>
    )
  }
];

const PrecheckCard = ({ result }: { result?: PrecheckResult }) => {
  if (!result) {
    return null;
  }

  return (
    <Card size="small" title="最近一次校验结果">
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        <Typography.Text>
          校验状态：<StatusTag value={result.status} />
        </Typography.Text>
        <div>
          <Typography.Text strong>阻断项</Typography.Text>
          <List
            size="small"
            dataSource={result.blockers}
            locale={{ emptyText: '无阻断项' }}
            renderItem={(item) => <List.Item>{item}</List.Item>}
          />
        </div>
        <div>
          <Typography.Text strong>风险提示</Typography.Text>
          <List
            size="small"
            dataSource={result.warnings}
            locale={{ emptyText: '无风险提示' }}
            renderItem={(item) => <List.Item>{item}</List.Item>}
          />
        </div>
        <Typography.Paragraph style={{ marginBottom: 0 }}>
          一致性提示：{result.consistencyNotice || '无'}
        </Typography.Paragraph>
      </Space>
    </Card>
  );
};

export const RestoreJobPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useBackupRestorePermissions();
  const { restoreMutation, precheckMutation } = useRestoreAction();
  const filters = {};
  const jobsQuery = useQuery({
    queryKey: backupRestoreQueryKeys.restoreJobs(restoreJobQueryScope(filters)),
    enabled: permissions.canRead,
    queryFn: () => listRestoreJobs(filters)
  });
  const restorePointsQuery = useQuery({
    queryKey: backupRestoreQueryKeys.restorePoints(),
    enabled: permissions.canRead,
    queryFn: () => listRestorePoints({})
  });

  const lastPrecheck = useMemo(() => precheckMutation.data, [precheckMutation.data]);

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无恢复任务访问权限。" />;
  }

  const jobs = jobsQuery.data?.items || [];
  const restorePoints = restorePointsQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="恢复任务中心"
        description="发起原地恢复、跨集群恢复或定向恢复，并查看恢复前校验结果。"
        actions={
          <Button type="primary" disabled={!permissions.canRestore} onClick={() => setDrawerOpen(true)}>
            发起恢复
          </Button>
        }
      />

      {jobsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="恢复任务列表加载失败"
          description={normalizeApiError(jobsQuery.error, '恢复任务列表加载失败，请稍后重试。')}
        />
      ) : null}
      {restoreMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="恢复任务创建失败"
          description={normalizeApiError(restoreMutation.error, '恢复任务创建失败，请稍后重试。')}
        />
      ) : null}
      {precheckMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="恢复前校验失败"
          description={normalizeApiError(precheckMutation.error, '恢复前校验失败，请稍后重试。')}
        />
      ) : null}

      {!permissions.canRestore ? (
        <Alert
          type="info"
          showIcon
          message="你当前只有查看权限，恢复动作已被禁用。"
        />
      ) : null}

      <Card size="small" title={`恢复任务（${jobs.length}）`}>
        <Table<RestoreJob>
          rowKey={(record) => record.id}
          dataSource={jobs}
          columns={columns((jobId) => precheckMutation.mutate(jobId), permissions.canRestore)}
          loading={jobsQuery.isLoading || jobsQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <PrecheckCard result={lastPrecheck} />

      <RestoreJobDrawer
        open={drawerOpen}
        submitting={restoreMutation.isPending}
        restorePoints={restorePoints}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          restoreMutation.mutate(payload, {
            onSuccess: () => {
              setDrawerOpen(false);
            }
          })
        }
      />
    </Space>
  );
};
