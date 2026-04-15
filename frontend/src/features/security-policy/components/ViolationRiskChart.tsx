import { Card, Empty, Progress, Space, Typography } from 'antd';
import type { PolicyHitRecordDTO } from '@/services/api/types';

type ViolationRiskChartProps = {
  hits: PolicyHitRecordDTO[];
};

const levels = ['critical', 'high', 'medium', 'low'] as const;

const levelColor: Record<(typeof levels)[number], string> = {
  critical: '#cf1322',
  high: '#d46b08',
  medium: '#faad14',
  low: '#52c41a'
};

export const ViolationRiskChart = ({ hits }: ViolationRiskChartProps) => {
  const total = hits.length;
  const counters = levels.reduce(
    (acc, level) => {
      acc[level] = hits.filter((item) => item.riskLevel === level).length;
      return acc;
    },
    { critical: 0, high: 0, medium: 0, low: 0 }
  );

  return (
    <Card size="small" title="风险级别分布">
      {total === 0 ? (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description="暂无违规数据，无法生成风险分布。"
        />
      ) : (
        <Space direction="vertical" size={10} style={{ width: '100%' }}>
          {levels.map((level) => {
            const count = counters[level];
            const percent = Math.round((count / total) * 100);
            return (
              <div key={level}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <Typography.Text>{level}</Typography.Text>
                  <Typography.Text type="secondary">{count} 条</Typography.Text>
                </div>
                <Progress
                  percent={percent}
                  strokeColor={levelColor[level]}
                  trailColor="#f5f5f5"
                  size="small"
                />
              </div>
            );
          })}
        </Space>
      )}
    </Card>
  );
};
