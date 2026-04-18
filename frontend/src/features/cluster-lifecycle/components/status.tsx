import { Tag } from 'antd';

const statusColorMap: Record<string, string> = {
  active: 'green',
  healthy: 'green',
  connected: 'green',
  succeeded: 'green',
  compatible: 'green',
  native: 'green',
  pending: 'gold',
  warning: 'gold',
  issued: 'gold',
  conditional: 'orange',
  partial: 'orange',
  upgrading: 'blue',
  running: 'blue',
  retiring: 'purple',
  disabled: 'default',
  retired: 'default',
  failed: 'red',
  critical: 'red',
  incompatible: 'red',
  unsupported: 'red',
  blocked: 'red',
  degraded: 'volcano',
  unknown: 'default'
};

const toText = (value?: string) => value || '未知';

export const LifecycleStatusTag = ({ value }: { value?: string }) => (
  <Tag color={statusColorMap[value || 'unknown'] || 'default'}>{toText(value)}</Tag>
);

export const HealthStatusTag = ({ value }: { value?: string }) => (
  <Tag color={statusColorMap[value || 'unknown'] || 'default'}>{toText(value)}</Tag>
);

export const CapabilityStatusTag = ({ value }: { value?: string }) => (
  <Tag color={statusColorMap[value || 'unknown'] || 'default'}>{toText(value)}</Tag>
);
