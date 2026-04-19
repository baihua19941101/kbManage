import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DelegationGrantDrawer } from '@/features/identity-tenancy/components/DelegationGrantDrawer';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { useRoleGovernanceAction } from '@/features/identity-tenancy/hooks/useRoleGovernanceAction';
import { normalizeApiError } from '@/services/api/client';
import {
  identityTenancyQueryKeys,
  listDelegationGrants,
  type DelegationGrant
} from '@/services/identityTenancy';

const columns: ColumnsType<DelegationGrant> = [
  { title: '委派人', dataIndex: 'grantorRef', key: 'grantorRef', render: (value?: string) => value || '—' },
  { title: '被委派人', dataIndex: 'delegateRef', key: 'delegateRef', render: (value?: string) => value || '—' },
  { title: '允许层级', key: 'allowedRoleLevels', render: (_, record) => record.allowedRoleLevels.join(', ') || '—' },
  { title: '有效期', key: 'validity', render: (_, record) => `${record.validFrom || '—'} ~ ${record.validUntil || '—'}` },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const DelegationPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useIdentityTenancyPermissions();
  const { createDelegationMutation } = useRoleGovernanceAction();
  const delegationsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.delegations(),
    enabled: permissions.canRead,
    queryFn: () => listDelegationGrants()
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无委派治理访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="委派与临时授权"
        description="维护委派链路、允许层级、到期时间和回收理由。"
        actions={
          <Button type="primary" disabled={!permissions.canDelegate} onClick={() => setDrawerOpen(true)}>
            新建委派
          </Button>
        }
      />

      {delegationsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="委派记录加载失败"
          description={normalizeApiError(delegationsQuery.error, '委派记录加载失败，请稍后重试。')}
        />
      ) : null}
      {createDelegationMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="委派创建失败"
          description={normalizeApiError(createDelegationMutation.error, '委派创建失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`委派记录（${delegationsQuery.data?.items.length ?? 0}）`}>
        <Table<DelegationGrant>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={delegationsQuery.data?.items || []}
          loading={delegationsQuery.isLoading || delegationsQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <DelegationGrantDrawer
        open={drawerOpen}
        submitting={createDelegationMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createDelegationMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
