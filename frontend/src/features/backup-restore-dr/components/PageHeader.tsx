import type { ReactNode } from 'react';
import { Typography } from 'antd';

type PageHeaderProps = {
  title: string;
  description: string;
  actions?: ReactNode;
};

export const PageHeader = ({ title, description, actions }: PageHeaderProps) => (
  <div
    style={{ width: '100%', flexWrap: 'wrap' }}
  >
    <div
      style={{
        display: 'flex',
        alignItems: 'flex-start',
        justifyContent: 'space-between',
        gap: 16,
        width: '100%',
        flexWrap: 'wrap'
      }}
    >
      <Typography.Title level={3} style={{ marginBottom: 8 }}>
        {title}
      </Typography.Title>
      {actions ? <div>{actions}</div> : null}
    </div>
    <div style={{ marginTop: 8 }}>
      <Typography.Text type="secondary">{description}</Typography.Text>
    </div>
  </div>
);
