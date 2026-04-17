import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { canReadCompliance, useAuthStore } from '@/features/auth/store';
import { ComplianceOverviewCards } from '@/features/compliance-hardening/components/ComplianceOverviewCards';
import { normalizeApiError } from '@/services/api/client';
import {
  getComplianceOverview,
  type ComplianceOverviewGroup,
  type OverviewGroupBy
} from '@/services/compliance';

const groupByOptions: Array<{ label: string; value: OverviewGroupBy }> = [
  { label: '按集群', value: 'cluster' },
  { label: '按工作空间', value: 'workspace' },
  { label: '按项目', value: 'project' },
  { label: '按基线', value: 'baseline' }
];

const columns: ColumnsType<ComplianceOverviewGroup> = [
  { title: '分组', dataIndex: 'groupKey', key: 'groupKey' },
  { title: '平均得分', dataIndex: 'scoreAvg', key: 'scoreAvg', render: (value?: number) => value ?? '—' },
  { title: '覆盖率', dataIndex: 'coverageRate', key: 'coverageRate', render: (value?: number) => (typeof value === 'number' ? `${value}%` : '—') },
  { title: '未关闭失败项', dataIndex: 'openFindingsCount', key: 'openFindingsCount', render: (value?: number) => value ?? 0 },
  { title: '高风险遗留', dataIndex: 'highRiskOpenCount', key: 'highRiskOpenCount', render: (value?: number) => value ?? 0 }
];

export const ComplianceOverviewPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const [groupBy, setGroupBy] = useState<OverviewGroupBy>('cluster');

  const queryInput = useMemo(() => ({ groupBy }), [groupBy]);
  const overviewQuery = useQuery({
    queryKey: ['compliance', 'overview', queryInput],
    enabled: canRead,
    queryFn: () => getComplianceOverview(queryInput)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无合规总览访问权限。" />;
  }

  const overview = overviewQuery.data;
  const groups = overview?.groups || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          合规覆盖率总览
        </Typography.Title>
        <Typography.Text type="secondary">
          以集群、团队或基线维度查看覆盖率、遗留风险和整改完成率。
        </Typography.Text>
      </div>

      {overviewQuery.error ? (
        <Alert type="error" showIcon message="合规总览加载失败" description={normalizeApiError(overviewQuery.error, '合规总览加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="分析维度">
        <Select style={{ width: 220 }} options={groupByOptions} value={groupBy} onChange={setGroupBy} />
      </Card>

      <ComplianceOverviewCards overview={overview} />

      <Card size="small" title={`分组列表（${groups.length}）`}>
        <Table<ComplianceOverviewGroup>
          rowKey={(record) => record.groupKey || 'unknown'}
          dataSource={groups}
          loading={overviewQuery.isLoading || overviewQuery.isFetching}
          columns={columns}
          pagination={{ pageSize: 8 }}
        />
      </Card>
    </Space>
  );
};
