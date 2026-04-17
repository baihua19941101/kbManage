import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Input, Space, Typography } from 'antd';
import { canReadCompliance, useAuthStore } from '@/features/auth/store';
import { ComplianceTrendChart } from '@/features/compliance-hardening/components/ComplianceTrendChart';
import { normalizeApiError } from '@/services/api/client';
import { getComplianceTrends } from '@/services/compliance';

export const ComplianceTrendPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const [baselineId, setBaselineId] = useState('');
  const [scopeRef, setScopeRef] = useState('');

  const queryInput = useMemo(
    () => ({
      baselineId: baselineId.trim() || undefined,
      scopeType: scopeRef.trim() ? ('cluster' as const) : undefined,
      scopeRef: scopeRef.trim() || undefined
    }),
    [baselineId, scopeRef]
  );

  const trendsQuery = useQuery({
    queryKey: ['compliance', 'trends', queryInput],
    enabled: canRead,
    queryFn: () => getComplianceTrends(queryInput)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无趋势复盘访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          合规趋势复盘
        </Typography.Title>
        <Typography.Text type="secondary">
          比较不同时间窗口的得分、覆盖率和整改完成率，识别基线版本变更影响。
        </Typography.Text>
      </div>

      {trendsQuery.error ? (
        <Alert type="error" showIcon message="趋势数据加载失败" description={normalizeApiError(trendsQuery.error, '趋势数据加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选">
        <Space wrap>
          <Input style={{ width: 240 }} placeholder="基线 ID（可选）" value={baselineId} onChange={(event) => setBaselineId(event.target.value)} />
          <Input style={{ width: 240 }} placeholder="范围标识（可选，默认按 cluster）" value={scopeRef} onChange={(event) => setScopeRef(event.target.value)} />
        </Space>
      </Card>

      <ComplianceTrendChart data={trendsQuery.data} />
    </Space>
  );
};
