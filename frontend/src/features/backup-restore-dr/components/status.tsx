import { Tag } from 'antd';

const colorMap: Record<string, string> = {
  active: 'green',
  succeeded: 'green',
  completed: 'green',
  running: 'blue',
  partial: 'gold',
  failed: 'red',
  blocked: 'volcano',
  paused: 'default',
  draft: 'default',
  pending: 'processing',
  scheduled: 'processing'
};

type StatusTagProps = {
  value?: string;
  fallback?: string;
};

export const StatusTag = ({ value, fallback = '未知' }: StatusTagProps) => (
  <Tag color={value ? colorMap[value] || 'default' : 'default'}>{value || fallback}</Tag>
);
