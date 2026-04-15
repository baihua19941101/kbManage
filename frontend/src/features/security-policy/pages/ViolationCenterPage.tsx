import { useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space, Typography } from 'antd';
import {
  canManagePolicy,
  canReadPolicy,
  useAuthStore
} from '@/features/auth/store';
import { RemediationUpdateDrawer } from '@/features/security-policy/components/RemediationUpdateDrawer';
import { ViolationRiskChart } from '@/features/security-policy/components/ViolationRiskChart';
import { ViolationTable } from '@/features/security-policy/components/ViolationTable';
import {
  isAuthorizationError,
  normalizeApiError,
  normalizeAuthorizationError
} from '@/services/api/client';
import {
  listPolicyHits,
  listSecurityPolicies,
  type PolicyHitListQuery
} from '@/services/securityPolicy';
import type {
  PolicyHitRecordDTO,
  PolicyRemediationStatus,
  SecurityPolicyRiskLevel
} from '@/services/api/types';

const EMPTY_VIOLATIONS: PolicyHitRecordDTO[] = [];

const riskLevelOptions = [
  { label: '全部风险', value: '' },
  { label: 'critical', value: 'critical' },
  { label: 'high', value: 'high' },
  { label: 'medium', value: 'medium' },
  { label: 'low', value: 'low' }
];

const remediationOptions = [
  { label: '全部状态', value: '' },
  { label: 'open', value: 'open' },
  { label: 'in_progress', value: 'in_progress' },
  { label: 'mitigated', value: 'mitigated' },
  { label: 'closed', value: 'closed' }
];

export const ViolationCenterPage = () => {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadPolicy(user);
  const canManage = canManagePolicy(user);

  const [policyId, setPolicyId] = useState<string>('');
  const [riskLevel, setRiskLevel] = useState<string>('');
  const [remediationStatus, setRemediationStatus] = useState<string>('');
  const [selectedViolation, setSelectedViolation] = useState<PolicyHitRecordDTO>();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const policiesQuery = useQuery({
    queryKey: ['securityPolicy', 'violations', 'policies'],
    enabled: canRead,
    queryFn: () => listSecurityPolicies({ status: 'active' })
  });

  const hitQueryInput = useMemo<PolicyHitListQuery>(
    () => ({
      policyId: policyId || undefined,
      riskLevel: (riskLevel || undefined) as SecurityPolicyRiskLevel | undefined,
      remediationStatus: (remediationStatus || undefined) as PolicyRemediationStatus | undefined
    }),
    [policyId, riskLevel, remediationStatus]
  );

  const hitsQuery = useQuery({
    queryKey: ['securityPolicy', 'violations', hitQueryInput],
    enabled: canRead,
    queryFn: () => listPolicyHits(hitQueryInput)
  });

  const violations = hitsQuery.data?.items ?? EMPTY_VIOLATIONS;
  const authorizationChanged =
    isAuthorizationError(policiesQuery.error) || isAuthorizationError(hitsQuery.error);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无违规中心访问权限，请联系管理员授予 policy:read 或对应平台角色。"
      />
    );
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          违规中心
        </Typography.Title>
        <Typography.Text type="secondary">
          统一查询违规对象、查看风险分布，并更新整改状态形成处置闭环。
        </Typography.Text>
      </div>

      {!canManage ? (
        <Alert
          type="info"
          showIcon
          message="当前为只读模式"
          description="你可查询违规与风险分布，但无法更新整改状态。"
        />
      ) : null}

      {authorizationChanged ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description={normalizeAuthorizationError(
            policiesQuery.error || hitsQuery.error,
            '当前账号可能缺少违规中心所需权限，请刷新页面或重新登录后重试。'
          )}
        />
      ) : null}

      {hitsQuery.error && !authorizationChanged ? (
        <Alert
          type="error"
          showIcon
          message="违规记录加载失败"
          description={normalizeApiError(hitsQuery.error, '违规记录加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="筛选条件">
        <Space wrap>
          <Select
            style={{ width: 280 }}
            value={policyId}
            placeholder="按策略筛选"
            options={[
              { label: '全部策略', value: '' },
              ...(policiesQuery.data?.items || []).map((item) => ({
                label: item.name,
                value: item.id
              }))
            ]}
            onChange={setPolicyId}
          />
          <Select
            style={{ width: 180 }}
            value={riskLevel}
            options={riskLevelOptions}
            onChange={setRiskLevel}
          />
          <Select
            style={{ width: 200 }}
            value={remediationStatus}
            options={remediationOptions}
            onChange={setRemediationStatus}
          />
        </Space>
      </Card>

      <ViolationRiskChart hits={violations} />

      <Card
        size="small"
        title={`违规列表（${violations.length}）`}
        extra={selectedViolation ? `当前选中：${selectedViolation.id}` : '未选择'}
      >
        <ViolationTable
          violations={violations}
          loading={hitsQuery.isLoading || hitsQuery.isFetching}
          selectedViolationId={selectedViolation?.id}
          readonly={!canManage || authorizationChanged}
          onSelectViolation={setSelectedViolation}
          onUpdateRemediation={(violation) => {
            setSelectedViolation(violation);
            setDrawerOpen(true);
          }}
        />
      </Card>

      <RemediationUpdateDrawer
        open={drawerOpen}
        violation={selectedViolation}
        readonly={!canManage || authorizationChanged}
        onClose={() => setDrawerOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'violations'] });
        }}
      />
    </Space>
  );
};
