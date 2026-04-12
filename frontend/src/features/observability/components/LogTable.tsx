import { Table } from 'antd';
import type { ObservabilityLogEntryDTO } from '@/services/api/types';

type LogTableProps = {
  loading?: boolean;
  items: ObservabilityLogEntryDTO[];
};

export const LogTable = ({ loading, items }: LogTableProps) => {
  return (
    <Table
      rowKey={(row) => `${row.timestamp}-${row.message}`}
      loading={loading}
      dataSource={items}
      pagination={{ pageSize: 10 }}
      columns={[
        { title: '时间', dataIndex: 'timestamp', key: 'timestamp', width: 220 },
        { title: 'Namespace', dataIndex: 'namespace', key: 'namespace', width: 150 },
        { title: 'Workload', dataIndex: 'workload', key: 'workload', width: 160 },
        { title: '内容', dataIndex: 'message', key: 'message' }
      ]}
    />
  );
};
