import { useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { IdentitySourceDrawer } from '@/features/identity-tenancy/components/IdentitySourceDrawer';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  createIdentitySource,
  identitySourceQueryScope,
  identityTenancyQueryKeys,
  listIdentitySources,
  type IdentitySource
} from '@/services/identityTenancy';

const columns: ColumnsType<IdentitySource> = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '来源类型', dataIndex: 'sourceType', key: 'sourceType', render: (value?: string) => value || '—' },
  {
    title: '登录方式',
    dataIndex: 'loginMode',
    key: 'loginMode',
    render: (value?: string) => <StatusTag value={value} />
  },
  {
    title: '健康状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <StatusTag value={value} />
  },
  { title: '同步状态', dataIndex: 'syncState', key: 'syncState', render: (value?: string) => value || '—' }
];

export const IdentitySourcePage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useIdentityTenancyPermissions();
  const createMutation = useMutation({ mutationFn: createIdentitySource });
  const sourcesQuery = useQuery({
    queryKey: identityTenancyQueryKeys.sources(identitySourceQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listIdentitySources({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无身份源治理访问权限。" />;
  }

  const items = sourcesQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="身份源中心"
        description="维护 OIDC、LDAP、SSO 与本地账号并存策略，统一查看来源状态与登录方式。"
        actions={
          <Button type="primary" disabled={!permissions.canManageSource} onClick={() => setDrawerOpen(true)}>
            接入身份源
          </Button>
        }
      />

      <Alert
        type="info"
        showIcon
        message="本地管理员账号与外部身份源可并存"
        description="首期建议保留至少一个本地 break-glass 管理员账号，用于外部身份源异常时兜底。"
      />

      {sourcesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="身份源加载失败"
          description={normalizeApiError(sourcesQuery.error, '身份源加载失败，请稍后重试。')}
        />
      ) : null}
      {createMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="身份源创建失败"
          description={normalizeApiError(createMutation.error, '身份源创建失败，请稍后重试。')}
        />
      ) : null}
      {!permissions.canManageSource ? (
        <Alert type="warning" showIcon message="当前账号仅有查看权限，无法新增或修改身份源。" />
      ) : null}

      <Card size="small" title={`身份源列表（${items.length}）`}>
        <Table<IdentitySource>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={items}
          loading={sourcesQuery.isLoading || sourcesQuery.isFetching}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      <IdentitySourceDrawer
        open={drawerOpen}
        submitting={createMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
