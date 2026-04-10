import { useState } from 'react';
import { Button, Card, Space, Typography, message } from 'antd';
import { OperationConfirmDrawer } from '@/features/operations/components/OperationConfirmDrawer';
import type { ResourceItem } from '@/features/resources/components/ResourceDetailDrawer';
import {
  createOperation,
  type CreateOperationPayload,
  type OperationType
} from '@/services/operations';

type ResourceActionPanelProps = {
  resource: ResourceItem;
  onOperationCreated?: () => void;
};

export const ResourceActionPanel = ({ resource, onOperationCreated }: ResourceActionPanelProps) => {
  const [actionType, setActionType] = useState<OperationType>('restart');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const openDrawer = (type: OperationType) => {
    setActionType(type);
    setDrawerOpen(true);
  };

  const handleConfirm = async (payload: Omit<CreateOperationPayload, 'target'>) => {
    setSubmitting(true);
    try {
      await createOperation({
        ...payload,
        target: {
          resourceId: resource.id,
          name: resource.name,
          resourceType: resource.resourceType,
          cluster: resource.cluster,
          namespace: resource.namespace
        }
      });
      message.success('操作已提交，状态将在操作中心刷新。');
      onOperationCreated?.();
    } catch (error) {
      message.error(error instanceof Error ? error.message : '操作提交失败');
      throw error;
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card size="small" title="资源操作">
      <Space direction="vertical" style={{ width: '100%' }}>
        <Typography.Text type="secondary">
          对资源执行高风险变更前，请确认维护窗口和影响范围。
        </Typography.Text>
        <Space wrap>
          <Button onClick={() => openDrawer('scale')}>扩缩容</Button>
          <Button onClick={() => openDrawer('restart')}>重启</Button>
          <Button danger onClick={() => openDrawer('node-maintenance')}>
            节点维护
          </Button>
        </Space>
      </Space>
      <OperationConfirmDrawer
        open={drawerOpen}
        actionType={actionType}
        resourceName={resource.name}
        onClose={() => setDrawerOpen(false)}
        onConfirm={handleConfirm}
        submitting={submitting}
      />
    </Card>
  );
};
