import { Descriptions, Drawer, Space, Tag, Typography } from 'antd';
import type { ObservabilityAlertDTO } from '@/services/api/types';

type AlertDetailDrawerProps = {
  open: boolean;
  alert?: ObservabilityAlertDTO;
  onClose: () => void;
};

export const AlertDetailDrawer = ({ open, alert, onClose }: AlertDetailDrawerProps) => {
  return (
    <Drawer
      title="告警详情"
      open={open}
      onClose={onClose}
      width={520}
      destroyOnHidden
      styles={{ body: { paddingTop: 12 } }}
    >
      {alert ? (
        <Space direction="vertical" style={{ width: '100%' }}>
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="告警 ID">{alert.id}</Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag>{alert.status ?? 'unknown'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="级别">
              <Tag>{alert.severity ?? 'unknown'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="资源类型">{alert.resourceKind ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="资源名称">{alert.resourceName ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="命名空间">{alert.namespace ?? '-'}</Descriptions.Item>
          </Descriptions>
          <Typography.Paragraph style={{ marginBottom: 0 }}>
            {alert.summary ?? '暂无摘要'}
          </Typography.Paragraph>
        </Space>
      ) : (
        <Typography.Text type="secondary">请选择一条告警查看详情。</Typography.Text>
      )}
    </Drawer>
  );
};
