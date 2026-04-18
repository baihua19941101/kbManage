import { Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { CapabilityStatusTag } from '@/features/cluster-lifecycle/components/status';
import type { CapabilityMatrixEntry } from '@/services/clusterLifecycle';

const columns: ColumnsType<CapabilityMatrixEntry> = [
  {
    title: '能力域',
    dataIndex: 'capabilityDomain',
    key: 'capabilityDomain'
  },
  {
    title: '支持级别',
    dataIndex: 'supportLevel',
    key: 'supportLevel',
    render: (value?: string) => <CapabilityStatusTag value={value} />
  },
  {
    title: '兼容结论',
    dataIndex: 'compatibilityStatus',
    key: 'compatibilityStatus',
    render: (value?: string) => <CapabilityStatusTag value={value} />
  },
  {
    title: '约束说明',
    dataIndex: 'constraintsSummary',
    key: 'constraintsSummary',
    render: (value?: string) => value || '—'
  },
  {
    title: '推荐场景',
    dataIndex: 'recommendedFor',
    key: 'recommendedFor',
    render: (value?: string) =>
      value ? <Typography.Text>{value}</Typography.Text> : <Typography.Text type="secondary">—</Typography.Text>
  }
];

export const CapabilityMatrixTable = ({
  data,
  loading
}: {
  data: CapabilityMatrixEntry[];
  loading?: boolean;
}) => (
  <Table<CapabilityMatrixEntry>
    rowKey={(record) => record.id}
    columns={columns}
    dataSource={data}
    loading={loading}
    pagination={{ pageSize: 10 }}
  />
);
