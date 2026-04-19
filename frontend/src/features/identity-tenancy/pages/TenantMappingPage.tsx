import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { TenantScopeDrawer } from '@/features/identity-tenancy/components/TenantScopeDrawer';
import { useOrganizationAction } from '@/features/identity-tenancy/hooks/useOrganizationAction';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  identityTenancyQueryKeys,
  listOrganizationUnits,
  listTenantScopeMappings,
  summarizeScopeBoundary,
  type OrganizationUnit,
  type TenantScopeMapping
} from '@/services/identityTenancy';

const columns: ColumnsType<TenantScopeMapping> = [
  { title: '边界', key: 'boundary', render: (_, record) => summarizeScopeBoundary(record.scopeType, record.scopeRef) },
  { title: '继承模式', dataIndex: 'inheritanceMode', key: 'inheritanceMode', render: (value?: string) => value || '—' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> },
  { title: '冲突提示', dataIndex: 'conflictSummary', key: 'conflictSummary', render: (value?: string) => value || '—' }
];

export const TenantMappingPage = () => {
  const permissions = useIdentityTenancyPermissions();
  const [selectedUnitId, setSelectedUnitId] = useState<string>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { createMappingMutation } = useOrganizationAction(selectedUnitId);

  const unitsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.organizations('mapping-source'),
    enabled: permissions.canRead,
    queryFn: () => listOrganizationUnits({})
  });
  const mappingsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.mappings(selectedUnitId),
    enabled: permissions.canRead && Boolean(selectedUnitId),
    queryFn: () => listTenantScopeMappings(selectedUnitId as string)
  });

  const selectedUnit = useMemo(
    () => (unitsQuery.data?.items || []).find((unit) => unit.id === selectedUnitId),
    [selectedUnitId, unitsQuery.data?.items]
  );

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无租户边界映射访问权限。" />;
  }

  const mappings = mappingsQuery.data?.items || [];
  const hasConflict = mappings.some((item) => item.status === 'conflict' || Boolean(item.conflictSummary));

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="租户边界映射"
        description="把组织、团队和用户组映射到工作空间、项目和资源级边界，识别潜在冲突。"
        actions={
          <Button
            type="primary"
            disabled={!permissions.canManageOrg || !selectedUnitId}
            onClick={() => setDrawerOpen(true)}
          >
            新建边界映射
          </Button>
        }
      />

      <Card size="small">
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          <Typography.Text strong>选择组织单元</Typography.Text>
          <Select
            allowClear
            placeholder="选择需要查看边界的组织单元"
            value={selectedUnitId}
            onChange={(value) => setSelectedUnitId(value)}
            options={(unitsQuery.data?.items || []).map((unit: OrganizationUnit) => ({
              label: `${unit.name} (${unit.unitType || '未分类'})`,
              value: unit.id
            }))}
          />
        </Space>
      </Card>

      {unitsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="组织单元加载失败"
          description={normalizeApiError(unitsQuery.error, '组织单元加载失败，请稍后重试。')}
        />
      ) : null}
      {mappingsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="边界映射加载失败"
          description={normalizeApiError(mappingsQuery.error, '边界映射加载失败，请稍后重试。')}
        />
      ) : null}
      {createMappingMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="边界映射创建失败"
          description={normalizeApiError(createMappingMutation.error, '边界映射创建失败，请稍后重试。')}
        />
      ) : null}
      {hasConflict ? (
        <Alert
          type="warning"
          showIcon
          message="发现租户边界冲突"
          description="至少一个映射存在范围冲突或继承模式歧义，正式生效前需要管理员进一步处理。"
        />
      ) : null}

      {!selectedUnitId ? (
        <Card size="small">
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="先选择组织单元，再查看租户边界。" />
        </Card>
      ) : (
        <Card size="small" title={`边界映射（${mappings.length}）`}>
          <Table<TenantScopeMapping>
            rowKey={(record) => record.id}
            columns={columns}
            dataSource={mappings}
            loading={mappingsQuery.isLoading || mappingsQuery.isFetching}
            pagination={{ pageSize: 8 }}
          />
        </Card>
      )}

      <TenantScopeDrawer
        open={drawerOpen}
        submitting={createMappingMutation.isPending}
        unitName={selectedUnit?.name}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          createMappingMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
