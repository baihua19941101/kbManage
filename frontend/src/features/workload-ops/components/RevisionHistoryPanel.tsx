import { Card, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { ReleaseRevisionDTO } from '@/services/api/types';

type RevisionHistoryPanelProps = {
  items: ReleaseRevisionDTO[];
  loading?: boolean;
  onRollback?: (item: ReleaseRevisionDTO) => void;
};

export const RevisionHistoryPanel = ({ items, loading, onRollback }: RevisionHistoryPanelProps) => {
  const columns: ColumnsType<ReleaseRevisionDTO> = [
    { title: 'Revision', dataIndex: 'revision', key: 'revision' },
    { title: 'Source', dataIndex: 'sourceName', key: 'sourceName' },
    {
      title: '当前版本',
      dataIndex: 'isCurrent',
      key: 'isCurrent',
      render: (value: boolean) => (value ? <Tag color="green">Current</Tag> : <Tag>History</Tag>)
    },
    {
      title: '可回滚',
      dataIndex: 'rollbackAvailable',
      key: 'rollbackAvailable',
      render: (value: boolean, row) =>
        value ? (
          <a onClick={() => onRollback?.(row)} role="button">
            回滚
          </a>
        ) : (
          <Tag color="default">不可用</Tag>
        )
    }
  ];

  return (
    <Card title="发布历史" size="small">
      <Table rowKey={(row) => `${row.sourceName}-${row.revision}`} loading={loading} dataSource={items} columns={columns} pagination={false} />
    </Card>
  );
};
