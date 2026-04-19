import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { OrganizationUnitDrawer } from '@/features/identity-tenancy/components/OrganizationUnitDrawer';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useOrganizationAction } from '@/features/identity-tenancy/hooks/useOrganizationAction';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  identityTenancyQueryKeys,
  listOrganizationUnits,
  organizationQueryScope,
  type OrganizationUnit
} from '@/services/identityTenancy';

const columns: ColumnsType<OrganizationUnit> = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '类型', dataIndex: 'unitType', key: 'unitType', render: (value?: string) => value || '—' },
  { title: '上级单元', dataIndex: 'parentUnitId', key: 'parentUnitId', render: (value?: string) => value || '根节点' },
  { title: '成员数', dataIndex: 'memberCount', key: 'memberCount', render: (value?: number) => value ?? '—' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const OrganizationModelPage = () => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const permissions = useIdentityTenancyPermissions();
  const { createUnitMutation } = useOrganizationAction();
  const unitsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.organizations(organizationQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listOrganizationUnits({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无组织模型访问权限。" />;
  }

  const items = unitsQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="组织模型"
        description="维护组织、团队和用户组层级，沉淀成员来源、归属责任和租户边界前置语义。"
        actions={
          <Button type="primary" disabled={!permissions.canManageOrg} onClick={() => setDrawerOpen(true)}>
            新建组织单元
          </Button>
        }
      />

      <Card size="small">
        <Typography.Paragraph style={{ marginBottom: 0 }}>
          成员来源视图需要结合身份源同步结果与组织映射关系综合判断。当前页面优先呈现组织层级、父子关系和成员规模。
        </Typography.Paragraph>
      </Card>

      {unitsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="组织模型加载失败"
          description={normalizeApiError(unitsQuery.error, '组织模型加载失败，请稍后重试。')}
        />
      ) : null}
      {createUnitMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="组织单元创建失败"
          description={normalizeApiError(createUnitMutation.error, '组织单元创建失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`组织单元（${items.length}）`}>
        <Table<OrganizationUnit>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={items}
          loading={unitsQuery.isLoading || unitsQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <OrganizationUnitDrawer
        open={drawerOpen}
        submitting={createUnitMutation.isPending}
        units={items}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createUnitMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
