import { useMemo, useState } from 'react';
import { Button, Card, Space, Table, Tabs, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { OperationCenterPage } from '@/features/operations/pages/OperationCenterPage';
import { ResourceDetailDrawer, type ResourceItem } from '@/features/resources/components/ResourceDetailDrawer';
import {
  ResourceFilters,
  type ResourceFilterValues
} from '@/features/resources/components/ResourceFilters';

const mockResources: ResourceItem[] = [
  {
    id: 'res-1',
    cluster: 'prod-cn',
    namespace: 'payments',
    resourceType: 'Deployment',
    name: 'payment-api',
    status: 'Running',
    labels: { app: 'payment-api', env: 'prod' },
    updatedAt: '2026-04-09 11:20'
  },
  {
    id: 'res-2',
    cluster: 'prod-cn',
    namespace: 'gateway',
    resourceType: 'Service',
    name: 'edge-gateway',
    status: 'Running',
    labels: { app: 'gateway', env: 'prod' },
    updatedAt: '2026-04-09 11:10'
  },
  {
    id: 'res-3',
    cluster: 'staging-us',
    namespace: 'payments',
    resourceType: 'Pod',
    name: 'payment-api-66d8b87f74-z2vfh',
    status: 'Pending',
    labels: { app: 'payment-api', env: 'staging' },
    updatedAt: '2026-04-09 10:58'
  }
];

const columns: ColumnsType<ResourceItem> = [
  { title: 'Name', dataIndex: 'name', key: 'name' },
  { title: 'Type', dataIndex: 'resourceType', key: 'resourceType' },
  { title: 'Cluster', dataIndex: 'cluster', key: 'cluster' },
  { title: 'Namespace', dataIndex: 'namespace', key: 'namespace' },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    render: (status: string) => (
      <Tag color={status === 'Running' ? 'green' : 'gold'}>{status}</Tag>
    )
  }
];

export const ResourceListPage = () => {
  const [filters, setFilters] = useState<ResourceFilterValues>({});
  const [selectedResource, setSelectedResource] = useState<ResourceItem | undefined>();
  const [operationRefreshSignal, setOperationRefreshSignal] = useState(0);

  const clusterOptions = useMemo(
    () => Array.from(new Set(mockResources.map((item) => item.cluster))),
    []
  );
  const namespaceOptions = useMemo(
    () => Array.from(new Set(mockResources.map((item) => item.namespace))),
    []
  );
  const resourceTypeOptions = useMemo(
    () => Array.from(new Set(mockResources.map((item) => item.resourceType))),
    []
  );

  const filteredResources = useMemo(
    () =>
      mockResources.filter((item) => {
        const matchCluster = !filters.cluster || item.cluster === filters.cluster;
        const matchNamespace = !filters.namespace || item.namespace === filters.namespace;
        const matchResourceType = !filters.resourceType || item.resourceType === filters.resourceType;
        const normalizedKeyword = filters.keyword?.trim().toLowerCase() ?? '';
        const matchKeyword =
          !normalizedKeyword ||
          item.name.toLowerCase().includes(normalizedKeyword) ||
          Object.entries(item.labels).some(([key, value]) =>
            `${key}=${value}`.toLowerCase().includes(normalizedKeyword)
          );

        return matchCluster && matchNamespace && matchResourceType && matchKeyword;
      }),
    [filters]
  );

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Tabs
        items={[
        {
          key: 'resources',
          label: '资源列表',
          children: (
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <div>
                <Typography.Title level={3} style={{ marginBottom: 8 }}>
                  资源列表
                </Typography.Title>
                <Typography.Text type="secondary">
                  支持按集群、命名空间、资源类型与关键字联合筛选。
                </Typography.Text>
              </div>

              <Card>
                <ResourceFilters
                  values={filters}
                  clusterOptions={clusterOptions}
                  namespaceOptions={namespaceOptions}
                  resourceTypeOptions={resourceTypeOptions}
                  onChange={setFilters}
                  onReset={() => setFilters({})}
                />
              </Card>

              <Table<ResourceItem>
                rowKey="id"
                dataSource={filteredResources}
                columns={[
                  ...columns,
                  {
                    title: 'Action',
                    key: 'action',
                    render: (_, record) => (
                      <Button type="link" onClick={() => setSelectedResource(record)}>
                        查看详情
                      </Button>
                    )
                  }
                ]}
                pagination={{ pageSize: 10 }}
              />
            </Space>
          )
        },
        {
          key: 'operations',
          label: '操作中心',
          children: <OperationCenterPage refreshSignal={operationRefreshSignal} />
        }
        ]}
      />
      <ResourceDetailDrawer
        open={Boolean(selectedResource)}
        resource={selectedResource}
        onClose={() => setSelectedResource(undefined)}
        onOperationCreated={() =>
          setOperationRefreshSignal((currentSignal) => currentSignal + 1)
        }
      />
    </Space>
  );
};
