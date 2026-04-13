import { Badge, Button, Card, Space, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { WorkloadInstanceDTO } from '@/services/api/types';

type InstanceListPanelProps = {
  items: WorkloadInstanceDTO[];
  loading?: boolean;
  onOpenTerminal?: (item: WorkloadInstanceDTO) => void;
  terminalDisabled?: boolean;
};

export const InstanceListPanel = ({
  items,
  loading,
  onOpenTerminal,
  terminalDisabled
}: InstanceListPanelProps) => {
  const columns: ColumnsType<WorkloadInstanceDTO> = [
    {
      title: 'Pod',
      dataIndex: 'podName',
      key: 'podName'
    },
    {
      title: 'Container',
      dataIndex: 'containerName',
      key: 'containerName',
      render: (value: string | undefined) => value || '-'
    },
    {
      title: 'Phase',
      dataIndex: 'phase',
      key: 'phase',
      render: (value: string) => <Badge status={value === 'Running' ? 'success' : 'warning'} text={value} />
    },
    {
      title: 'Ready',
      dataIndex: 'ready',
      key: 'ready',
      render: (value: boolean) => (value ? <Tag color="green">Ready</Tag> : <Tag color="gold">NotReady</Tag>)
    },
    {
      title: '重启次数',
      dataIndex: 'restartCount',
      key: 'restartCount',
      render: (value: number | undefined) => value ?? 0
    },
    {
      title: '操作',
      key: 'action',
      render: (_, row) => (
        <Space>
          <Button
            size="small"
            disabled={terminalDisabled || !row.terminalAvailable}
            onClick={() => onOpenTerminal?.(row)}
          >
            终端
          </Button>
        </Space>
      )
    }
  ];

  return (
    <Card title="实例列表" size="small">
      <Table
        rowKey={(record) => `${record.podName}-${record.containerName ?? 'default'}`}
        loading={loading}
        dataSource={items}
        columns={columns}
        pagination={false}
      />
    </Card>
  );
};
