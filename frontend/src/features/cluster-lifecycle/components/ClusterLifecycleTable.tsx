import { Button, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { Link } from 'react-router-dom';
import { HealthStatusTag, LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import type { ClusterLifecycleRecord } from '@/services/clusterLifecycle';

type Props = {
  data: ClusterLifecycleRecord[];
  loading?: boolean;
};

const columns: ColumnsType<ClusterLifecycleRecord> = [
  {
    title: '集群',
    key: 'name',
    render: (_, record) => record.displayName || record.name
  },
  {
    title: '模式',
    dataIndex: 'lifecycleMode',
    key: 'lifecycleMode'
  },
  {
    title: '基础设施',
    dataIndex: 'infrastructureType',
    key: 'infrastructureType'
  },
  {
    title: 'K8s 版本',
    dataIndex: 'kubernetesVersion',
    key: 'kubernetesVersion',
    render: (value?: string) => value || '—'
  },
  {
    title: '生命周期状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  },
  {
    title: '健康状态',
    dataIndex: 'healthStatus',
    key: 'healthStatus',
    render: (value?: string) => <HealthStatusTag value={value} />
  },
  {
    title: '接入状态',
    dataIndex: 'registrationStatus',
    key: 'registrationStatus',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Space wrap>
        <Link to={`/cluster-lifecycle/${record.id}`}>查看详情</Link>
        <Button type="link" size="small" href={`/cluster-lifecycle/upgrades?clusterId=${record.id}`}>
          升级
        </Button>
        <Button type="link" size="small" href={`/cluster-lifecycle/node-pools?clusterId=${record.id}`}>
          节点池
        </Button>
      </Space>
    )
  }
];

export const ClusterLifecycleTable = ({ data, loading }: Props) => (
  <Table<ClusterLifecycleRecord>
    rowKey={(record) => record.id}
    loading={loading}
    columns={columns}
    dataSource={data}
    pagination={{ pageSize: 8 }}
  />
);
