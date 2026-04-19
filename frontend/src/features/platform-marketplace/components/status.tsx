import { Tag } from 'antd';

const colorMap: Record<string, string> = {
  active: 'green',
  enabled: 'green',
  available: 'green',
  succeeded: 'green',
  healthy: 'green',
  synced: 'green',
  pending: 'processing',
  syncing: 'processing',
  running: 'processing',
  draft: 'default',
  disabled: 'default',
  offline: 'volcano',
  retired: 'volcano',
  failed: 'red',
  blocked: 'red',
  incompatible: 'red',
  warning: 'gold',
  degraded: 'gold'
};

type StatusTagProps = {
  value?: string;
  fallback?: string;
};

export const StatusTag = ({ value, fallback = '未知' }: StatusTagProps) => (
  <Tag color={value ? colorMap[value] || 'default' : 'default'}>{value || fallback}</Tag>
);
