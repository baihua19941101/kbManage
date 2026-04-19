import { Card, List, Typography } from 'antd';
import type { CapacityBaseline, ScaleEvidence } from '@/services/sreScale';

export const CapacityTrendChart = ({
  baselines,
  evidence
}: {
  baselines: CapacityBaseline[];
  evidence: ScaleEvidence[];
}) => (
  <Card size="small" title="容量与趋势摘要">
    <List
      dataSource={baselines}
      locale={{ emptyText: '暂无容量基线。' }}
      renderItem={(item) => (
        <List.Item>
          <Typography.Text>
            {item.name} / {item.status || 'unknown'} / {item.forecastResult || '无预测结果'}
          </Typography.Text>
        </List.Item>
      )}
    />
    <Typography.Paragraph style={{ marginTop: 12, marginBottom: 0 }}>
      最近证据数：{evidence.length}
    </Typography.Paragraph>
  </Card>
);
