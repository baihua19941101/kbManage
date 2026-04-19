import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { RoleAssignmentDrawer } from '@/features/identity-tenancy/components/RoleAssignmentDrawer';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { useRoleGovernanceAction } from '@/features/identity-tenancy/hooks/useRoleGovernanceAction';
import { normalizeApiError } from '@/services/api/client';
import {
  identityTenancyQueryKeys,
  listRoleAssignments,
  listRoleDefinitions,
  roleAssignmentQueryScope,
  summarizeRoleBoundary,
  type RoleAssignment
} from '@/services/identityTenancy';

const columns: ColumnsType<RoleAssignment> = [
  { title: '主体', key: 'subject', render: (_, record) => `${record.subjectType || '—'} / ${record.subjectRef || '—'}` },
  { title: '角色', key: 'role', render: (_, record) => record.roleDefinitionName || record.roleDefinitionId || '—' },
  { title: '权限边界', key: 'boundary', render: (_, record) => summarizeRoleBoundary(record) },
  { title: '来源', dataIndex: 'sourceType', key: 'sourceType', render: (value?: string) => value || 'direct' },
  {
    title: '到期状态',
    key: 'expiry',
    render: (_, record) => (record.validUntil ? <Tag color="gold">{record.validUntil}</Tag> : '长期有效')
  },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const RoleAssignmentPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useIdentityTenancyPermissions();
  const { createAssignmentMutation } = useRoleGovernanceAction();
  const assignmentsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.assignments(roleAssignmentQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listRoleAssignments({})
  });
  const rolesQuery = useQuery({
    queryKey: identityTenancyQueryKeys.roles('assignment-form'),
    enabled: permissions.canRead,
    queryFn: () => listRoleDefinitions({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无授权分配访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="授权分配"
        description="查看主体、角色、范围边界和临时授权到期状态，统一识别跨租户扩散风险。"
        actions={
          <Button type="primary" disabled={!permissions.canManageRole} onClick={() => setDrawerOpen(true)}>
            新建授权
          </Button>
        }
      />

      {assignmentsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="授权分配加载失败"
          description={normalizeApiError(assignmentsQuery.error, '授权分配加载失败，请稍后重试。')}
        />
      ) : null}
      {createAssignmentMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="授权创建失败"
          description={normalizeApiError(createAssignmentMutation.error, '授权创建失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`授权列表（${assignmentsQuery.data?.items.length ?? 0}）`}>
        <Table<RoleAssignment>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={assignmentsQuery.data?.items || []}
          loading={assignmentsQuery.isLoading || assignmentsQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <RoleAssignmentDrawer
        open={drawerOpen}
        submitting={createAssignmentMutation.isPending}
        roleOptions={rolesQuery.data?.items || []}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createAssignmentMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
