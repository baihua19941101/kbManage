import { Tag } from 'antd';

const COLOR_MAP: Record<string, string> = {
  active: 'green',
  healthy: 'green',
  ready: 'green',
  accepted: 'green',
  warning: 'orange',
  degraded: 'orange',
  scheduled: 'blue',
  rolling: 'processing',
  maintenance: 'gold',
  critical: 'red',
  failed: 'red',
  exception: 'red'
};

export const StatusTag = ({ value }: { value?: string }) => (
  <Tag color={COLOR_MAP[(value || '').toLowerCase()] || 'default'}>{value || 'unknown'}</Tag>
);
