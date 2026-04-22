import { Card, Typography } from 'antd';

export const ReadinessSummaryCard = ({
  summary,
  conclusion
}: {
  summary?: string;
  conclusion?: string;
}) => (
  <Card size="small" title="交付就绪摘要">
    <Typography.Text>{summary || '暂无摘要'}</Typography.Text>
    <br />
    <Typography.Text type="secondary">{conclusion || '未知结论'}</Typography.Text>
  </Card>
);
