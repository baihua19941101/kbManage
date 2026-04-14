import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { queryKeys } from '@/app/queryClient';
import { normalizeApiError } from '@/services/api/client';
import {
  listGitOpsReleaseRevisions,
  type GitOpsReleaseRevision,
  type ResourceId
} from '@/services/gitops';

type RevisionHistoryPanelProps = {
  unitId: ResourceId;
  environment?: string;
  onRollback?: (revision: GitOpsReleaseRevision) => void;
};

const statusColorMap: Record<string, string> = {
  active: 'green',
  historical: 'default',
  failed: 'red',
  rolled_back: 'orange'
};

const formatDateTime = (value?: string) => {
  if (!value) {
    return '-';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString('zh-CN', {
    hour12: false,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

export const RevisionHistoryPanel = ({ unitId, environment, onRollback }: RevisionHistoryPanelProps) => {
  const revisionsQuery = useQuery({
    queryKey: queryKeys.gitops.deliveryUnits(`revisions:${unitId}:${environment || 'all'}`),
    queryFn: () => listGitOpsReleaseRevisions(unitId, { environment }),
    meta: { suppressGlobalError: true }
  });

  const items = revisionsQuery.data?.items || [];

  const columns: ColumnsType<GitOpsReleaseRevision> = [
    {
      title: '发布 ID',
      dataIndex: 'id',
      key: 'id',
      width: 88
    },
    {
      title: '源版本',
      dataIndex: 'sourceRevision',
      key: 'sourceRevision',
      render: (value?: string) => value || '-'
    },
    {
      title: '应用版本',
      dataIndex: 'appVersion',
      key: 'appVersion',
      render: (value?: string) => value || '-'
    },
    {
      title: '配置版本',
      dataIndex: 'configVersion',
      key: 'configVersion',
      render: (value?: string) => value || '-'
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (value?: string) => (
        <Tag color={statusColorMap[value || 'historical'] || 'default'}>{value || 'historical'}</Tag>
      )
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (value?: string) => formatDateTime(value)
    },
    {
      title: '操作',
      key: 'actions',
      width: 92,
      render: (_, row) =>
        row.rollbackAvailable ? (
          <a onClick={() => onRollback?.(row)} role="button">
            回滚
          </a>
        ) : (
          <Tag>不可回滚</Tag>
        )
    }
  ];

  return (
    <Card size="small" title="发布历史" loading={revisionsQuery.isLoading || revisionsQuery.isFetching}>
      {revisionsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="发布历史加载失败"
          description={normalizeApiError(revisionsQuery.error, '发布历史加载失败，请稍后重试。')}
          style={{ marginBottom: 12 }}
        />
      ) : null}

      {items.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无发布历史" />
      ) : (
        <Table<GitOpsReleaseRevision>
          rowKey={(record) => String(record.id)}
          columns={columns}
          dataSource={items}
          pagination={false}
        />
      )}
    </Card>
  );
};
