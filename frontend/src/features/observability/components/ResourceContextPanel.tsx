import { Card, Descriptions, Space, Typography } from 'antd';
import type { ObservabilityOverviewDTO } from '@/services/api/types';

type ResourceContextPanelProps = {
  clusterId: string;
  namespace: string;
  resourceKind: string;
  resourceName: string;
  overview?: ObservabilityOverviewDTO;
};

export const ResourceContextPanel = ({
  clusterId,
  namespace,
  resourceKind,
  resourceName,
  overview
}: ResourceContextPanelProps) => {
  return (
    <Card>
      <Space direction="vertical" style={{ width: '100%' }}>
        <Typography.Title level={5} style={{ margin: 0 }}>
          资源上下文
        </Typography.Title>
        <Descriptions size="small" bordered column={1}>
          <Descriptions.Item label="Cluster">{clusterId}</Descriptions.Item>
          <Descriptions.Item label="Namespace">{namespace}</Descriptions.Item>
          <Descriptions.Item label="Kind">{resourceKind}</Descriptions.Item>
          <Descriptions.Item label="Name">{resourceName}</Descriptions.Item>
        </Descriptions>
        <Typography.Text type="secondary">
          当前概览卡片数量：{overview?.cards?.length ?? 0}
        </Typography.Text>
      </Space>
    </Card>
  );
};
