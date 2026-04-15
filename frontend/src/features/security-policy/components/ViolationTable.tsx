import { Button, Empty, Space, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { PolicyHitRecordDTO } from '@/services/api/types';

type ViolationTableProps = {
  violations: PolicyHitRecordDTO[];
  loading?: boolean;
  selectedViolationId?: string;
  readonly?: boolean;
  onSelectViolation?: (violation: PolicyHitRecordDTO) => void;
  onUpdateRemediation?: (violation: PolicyHitRecordDTO) => void;
};

const riskColor: Record<string, string> = {
  critical: 'red',
  high: 'volcano',
  medium: 'gold',
  low: 'green'
};

const remediationColor: Record<string, string> = {
  open: 'default',
  in_progress: 'processing',
  mitigated: 'success',
  closed: 'blue'
};

export const ViolationTable = ({
  violations,
  loading,
  selectedViolationId,
  readonly,
  onSelectViolation,
  onUpdateRemediation
}: ViolationTableProps) => {
  const columns: ColumnsType<PolicyHitRecordDTO> = [
    {
      title: '违规对象',
      key: 'target',
      render: (_value, record) =>
        `${record.clusterId || '-'} / ${record.namespace || '-'} / ${record.resourceKind || '-'} / ${record.resourceName || '-'}`
    },
    {
      title: '策略 ID',
      dataIndex: 'policyId',
      key: 'policyId'
    },
    {
      title: '风险级别',
      dataIndex: 'riskLevel',
      key: 'riskLevel',
      render: (value: string) => <Tag color={riskColor[value] || 'default'}>{value}</Tag>
    },
    {
      title: '整改状态',
      dataIndex: 'remediationStatus',
      key: 'remediationStatus',
      render: (value: string) => <Tag color={remediationColor[value] || 'default'}>{value}</Tag>
    },
    {
      title: '发现时间',
      dataIndex: 'detectedAt',
      key: 'detectedAt'
    },
    {
      title: '操作',
      key: 'actions',
      width: 132,
      render: (_value, record) => (
        <Space>
          <Button
            size="small"
            disabled={readonly}
            onClick={(event) => {
              event.stopPropagation();
              onUpdateRemediation?.(record);
            }}
          >
            更新整改
          </Button>
        </Space>
      )
    }
  ];

  return (
    <Table<PolicyHitRecordDTO>
      rowKey={(record) => record.id}
      columns={columns}
      dataSource={violations}
      loading={loading}
      pagination={{ pageSize: 8 }}
      rowClassName={(record) =>
        record.id === selectedViolationId ? 'ant-table-row-selected' : ''
      }
      onRow={(record) => ({
        onClick: () => onSelectViolation?.(record)
      })}
      locale={{
        emptyText: (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={loading ? '正在加载违规记录...' : '暂无违规记录，请调整筛选条件后重试。'}
          />
        )
      }}
    />
  );
};
