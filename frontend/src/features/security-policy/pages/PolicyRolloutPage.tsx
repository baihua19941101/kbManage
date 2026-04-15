import { useEffect, useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Card,
  Empty,
  Select,
  Space,
  Statistic,
  Table,
  Tag,
  Typography
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { queryKeys } from '@/app/queryClient';
import {
  canApprovePolicyException,
  canManagePolicy,
  canReadPolicy,
  useAuthStore
} from '@/features/auth/store';
import { ExceptionRequestDrawer } from '@/features/security-policy/components/ExceptionRequestDrawer';
import { ExceptionReviewDrawer } from '@/features/security-policy/components/ExceptionReviewDrawer';
import { ModeSwitchDrawer } from '@/features/security-policy/components/ModeSwitchDrawer';
import { usePolicyRollout } from '@/features/security-policy/hooks/usePolicyRollout';
import {
  isAuthorizationError,
  normalizeApiError,
  normalizeAuthorizationError
} from '@/services/api/client';
import { listSecurityPolicies } from '@/services/securityPolicy';
import type {
  PolicyExceptionRequestDTO,
  PolicyHitRecordDTO,
  SecurityPolicyDTO
} from '@/services/api/types';

const statusColorMap: Record<string, string> = {
  pending: 'gold',
  approved: 'blue',
  rejected: 'red',
  active: 'green',
  expired: 'default',
  revoked: 'volcano'
};

const hitResultColorMap: Record<string, string> = {
  pass: 'green',
  warn: 'orange',
  block: 'red'
};

const EMPTY_POLICIES: SecurityPolicyDTO[] = [];
const EMPTY_HITS: PolicyHitRecordDTO[] = [];
const EMPTY_EXCEPTIONS: PolicyExceptionRequestDTO[] = [];

const hitColumns: ColumnsType<PolicyHitRecordDTO> = [
  {
    title: '资源',
    key: 'resource',
    render: (_value, record) =>
      `${record.clusterId || '-'} / ${record.namespace || '-'} / ${record.resourceKind || '-'} / ${record.resourceName || '-'}`
  },
  {
    title: '命中结果',
    dataIndex: 'hitResult',
    key: 'hitResult',
    render: (value: string) => <Tag color={hitResultColorMap[value] || 'default'}>{value}</Tag>
  },
  {
    title: '风险级别',
    dataIndex: 'riskLevel',
    key: 'riskLevel'
  },
  {
    title: '整改状态',
    dataIndex: 'remediationStatus',
    key: 'remediationStatus'
  },
  {
    title: '发现时间',
    dataIndex: 'detectedAt',
    key: 'detectedAt'
  }
];

const exceptionColumns: ColumnsType<PolicyExceptionRequestDTO> = [
  {
    title: '申请 ID',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '命中 ID',
    dataIndex: 'hitRecordId',
    key: 'hitRecordId'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value: string) => <Tag color={statusColorMap[value] || 'default'}>{value}</Tag>
  },
  {
    title: '有效期',
    key: 'timeRange',
    render: (_value, record) => `${record.startsAt} ~ ${record.expiresAt}`
  },
  {
    title: '备注',
    dataIndex: 'reviewComment',
    key: 'reviewComment',
    render: (value?: string) => value || '-'
  }
];

export const PolicyRolloutPage = () => {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadPolicy(user);
  const canManage = canManagePolicy(user);
  const canApproveException = canApprovePolicyException(user);

  const [selectedPolicyId, setSelectedPolicyId] = useState<string>();
  const [modeDrawerOpen, setModeDrawerOpen] = useState(false);
  const [requestDrawerOpen, setRequestDrawerOpen] = useState(false);
  const [reviewDrawerOpen, setReviewDrawerOpen] = useState(false);
  const [selectedHit, setSelectedHit] = useState<PolicyHitRecordDTO>();
  const [selectedException, setSelectedException] = useState<PolicyExceptionRequestDTO>();

  const policiesQuery = useQuery({
    queryKey: queryKeys.securityPolicy.list('rollout-page'),
    enabled: canRead,
    queryFn: () => listSecurityPolicies({ status: 'active' })
  });

  const policies = policiesQuery.data?.items ?? EMPTY_POLICIES;
  useEffect(() => {
    if (!selectedPolicyId && policies.length > 0) {
      setSelectedPolicyId(policies[0]?.id);
    }
  }, [policies, selectedPolicyId]);
  const selectedPolicy = policies.find((item) => item.id === selectedPolicyId);
  const { hitsQuery, exceptionsQuery, permissionChanged } = usePolicyRollout(selectedPolicyId);

  const hits = hitsQuery.data?.items ?? EMPTY_HITS;
  const exceptions = exceptionsQuery.data?.items ?? EMPTY_EXCEPTIONS;

  const exceptionSummary = useMemo(() => {
    return exceptions.reduce(
      (acc, exception) => {
        if (exception.status in acc) {
          acc[exception.status as keyof typeof acc] += 1;
        }
        return acc;
      },
      { pending: 0, active: 0, expired: 0, revoked: 0 }
    );
  }, [exceptions]);

  const hasAuthError =
    permissionChanged || isAuthorizationError(policiesQuery.error) || isAuthorizationError(hitsQuery.error);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无策略灰度与例外治理访问权限，请联系管理员授予 policy:read 或对应平台角色。"
      />
    );
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          策略灰度与例外治理
        </Typography.Title>
        <Typography.Text type="secondary">
          支持模式切换、灰度验证、例外申请与审批，观察例外生命周期并跟踪状态变化。
        </Typography.Text>
      </div>

      {!canManage ? (
        <Alert
          type="info"
          showIcon
          message="当前为只读模式"
          description="你可查看命中、例外和状态统计，但无法提交切换、申请或审批。"
        />
      ) : null}

      {hasAuthError ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description={normalizeAuthorizationError(
            policiesQuery.error || hitsQuery.error || exceptionsQuery.error,
            '当前账号可能缺少灰度治理所需权限，请刷新页面或重新登录后重试。'
          )}
        />
      ) : null}

      {policiesQuery.error && !isAuthorizationError(policiesQuery.error) ? (
        <Alert
          type="error"
          showIcon
          message="策略列表加载失败"
          description={normalizeApiError(policiesQuery.error, '策略列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="选择策略">
        <Space wrap>
          <Select
            style={{ width: 320 }}
            placeholder="选择要执行灰度切换的策略"
            value={selectedPolicyId}
            options={policies.map((policy) => ({
              label: `${policy.name} (${policy.defaultEnforcementMode})`,
              value: policy.id
            }))}
            onChange={setSelectedPolicyId}
          />

          <Button
            type="primary"
            disabled={!selectedPolicy || !canManage || hasAuthError}
            onClick={() => setModeDrawerOpen(true)}
          >
            模式切换
          </Button>

          <Button
            disabled={!selectedHit || !canManage || hasAuthError}
            onClick={() => setRequestDrawerOpen(true)}
          >
            申请例外
          </Button>

          <Button
            disabled={!selectedException || !canApproveException || hasAuthError}
            onClick={() => setReviewDrawerOpen(true)}
          >
            审批例外
          </Button>
        </Space>
      </Card>

      <Card size="small" title="例外状态概览">
        <Space wrap size={24}>
          <Statistic title="待审批" value={exceptionSummary.pending} />
          <Statistic title="生效中" value={exceptionSummary.active} />
          <Statistic title="已过期" value={exceptionSummary.expired} />
          <Statistic title="已撤销" value={exceptionSummary.revoked} />
        </Space>
      </Card>

      {selectedPolicyId ? (
        <Card size="small" title={`策略命中（${hits.length}）`}>
          {hitsQuery.error && !isAuthorizationError(hitsQuery.error) ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(hitsQuery.error, '策略命中加载失败，请稍后重试。')}
            />
          ) : null}

          <Table<PolicyHitRecordDTO>
            rowKey={(record) => record.id}
            loading={hitsQuery.isLoading || hitsQuery.isFetching}
            columns={hitColumns}
            dataSource={hits}
            pagination={{ pageSize: 6 }}
            rowClassName={(record) => (record.id === selectedHit?.id ? 'ant-table-row-selected' : '')}
            onRow={(record) => ({
              onClick: () => setSelectedHit(record)
            })}
          />
        </Card>
      ) : null}

      {selectedPolicyId ? (
        <Card size="small" title={`例外申请（${exceptions.length}）`}>
          {exceptionsQuery.error && !isAuthorizationError(exceptionsQuery.error) ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(exceptionsQuery.error, '例外申请加载失败，请稍后重试。')}
            />
          ) : null}

          <Table<PolicyExceptionRequestDTO>
            rowKey={(record) => record.id}
            loading={exceptionsQuery.isLoading || exceptionsQuery.isFetching}
            columns={exceptionColumns}
            dataSource={exceptions}
            pagination={{ pageSize: 6 }}
            rowClassName={(record) =>
              record.id === selectedException?.id ? 'ant-table-row-selected' : ''
            }
            onRow={(record) => ({
              onClick: () => setSelectedException(record)
            })}
          />

          {exceptions.some((item) => item.status === 'expired') ? (
            <Alert
              type="warning"
              showIcon
              style={{ marginTop: 12 }}
              message="检测到已过期例外"
              description="已过期例外不再放行，请及时确认工作负载是否已完成整改。"
            />
          ) : null}
        </Card>
      ) : null}

      <ModeSwitchDrawer
        open={modeDrawerOpen}
        policy={selectedPolicy}
        readonly={!canManage || hasAuthError}
        onClose={() => setModeDrawerOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'hits', selectedPolicyId] });
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'exceptions', selectedPolicyId] });
        }}
      />

      <ExceptionRequestDrawer
        open={requestDrawerOpen}
        hit={selectedHit}
        readonly={!canManage || hasAuthError}
        onClose={() => setRequestDrawerOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'exceptions', selectedPolicyId] });
        }}
      />

      <ExceptionReviewDrawer
        open={reviewDrawerOpen}
        exception={selectedException}
        readonly={!canApproveException || hasAuthError}
        onClose={() => setReviewDrawerOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'exceptions', selectedPolicyId] });
          void queryClient.invalidateQueries({ queryKey: ['securityPolicy', 'hits', selectedPolicyId] });
        }}
      />
    </Space>
  );
};
