import { Button, Space, Table, Tag } from 'antd';
import type { ObservabilityAlertDTO } from '@/services/api/types';

type AlertTableProps = {
  loading?: boolean;
  items: ObservabilityAlertDTO[];
  acknowledgeDisabled?: boolean;
  onViewDetail?: (item: ObservabilityAlertDTO) => void;
  onAcknowledge?: (item: ObservabilityAlertDTO) => void;
};

const severityColor: Record<string, string> = {
  critical: 'red',
  warning: 'orange',
  info: 'blue'
};

export const AlertTable = ({
  loading,
  items,
  acknowledgeDisabled,
  onViewDetail,
  onAcknowledge
}: AlertTableProps) => {
  return (
    <Table
      rowKey={(row) => `${row.id}`}
      loading={loading}
      dataSource={items}
      pagination={{ pageSize: 10 }}
      columns={[
        { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
        {
          title: '级别',
          dataIndex: 'severity',
          key: 'severity',
          width: 120,
          render: (value: string) => (
            <Tag color={severityColor[value] ?? 'default'}>{value || 'unknown'}</Tag>
          )
        },
        {
          title: '状态',
          dataIndex: 'status',
          key: 'status',
          width: 140,
          render: (value: string) => <Tag>{value || 'unknown'}</Tag>
        },
        { title: '摘要', dataIndex: 'summary', key: 'summary' },
        {
          title: '操作',
          key: 'actions',
          width: 200,
          render: (_: unknown, record: ObservabilityAlertDTO) => (
            <Space>
              <Button size="small" onClick={() => onViewDetail?.(record)}>
                详情
              </Button>
              <Button
                size="small"
                type="primary"
                disabled={acknowledgeDisabled}
                onClick={() => onAcknowledge?.(record)}
              >
                确认
              </Button>
            </Space>
          )
        }
      ]}
    />
  );
};
