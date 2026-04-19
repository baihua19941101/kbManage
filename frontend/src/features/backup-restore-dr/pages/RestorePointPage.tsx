import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Button, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { RestorePointDetailDrawer } from '@/features/backup-restore-dr/components/RestorePointDetailDrawer';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  backupRestoreQueryKeys,
  listRestorePoints,
  restorePointQueryScope,
  type RestorePoint
} from '@/services/backupRestore';

const columns = (onView: (restorePoint: RestorePoint) => void): ColumnsType<RestorePoint> => [
  { title: '恢复点 ID', dataIndex: 'id', key: 'id' },
  { title: '策略 ID', dataIndex: 'policyId', key: 'policyId', render: (value?: string) => value || '—' },
  {
    title: '结果',
    dataIndex: 'result',
    key: 'result',
    render: (value?: string) => <StatusTag value={value} />
  },
  {
    title: '耗时',
    dataIndex: 'durationSeconds',
    key: 'durationSeconds',
    render: (value?: number) => (typeof value === 'number' ? `${value} 秒` : '—')
  },
  { title: '到期时间', dataIndex: 'expiresAt', key: 'expiresAt', render: (value?: string) => value || '—' },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button size="small" onClick={() => onView(record)}>
        查看详情
      </Button>
    )
  }
];

export const RestorePointPage = () => {
  const permissions = useBackupRestorePermissions();
  const [selectedRestorePoint, setSelectedRestorePoint] = useState<RestorePoint>();
  const restorePointQuery = useQuery({
    queryKey: backupRestoreQueryKeys.restorePoints(restorePointQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listRestorePoints({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无恢复点访问权限。" />;
  }

  const restorePoints = restorePointQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="恢复点目录"
        description="查看备份结果、失败原因、一致性说明和恢复点有效期。"
      />

      {restorePointQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="恢复点列表加载失败"
          description={normalizeApiError(
            restorePointQuery.error,
            '恢复点列表加载失败，请稍后重试。'
          )}
        />
      ) : null}

      <Card size="small" title={`恢复点列表（${restorePoints.length}）`}>
        <Table<RestorePoint>
          rowKey={(record) => record.id}
          dataSource={restorePoints}
          columns={columns((record) => setSelectedRestorePoint(record))}
          loading={restorePointQuery.isLoading || restorePointQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <RestorePointDetailDrawer
        open={Boolean(selectedRestorePoint)}
        restorePoint={selectedRestorePoint}
        onClose={() => setSelectedRestorePoint(undefined)}
      />
    </Space>
  );
};
