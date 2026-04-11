import { useCallback, useEffect, useState } from 'react';
import { Alert, Button, Card, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { listOperations, type OperationRecord, type OperationStatus } from '@/services/operations';

const actionLabel: Record<OperationRecord['type'], string> = {
  scale: '扩缩容',
  restart: '重启',
  'node-maintenance': '节点维护'
};

const statusColor: Record<OperationStatus, string> = {
  pending: 'default',
  running: 'processing',
  succeeded: 'success',
  failed: 'error'
};

const statusLabel: Record<OperationStatus, string> = {
  pending: '待执行',
  running: '执行中',
  succeeded: '已成功',
  failed: '已失败'
};

const riskColor = {
  low: 'green',
  medium: 'gold',
  high: 'red'
} as const;

const formatDateTime = (value: string) =>
  new Date(value).toLocaleString('zh-CN', { hour12: false });

type OperationCenterPageProps = {
  refreshSignal?: number;
};

export const OperationCenterPage = ({ refreshSignal }: OperationCenterPageProps) => {
  const [data, setData] = useState<OperationRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();

  const reload = useCallback(async () => {
    setLoading(true);
    try {
      const list = await listOperations();
      setData(list);
      setError(undefined);
    } catch (err) {
      setError(err instanceof Error ? err.message : '操作状态刷新失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void reload();
  }, [refreshSignal, reload]);

  useEffect(() => {
    const timer = window.setInterval(() => {
      void reload();
    }, 2500);
    return () => window.clearInterval(timer);
  }, [reload]);

  const columns: ColumnsType<OperationRecord> = [
    {
      title: '操作ID',
      dataIndex: 'id',
      key: 'id',
      width: 220
    },
    {
      title: '操作类型',
      key: 'type',
      render: (_, record) => actionLabel[record.type]
    },
    {
      title: '目标资源',
      key: 'target',
      render: (_, record) =>
        `${record.target.cluster}/${record.target.namespace}/${record.target.name}`
    },
    {
      title: '风险',
      dataIndex: 'riskLevel',
      key: 'riskLevel',
      render: (value: OperationRecord['riskLevel']) => <Tag color={riskColor[value]}>{value}</Tag>
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (value: OperationStatus) => (
        <Tag color={statusColor[value]}>{statusLabel[value]}</Tag>
      )
    },
    {
      title: '原因',
      dataIndex: 'reason',
      key: 'reason',
      ellipsis: true
    },
    {
      title: '结果',
      dataIndex: 'resultMessage',
      key: 'resultMessage',
      ellipsis: true,
      render: (value?: string) => value || '-'
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (value: string) => formatDateTime(value),
      width: 180
    }
  ];

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Card
        size="small"
        title="操作中心"
        extra={
          <Button loading={loading} onClick={() => void reload()}>
            刷新
          </Button>
        }
      >
        <Typography.Text type="secondary">
          展示资源操作的最新状态，包含待执行、执行中、成功与失败。
        </Typography.Text>
      </Card>
      {error ? <Alert type="error" showIcon message="操作中心刷新失败" description={error} /> : null}
      <Table<OperationRecord>
        rowKey="id"
        loading={loading}
        dataSource={data}
        columns={columns}
        pagination={{ pageSize: 8 }}
      />
    </Space>
  );
};
