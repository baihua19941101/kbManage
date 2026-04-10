import { Button, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { AuditEvent } from '@/services/audit';

export type AuditEventTableProps = {
  data: AuditEvent[];
  loading?: boolean;
  onRefresh?: () => void;
};

const resultColor: Record<string, string> = {
  success: 'success',
  failed: 'error',
  denied: 'warning',
  pending: 'processing'
};

const normalizeText = (value?: string) => value || '-';

const formatDateTime = (value: string) =>
  new Date(value).toLocaleString('zh-CN', { hour12: false });

export const AuditEventTable = ({ data, loading, onRefresh }: AuditEventTableProps) => {
  const columns: ColumnsType<AuditEvent> = [
    {
      title: '发生时间',
      dataIndex: 'occurredAt',
      key: 'occurredAt',
      width: 188,
      render: (value: string) => formatDateTime(value)
    },
    {
      title: '操作者',
      dataIndex: 'actorUserId',
      key: 'actorUserId',
      width: 120,
      render: (value?: string) => normalizeText(value)
    },
    {
      title: '事件类型',
      dataIndex: 'eventType',
      key: 'eventType',
      width: 180,
      render: (value: string) => <Typography.Text code>{value}</Typography.Text>
    },
    {
      title: '集群',
      dataIndex: 'clusterId',
      key: 'clusterId',
      width: 120,
      render: (value?: string) => normalizeText(value)
    },
    {
      title: '目标资源',
      key: 'resource',
      render: (_, record) => {
        const resource = [record.resourceNamespace, record.resourceKind, record.resourceName]
          .filter(Boolean)
          .join('/');
        return normalizeText(resource || record.scopeId);
      }
    },
    {
      title: '结果',
      dataIndex: 'result',
      key: 'result',
      width: 100,
      render: (value: string) => (
        <Tag color={resultColor[value] || 'default'}>{value || '-'}</Tag>
      )
    }
  ];

  return (
    <Space direction="vertical" size="small" style={{ width: '100%' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography.Text type="secondary">共 {data.length} 条审计记录</Typography.Text>
        <Button loading={loading} onClick={onRefresh}>
          刷新
        </Button>
      </div>
      <Table<AuditEvent>
        rowKey="id"
        loading={loading}
        dataSource={data}
        columns={columns}
        pagination={{ pageSize: 10 }}
      />
    </Space>
  );
};
