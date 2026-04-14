import { Card, Empty, List, Space, Tag, Typography } from 'antd';
import type { GitOpsConfigurationOverlay } from '@/services/gitops';

type OverlaySummaryPanelProps = {
  overlays?: GitOpsConfigurationOverlay[];
};

export const OverlaySummaryPanel = ({ overlays = [] }: OverlaySummaryPanelProps) => {
  return (
    <Card size="small" title="配置覆盖摘要">
      {overlays.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无配置覆盖" />
      ) : (
        <List
          dataSource={overlays}
          renderItem={(overlay, index) => (
            <List.Item key={`${overlay.overlayRef}-${index}`}>
              <Space direction="vertical" size={0} style={{ width: '100%' }}>
                <Typography.Text strong>{overlay.overlayRef}</Typography.Text>
                <Space wrap>
                  <Tag>{overlay.overlayType}</Tag>
                  <Typography.Text type="secondary">
                    作用域：{overlay.effectiveScope || 'global'}
                  </Typography.Text>
                  <Typography.Text type="secondary">
                    优先级：{overlay.precedence ?? '-'}
                  </Typography.Text>
                </Space>
              </Space>
            </List.Item>
          )}
        />
      )}
    </Card>
  );
};
