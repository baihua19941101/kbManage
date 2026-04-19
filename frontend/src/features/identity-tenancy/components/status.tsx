import { Tag } from 'antd';

const statusColorMap: Record<string, string> = {
  active: 'green',
  healthy: 'green',
  succeeded: 'green',
  completed: 'green',
  mixed: 'blue',
  external: 'blue',
  local: 'gold',
  degraded: 'orange',
  warning: 'orange',
  paused: 'orange',
  draft: 'default',
  blocked: 'volcano',
  revoked: 'red',
  failed: 'red',
  critical: 'red',
  high: 'red',
  medium: 'orange',
  low: 'green'
};

export const StatusTag = ({ value }: { value?: string }) => (
  <Tag color={statusColorMap[value || ''] || 'default'}>{value || '未知'}</Tag>
);
