import { Card, List } from 'antd';
import type { KeyOperationTrace } from '@/services/enterprisePolish';

export const RiskTrendChart = ({ items }: { items: KeyOperationTrace[] }) => (
  <Card size="small" title="风险趋势摘要">
    <List
      dataSource={items}
      renderItem={(item) => (
        <List.Item>
          {item.operationType || '未知操作'} / {item.riskLevel || '未知等级'} / {item.outcome || '未知结果'}
        </List.Item>
      )}
    />
  </Card>
);
