import { Button, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { ComplianceFinding } from '@/services/compliance';
import { findingResultColorMap, remediationStatusColorMap, riskColorMap } from '@/features/compliance-hardening/utils';

const columns = (
  onView: (finding: ComplianceFinding) => void
): ColumnsType<ComplianceFinding> => [
  {
    title: '控制项',
    key: 'control',
    render: (_, record) => record.controlTitle || record.controlId || '未命名控制项'
  },
  {
    title: '结果',
    dataIndex: 'result',
    key: 'result',
    render: (value?: ComplianceFinding['result']) =>
      value ? <Tag color={findingResultColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '风险',
    dataIndex: 'riskLevel',
    key: 'riskLevel',
    render: (value?: ComplianceFinding['riskLevel']) =>
      value ? <Tag color={riskColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '资源',
    key: 'resource',
    render: (_, record) =>
      [record.clusterId, record.namespace, record.resourceKind, record.resourceName].filter(Boolean).join(' / ') || '—'
  },
  {
    title: '治理状态',
    dataIndex: 'remediationStatus',
    key: 'remediationStatus',
    render: (value?: ComplianceFinding['remediationStatus']) =>
      value ? <Tag color={remediationStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" onClick={() => onView(record)}>
        查看详情
      </Button>
    )
  }
];

type FindingTableProps = {
  findings: ComplianceFinding[];
  loading?: boolean;
  onView: (finding: ComplianceFinding) => void;
};

export const FindingTable = ({ findings, loading, onView }: FindingTableProps) => {
  return (
    <Table<ComplianceFinding>
      rowKey="id"
      loading={loading}
      dataSource={findings}
      columns={columns(onView)}
      pagination={{ pageSize: 8 }}
      scroll={{ x: 960 }}
    />
  );
};
