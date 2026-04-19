import { Alert, Card, List, Typography } from 'antd';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import type { CompatibilityStatement } from '@/services/platformMarketplace';

type CompatibilitySummaryCardProps = {
  compatibility?: CompatibilityStatement;
  loading?: boolean;
};

export const CompatibilitySummaryCard = ({
  compatibility,
  loading
}: CompatibilitySummaryCardProps) => (
  <Card loading={loading} size="small" title="兼容性摘要">
    {compatibility ? (
      <>
        <Typography.Paragraph>
          平台版本：{compatibility.platformVersion || '—'}，结论：
          <StatusTag value={compatibility.compatibilityStatus} />
        </Typography.Paragraph>
        <Typography.Paragraph>{compatibility.summary || '暂无兼容性摘要。'}</Typography.Paragraph>
        {compatibility.permissionImpact ? (
          <Alert
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
            message="权限影响"
            description={compatibility.permissionImpact}
          />
        ) : null}
        <List
          size="small"
          header="阻断原因"
          locale={{ emptyText: '当前没有阻断原因。' }}
          dataSource={compatibility.blockedReasons}
          renderItem={(item) => <List.Item>{item}</List.Item>}
        />
        <List
          size="small"
          header="建议动作"
          locale={{ emptyText: '当前没有额外建议。' }}
          dataSource={compatibility.suggestedActions}
          renderItem={(item) => <List.Item>{item}</List.Item>}
        />
      </>
    ) : (
      <Typography.Text type="secondary">请选择扩展查看兼容性结论。</Typography.Text>
    )}
  </Card>
);
