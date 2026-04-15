import { useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Statistic, Typography } from 'antd';
import { queryKeys } from '@/app/queryClient';
import {
  canAssignPolicy,
  canManagePolicy,
  canReadPolicy,
  useAuthStore
} from '@/features/auth/store';
import { PolicyEditorDrawer } from '@/features/security-policy/components/PolicyEditorDrawer';
import { PolicyScopeDrawer } from '@/features/security-policy/components/PolicyScopeDrawer';
import { PolicyTable } from '@/features/security-policy/components/PolicyTable';
import {
  isAuthorizationError,
  normalizeApiError
} from '@/services/api/client';
import {
  listSecurityPolicies,
  type SecurityPolicyListQuery
} from '@/services/securityPolicy';
import type { SecurityPolicyDTO } from '@/services/api/types';

const scopeLevelOptions = [
  { label: '全部层级', value: '' },
  { label: 'platform', value: 'platform' },
  { label: 'workspace', value: 'workspace' },
  { label: 'project', value: 'project' }
];

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: 'draft', value: 'draft' },
  { label: 'active', value: 'active' },
  { label: 'disabled', value: 'disabled' },
  { label: 'archived', value: 'archived' }
];

const categoryOptions = [
  { label: '全部类别', value: '' },
  { label: 'pod-security', value: 'pod-security' },
  { label: 'image', value: 'image' },
  { label: 'resource', value: 'resource' },
  { label: 'label', value: 'label' },
  { label: 'network', value: 'network' },
  { label: 'admission', value: 'admission' }
];

const EMPTY_POLICIES: SecurityPolicyDTO[] = [];

export const PolicyCenterPage = () => {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadPolicy(user);
  const canManage = canManagePolicy(user);
  const canAssign = canAssignPolicy(user);
  const [scopeLevel, setScopeLevel] = useState<string>('');
  const [status, setStatus] = useState<string>('');
  const [category, setCategory] = useState<string>('');
  const [editorOpen, setEditorOpen] = useState(false);
  const [scopeDrawerOpen, setScopeDrawerOpen] = useState(false);
  const [editingPolicy, setEditingPolicy] = useState<SecurityPolicyDTO>();
  const [selectedPolicy, setSelectedPolicy] = useState<SecurityPolicyDTO>();

  const listQuery = useMemo<SecurityPolicyListQuery>(
    () => ({
      scopeLevel: scopeLevel ? (scopeLevel as SecurityPolicyListQuery['scopeLevel']) : undefined,
      status: status ? (status as SecurityPolicyListQuery['status']) : undefined,
      category: category ? (category as SecurityPolicyListQuery['category']) : undefined
    }),
    [scopeLevel, status, category]
  );

  const scopeKey = `${scopeLevel || 'all'}:${status || 'all'}:${category || 'all'}`;

  const policiesQuery = useQuery({
    queryKey: queryKeys.securityPolicy.list(scopeKey),
    enabled: canRead,
    queryFn: () => listSecurityPolicies(listQuery)
  });

  const policies = policiesQuery.data?.items ?? EMPTY_POLICIES;
  const permissionChanged = isAuthorizationError(policiesQuery.error);

  const summary = useMemo(() => {
    return policies.reduce(
      (acc, policy) => {
        acc.total += 1;
        if (policy.scopeLevel === 'platform') {
          acc.platform += 1;
        }
        if (policy.scopeLevel === 'workspace') {
          acc.workspace += 1;
        }
        if (policy.scopeLevel === 'project') {
          acc.project += 1;
        }
        if (policy.status === 'active') {
          acc.applicable += 1;
        }
        return acc;
      },
      { total: 0, platform: 0, workspace: 0, project: 0, applicable: 0 }
    );
  }, [policies]);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无安全策略访问权限，请联系管理员授予 policy:read 或对应平台角色。"
      />
    );
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          安全策略中心
        </Typography.Title>
        <Typography.Text type="secondary">
          统一维护平台级、工作空间级、项目级策略，完成策略创建、编辑和范围分配。
        </Typography.Text>
      </div>

      {!canManage ? (
        <Alert
          type="info"
          showIcon
          message="当前为只读模式"
          description="你可查看策略与分配结果，但无法创建或编辑策略。"
        />
      ) : null}

      {permissionChanged ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description="当前账号可能缺少策略中心所需权限，请刷新页面或重新登录后重试。"
        />
      ) : null}

      {policiesQuery.error && !permissionChanged ? (
        <Alert
          type="error"
          showIcon
          message="策略列表加载失败"
          description={normalizeApiError(policiesQuery.error, '策略列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Space wrap>
        <Button
          type="primary"
          disabled={!canManage || permissionChanged}
          onClick={() => {
            setEditingPolicy(undefined);
            setEditorOpen(true);
          }}
        >
          新建策略
        </Button>
        <Button
          disabled={!selectedPolicy || !canAssign || permissionChanged}
          onClick={() => setScopeDrawerOpen(true)}
        >
          分配策略
        </Button>
      </Space>

      <Card size="small" title="筛选条件">
        <Space wrap>
          <Select
            style={{ width: 180 }}
            value={scopeLevel}
            options={scopeLevelOptions}
            onChange={setScopeLevel}
          />
          <Select
            style={{ width: 180 }}
            value={status}
            options={statusOptions}
            onChange={setStatus}
          />
          <Select
            style={{ width: 200 }}
            value={category}
            options={categoryOptions}
            onChange={setCategory}
          />
        </Space>
      </Card>

      <Card size="small" title="策略层级与最终适用集合">
        <Space wrap size={24}>
          <Statistic title="策略总数" value={summary.total} />
          <Statistic title="平台级" value={summary.platform} />
          <Statistic title="工作空间级" value={summary.workspace} />
          <Statistic title="项目级" value={summary.project} />
          <Statistic title="当前最终适用(Active)" value={summary.applicable} />
        </Space>
      </Card>

      <Card
        size="small"
        title={`策略列表（${policies.length}）`}
        extra={
          selectedPolicy ? (
            <Typography.Text type="secondary">当前选中：{selectedPolicy.name}</Typography.Text>
          ) : null
        }
      >
        <PolicyTable
          loading={policiesQuery.isLoading || policiesQuery.isFetching}
          policies={policies}
          selectedPolicyId={selectedPolicy?.id}
          readonly={!canManage || permissionChanged}
          onSelectPolicy={setSelectedPolicy}
          onEditPolicy={(policy) => {
            setEditingPolicy(policy);
            setEditorOpen(true);
          }}
          onAssignPolicy={(policy) => {
            setSelectedPolicy(policy);
            setScopeDrawerOpen(true);
          }}
        />
      </Card>

      <PolicyEditorDrawer
        open={editorOpen}
        policy={editingPolicy}
        onClose={() => setEditorOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: queryKeys.securityPolicy.list() });
        }}
      />

      <PolicyScopeDrawer
        open={scopeDrawerOpen}
        policy={selectedPolicy}
        readonly={!canAssign || permissionChanged}
        onClose={() => setScopeDrawerOpen(false)}
        onSuccess={() => {
          void queryClient.invalidateQueries({ queryKey: queryKeys.securityPolicy.list() });
        }}
      />
    </Space>
  );
};
