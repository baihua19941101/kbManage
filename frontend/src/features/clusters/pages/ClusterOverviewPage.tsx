import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, List, Space, Tag, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import { normalizeErrorMessage } from '@/app/queryClient';
import { ClusterOnboardDrawer } from '@/features/clusters/components/ClusterOnboardDrawer';
import { listClusters } from '@/services/clusters';

type ClusterSummary = {
  id: string;
  name: string;
  status: 'Connected' | 'Syncing' | 'Degraded';
  namespaces: number;
};

const mapStatus = (status: string): ClusterSummary['status'] => {
  const normalized = status.trim().toLowerCase();

  if (normalized === 'healthy' || normalized === 'connected' || normalized === 'ready') {
    return 'Connected';
  }

  if (
    normalized === 'degraded' ||
    normalized === 'error' ||
    normalized === 'failed' ||
    normalized === 'unhealthy'
  ) {
    return 'Degraded';
  }

  return 'Syncing';
};

const statusColor = (status: ClusterSummary['status']) => {
  if (status === 'Connected') {
    return 'green';
  }
  if (status === 'Degraded') {
    return 'red';
  }
  return 'processing';
};

export const ClusterOverviewPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const navigate = useNavigate();
  const { data, isFetching, error, refetch } = useQuery({
    queryKey: ['clusters'],
    queryFn: listClusters,
    meta: {
      suppressGlobalError: true
    }
  });

  const clusters: ClusterSummary[] = (data || []).map((item) => ({
    id: item.id,
    name: item.name,
    status: mapStatus(item.status),
    namespaces: item.namespaces
  }));

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
        {error ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="集群列表加载失败"
            description={normalizeErrorMessage(error)}
          />
        ) : null}
        <List
          loading={isFetching}
          dataSource={clusters}
          rowKey="id"
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                title={
                  <Space>
                    <span>{item.name}</span>
                    <Tag color={statusColor(item.status)}>{item.status}</Tag>
                  </Space>
                }
                description={`Namespaces: ${item.namespaces}`}
              />
              <Space wrap>
                <Button
                  onClick={() => {
                    const params = new URLSearchParams({
                      clusterId: item.id,
                      subjectType: 'cluster',
                      subjectRef: item.id
                    });
                    void navigate(`/observability?${params.toString()}`);
                  }}
                >
                  总览
                </Button>
                <Button
                  onClick={() => {
                    const params = new URLSearchParams({
                      clusterId: item.id,
                      subjectType: 'cluster',
                      subjectRef: item.id
                    });
                    void navigate(`/observability/logs?${params.toString()}`);
                  }}
                >
                  日志
                </Button>
                <Button
                  onClick={() => {
                    const params = new URLSearchParams({
                      clusterId: item.id,
                      subjectType: 'cluster',
                      subjectRef: item.id
                    });
                    void navigate(`/observability/events?${params.toString()}`);
                  }}
                >
                  事件
                </Button>
                <Button
                  onClick={() => {
                    const params = new URLSearchParams({
                      clusterId: item.id,
                      subjectType: 'cluster',
                      subjectRef: item.id
                    });
                    void navigate(`/observability/metrics?${params.toString()}`);
                  }}
                >
                  指标
                </Button>
              </Space>
            </List.Item>
          )}
        />
      </Card>

      <ClusterOnboardDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        onSuccess={() => {
          void refetch();
        }}
      />
    </Space>
  );
};
