import { useEffect, useMemo, useState } from 'react';
import { Button, Card, Space, Table, Tabs, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { OperationCenterPage } from '@/features/operations/pages/OperationCenterPage';
import { ResourceDetailDrawer, type ResourceItem } from '@/features/resources/components/ResourceDetailDrawer';
import {
  ResourceFilters,
  type ResourceFilterValues
} from '@/features/resources/components/ResourceFilters';
import { listResources } from '@/services/resources';
import type { ResourceListQueryDTO } from '@/services/api/types';

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

const normalizeValue = (value?: string): string | undefined => {
  const normalized = value?.trim();
  return normalized && normalized.length > 0 ? normalized : undefined;
};

const mapFiltersToQuery = (filters: ResourceFilterValues): ResourceListQueryDTO => {
  const query: ResourceListQueryDTO = {};

  const clusterId = normalizeValue(filters.cluster);
  if (clusterId) {
    query.clusterId = clusterId;
  }

  const namespace = normalizeValue(filters.namespace);
  if (namespace) {
    query.namespace = namespace;
  }

  const kind = normalizeValue(filters.resourceType);
  if (kind) {
    query.kind = kind;
  }

  const keyword = normalizeValue(filters.keyword);
  if (keyword) {
    query.keyword = keyword;
  }

  const health = normalizeValue(filters.health);
  if (health) {
    query.health = health;
  }

  return query;
};

export const ResourceListPage = () => {
  const [filters, setFilters] = useState<ResourceFilterValues>({});
  const [resources, setResources] = useState<ResourceItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedResource, setSelectedResource] = useState<ResourceItem | undefined>();
  const [operationRefreshSignal, setOperationRefreshSignal] = useState(0);
  const query = useMemo(() => mapFiltersToQuery(filters), [filters]);

  useEffect(() => {
    let active = true;
    setLoading(true);

    void listResources(query)
      .then((items) => {
        if (!active) {
          return;
        }
        setResources(items);
      })
      .catch(() => {
        if (!active) {
          return;
        }
        setResources([]);
      })
      .finally(() => {
        if (active) {
          setLoading(false);
        }
      });

    return () => {
      active = false;
    };
  }, [query]);

  const clusterOptions = useMemo(
    () => Array.from(new Set(resources.map((item) => item.cluster))),
    [resources]
  );
  const namespaceOptions = useMemo(
    () => Array.from(new Set(resources.map((item) => item.namespace))),
    [resources]
  );
  const resourceTypeOptions = useMemo(
    () => Array.from(new Set(resources.map((item) => item.resourceType))),
    [resources]
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
                dataSource={resources}
                loading={loading}
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
