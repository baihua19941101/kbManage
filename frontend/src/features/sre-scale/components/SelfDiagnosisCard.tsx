import { Card, Typography } from 'antd';
import type { ScaleEvidence } from '@/services/sreScale';

export const SelfDiagnosisCard = ({ evidence }: { evidence?: ScaleEvidence }) => (
  <Card size="small" title="自诊断摘要">
    <Typography.Paragraph>{evidence?.summary || '暂无自诊断摘要。'}</Typography.Paragraph>
    <Typography.Text type="secondary">
      瓶颈：{evidence?.bottleneckSummary || '未识别'}
    </Typography.Text>
  </Card>
);
