import type { ReactNode } from 'react';
import { Space, Typography } from 'antd';

type PageHeaderProps = {
  title: string;
  description: string;
  actions?: ReactNode;
};

export const PageHeader = ({ title, description, actions }: PageHeaderProps) => (
  <Space
    align="start"
    style={{ width: '100%', justifyContent: 'space-between', gap: 16 }}
    wrap
  >
    <Space direction="vertical" size={4}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        {title}
      </Typography.Title>
      <Typography.Text type="secondary">{description}</Typography.Text>
    </Space>
    {actions}
  </Space>
);
