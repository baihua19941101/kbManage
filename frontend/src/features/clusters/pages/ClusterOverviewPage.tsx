import { useState } from 'react';
import { Button, Card, List, Space, Tag, Typography } from 'antd';
import { ClusterOnboardDrawer } from '@/features/clusters/components/ClusterOnboardDrawer';

type ClusterSummary = {
  id: string;
  name: string;
  status: 'Connected' | 'Syncing';
  namespaces: number;
};

const initialClusters: ClusterSummary[] = [
  { id: 'c-1', name: 'prod-cn', status: 'Connected', namespaces: 12 },
  { id: 'c-2', name: 'staging-us', status: 'Syncing', namespaces: 5 }
];

export const ClusterOverviewPage = () => {
  const [clusters, setClusters] = useState<ClusterSummary[]>(initialClusters);
  const [drawerOpen, setDrawerOpen] = useState(false);

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            集群总览
          </Typography.Title>
          <Typography.Text type="secondary">管理已接入集群，并发起新的集群接入。</Typography.Text>
        </div>
        <Button type="primary" onClick={() => setDrawerOpen(true)}>
          接入集群
        </Button>
      </div>

      <Card>
        <List
          dataSource={clusters}
          rowKey="id"
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                title={
                  <Space>
                    <span>{item.name}</span>
                    <Tag color={item.status === 'Connected' ? 'green' : 'processing'}>
                      {item.status}
                    </Tag>
                  </Space>
                }
                description={`Namespaces: ${item.namespaces}`}
              />
            </List.Item>
          )}
        />
      </Card>

      <ClusterOnboardDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        onSuccess={(clusterName) =>
          setClusters((previous) => [
            ...previous,
            {
              id: `c-${previous.length + 1}`,
              name: clusterName,
              status: 'Syncing',
              namespaces: 0
            }
          ])
        }
      />
    </Space>
  );
};
