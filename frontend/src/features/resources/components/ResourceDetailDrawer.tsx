import { Descriptions, Drawer, Space, Tag, Typography } from 'antd';
import { Button } from 'antd';
import { useNavigate } from 'react-router-dom';
import { ResourceActionPanel } from '@/features/resources/components/ResourceActionPanel';

export type ResourceItem = {
  id: string;
  cluster: string;
  namespace: string;
  resourceType: string;
  name: string;
  status: string;
  labels: Record<string, string>;
  updatedAt: string;
};

type ResourceDetailDrawerProps = {
  open: boolean;
  resource?: ResourceItem;
  onClose: () => void;
  onOperationCreated?: () => void;
};

export const ResourceDetailDrawer = ({
  open,
  resource,
  onClose,
  onOperationCreated
}: ResourceDetailDrawerProps) => {
  const navigate = useNavigate();

  const buildObservabilityParams = (item: ResourceItem) =>
    new URLSearchParams({
      clusterId: item.cluster,
      namespace: item.namespace,
      resourceKind: item.resourceType,
      resourceName: item.name,
      workload: item.name,
      subjectType: item.resourceType.toLowerCase() === 'pod' ? 'pod' : 'workload',
      subjectRef: item.name
    });

  return (
    <Drawer
      title={resource ? `资源详情：${resource.name}` : '资源详情'}
      width={560}
      open={open}
      onClose={onClose}
    >
      {resource ? (
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="Name">{resource.name}</Descriptions.Item>
            <Descriptions.Item label="Type">{resource.resourceType}</Descriptions.Item>
            <Descriptions.Item label="Cluster">{resource.cluster}</Descriptions.Item>
            <Descriptions.Item label="Namespace">{resource.namespace}</Descriptions.Item>
            <Descriptions.Item label="Status">
              <Tag color={resource.status === 'Running' ? 'green' : 'gold'}>{resource.status}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Updated">{resource.updatedAt}</Descriptions.Item>
          </Descriptions>
          <div>
            <Typography.Text strong>Labels</Typography.Text>
            <div style={{ marginTop: 8 }}>
              {Object.entries(resource.labels).map(([key, value]) => (
                <Tag key={key}>{`${key}=${value}`}</Tag>
              ))}
            </div>
          </div>
          <Space wrap>
            <Button
              onClick={() => {
                const params = buildObservabilityParams(resource);
                void navigate(`/observability?${params.toString()}`);
              }}
            >
              可观测总览
            </Button>
            <Button
              onClick={() => {
                const params = buildObservabilityParams(resource);
                void navigate(`/observability/logs?${params.toString()}`);
              }}
            >
              日志
            </Button>
            <Button
              onClick={() => {
                const params = buildObservabilityParams(resource);
                void navigate(`/observability/events?${params.toString()}`);
              }}
            >
              事件
            </Button>
            <Button
              onClick={() => {
                const params = buildObservabilityParams(resource);
                void navigate(`/observability/metrics?${params.toString()}`);
              }}
            >
              指标
            </Button>
            <Button
              onClick={() => {
                const params = buildObservabilityParams(resource);
                void navigate(`/observability/context?${params.toString()}`);
              }}
            >
              资源上下文
            </Button>
          </Space>
          <ResourceActionPanel resource={resource} onOperationCreated={onOperationCreated} />
        </Space>
      ) : null}
    </Drawer>
  );
};
