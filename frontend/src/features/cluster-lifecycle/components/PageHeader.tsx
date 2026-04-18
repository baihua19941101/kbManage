import { Button, Space, Typography } from 'antd';
import type { ReactNode } from 'react';

type PageHeaderProps = {
  title: string;
  description: string;
  actions?: ReactNode;
  extra?: ReactNode;
};

export const PageHeader = ({ title, description, actions, extra }: PageHeaderProps) => (
  <Space direction="vertical" size={12} style={{ width: '100%' }}>
    <Space style={{ width: '100%', justifyContent: 'space-between' }} align="start">
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          {title}
        </Typography.Title>
        <Typography.Text type="secondary">{description}</Typography.Text>
      </div>
      {actions ? <Space wrap>{actions}</Space> : null}
    </Space>
    {extra}
  </Space>
);

export const ComingSoonButton = ({ children }: { children: ReactNode }) => (
  <Button disabled>{children}</Button>
);
