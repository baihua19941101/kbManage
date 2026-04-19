import { Card, Descriptions } from 'antd';
import type { PlatformHealthOverview } from '@/services/sreScale';
import { StatusTag } from '@/features/sre-scale/components/status';

export const HealthStatusCard = ({ overview }: { overview?: PlatformHealthOverview }) => (
  <Card size="small" title="平台健康总览">
    <Descriptions column={1} size="small" bordered>
      <Descriptions.Item label="总体状态">
        <StatusTag value={overview?.overallStatus} />
      </Descriptions.Item>
      <Descriptions.Item label="组件健康">{overview?.componentHealthSummary || '—'}</Descriptions.Item>
      <Descriptions.Item label="依赖状态">{overview?.dependencyHealthSummary || '—'}</Descriptions.Item>
      <Descriptions.Item label="任务积压">{overview?.taskBacklogSummary || '—'}</Descriptions.Item>
      <Descriptions.Item label="容量风险">{overview?.capacityRiskLevel || '—'}</Descriptions.Item>
      <Descriptions.Item label="限流状态">{overview?.throttlingStatus || '—'}</Descriptions.Item>
      <Descriptions.Item label="恢复摘要">{overview?.recoverySummary || '—'}</Descriptions.Item>
    </Descriptions>
  </Card>
);
