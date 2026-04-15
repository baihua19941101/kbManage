import { Button, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { SecurityPolicyDTO } from '@/services/api/types';

type PolicyTableProps = {
  loading?: boolean;
  policies: SecurityPolicyDTO[];
  selectedPolicyId?: string;
  readonly?: boolean;
  onSelectPolicy?: (policy: SecurityPolicyDTO) => void;
  onEditPolicy?: (policy: SecurityPolicyDTO) => void;
  onAssignPolicy?: (policy: SecurityPolicyDTO) => void;
};

const modeColorMap: Record<string, string> = {
  audit: 'blue',
  alert: 'gold',
  warn: 'orange',
  enforce: 'red'
};

const statusColorMap: Record<string, string> = {
  draft: 'default',
  active: 'green',
  disabled: 'orange',
  archived: 'default'
};

const riskColorMap: Record<string, string> = {
  low: 'green',
  medium: 'gold',
  high: 'orange',
  critical: 'red'
};

export const PolicyTable = ({
  loading,
  policies,
  selectedPolicyId,
  readonly,
  onSelectPolicy,
  onEditPolicy,
  onAssignPolicy
}: PolicyTableProps) => {
  const columns: ColumnsType<SecurityPolicyDTO> = [
    {
      title: '策略名称',
      dataIndex: 'name',
      key: 'name',
      render: (_value, record) => (
        <Typography.Link onClick={() => onSelectPolicy?.(record)}>{record.name}</Typography.Link>
      )
    },
    {
      title: '层级',
      dataIndex: 'scopeLevel',
      key: 'scopeLevel'
    },
    {
      title: '类别',
      dataIndex: 'category',
      key: 'category'
    },
    {
      title: '默认模式',
      dataIndex: 'defaultEnforcementMode',
      key: 'defaultEnforcementMode',
      render: (mode?: string) => <Tag color={modeColorMap[mode || ''] || 'default'}>{mode || '-'}</Tag>
    },
    {
      title: '风险',
      dataIndex: 'riskLevel',
      key: 'riskLevel',
      render: (risk?: string) => <Tag color={riskColorMap[risk || ''] || 'default'}>{risk || '-'}</Tag>
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status?: string) => <Tag color={statusColorMap[status || ''] || 'default'}>{status || '-'}</Tag>
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      render: (updatedAt?: string) => updatedAt || '-'
    },
    {
      title: '操作',
      key: 'actions',
      render: (_value, record) => (
        <Space size={8}>
          <Button size="small" disabled={readonly} onClick={() => onEditPolicy?.(record)}>
            编辑
          </Button>
          <Button size="small" disabled={readonly} onClick={() => onAssignPolicy?.(record)}>
            分配策略
          </Button>
        </Space>
      )
    }
  ];

  return (
    <Table<SecurityPolicyDTO>
      rowKey={(record) => record.id}
      loading={loading}
      columns={columns}
      dataSource={policies}
      pagination={{ pageSize: 10 }}
      rowClassName={(record) => (record.id === selectedPolicyId ? 'ant-table-row-selected' : '')}
      onRow={(record) => ({
        onClick: () => onSelectPolicy?.(record)
      })}
    />
  );
};
