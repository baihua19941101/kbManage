import { Card, Skeleton, Typography } from 'antd';

type GitOpsSkeletonPanelProps = {
  title: string;
  description: string;
  loading?: boolean;
};

export const GitOpsSkeletonPanel = ({
  title,
  description,
  loading = false
}: GitOpsSkeletonPanelProps) => {
  return (
    <Card size="small" title={title}>
      {loading ? (
        <Skeleton active paragraph={{ rows: 2 }} title={false} />
      ) : (
        <Typography.Text type="secondary">{description}</Typography.Text>
      )}
    </Card>
  );
};
