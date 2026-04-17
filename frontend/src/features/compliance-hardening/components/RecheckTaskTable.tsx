import { Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { RecheckTask } from '@/services/compliance';
import { recheckStatusColorMap } from '@/features/compliance-hardening/utils';

const columns: ColumnsType<RecheckTask> = [
  {
    title: '复检任务',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '触发来源',
    dataIndex: 'triggerSource',
    key: 'triggerSource',
    render: (value?: string) => value || '—'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: RecheckTask['status']) =>
      value ? <Tag color={recheckStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '结果扫描',
    dataIndex: 'resultScanExecutionId',
    key: 'resultScanExecutionId',
    render: (value?: string) => value || '—'
  },
  {
    title: '摘要',
    dataIndex: 'summary',
    key: 'summary',
    render: (value?: string) => value || '—'
  }
];

export const RecheckTaskTable = ({ tasks, loading }: { tasks: RecheckTask[]; loading?: boolean }) => {
  return (
    <Table<RecheckTask>
      rowKey="id"
      dataSource={tasks}
      loading={loading}
      columns={columns}
      pagination={{ pageSize: 8 }}
      scroll={{ x: 880 }}
    />
  );
};
