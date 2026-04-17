import { Button, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { RemediationTask } from '@/services/compliance';
import { formatDateTime, remediationTaskStatusColorMap, riskColorMap } from '@/features/compliance-hardening/utils';

type RemediationTaskTableProps = {
  tasks: RemediationTask[];
  loading?: boolean;
  readonly?: boolean;
  onEdit: (task: RemediationTask) => void;
};

const columns = (
  readonly: boolean,
  onEdit: (task: RemediationTask) => void
): ColumnsType<RemediationTask> => [
  {
    title: '任务',
    key: 'title',
    render: (_, record) => record.title || record.id
  },
  {
    title: '负责人',
    dataIndex: 'owner',
    key: 'owner',
    render: (value?: string) => value || '—'
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    render: (value?: RemediationTask['priority']) =>
      value ? <Tag color={riskColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: RemediationTask['status']) =>
      value ? <Tag color={remediationTaskStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '到期时间',
    dataIndex: 'dueAt',
    key: 'dueAt',
    render: formatDateTime
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" disabled={readonly} onClick={() => onEdit(record)}>
        更新
      </Button>
    )
  }
];

export const RemediationTaskTable = ({ tasks, loading, readonly, onEdit }: RemediationTaskTableProps) => {
  return (
    <Table<RemediationTask>
      rowKey="id"
      dataSource={tasks}
      loading={loading}
      columns={columns(Boolean(readonly), onEdit)}
      pagination={{ pageSize: 8 }}
      scroll={{ x: 880 }}
    />
  );
};
