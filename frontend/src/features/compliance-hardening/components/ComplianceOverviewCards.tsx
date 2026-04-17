import { Card, Space, Statistic } from 'antd';
import type { ComplianceOverview } from '@/services/compliance';

export const ComplianceOverviewCards = ({ overview }: { overview?: ComplianceOverview }) => {
  return (
    <Card size="small" title="覆盖率与风险概览">
      <Space wrap size={24}>
        <Statistic title="覆盖率" value={overview?.coverageRate ?? 0} suffix="%" precision={1} />
        <Statistic title="未关闭失败项" value={overview?.openFindingsCount ?? 0} />
        <Statistic title="高风险遗留" value={overview?.highRiskOpenCount ?? 0} />
        <Statistic
          title="整改完成率"
          value={overview?.remediationCompletionRate ?? 0}
          suffix="%"
          precision={1}
        />
      </Space>
    </Card>
  );
};
