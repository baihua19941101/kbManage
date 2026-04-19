import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Form, Input, Select, Space, Switch, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  createRoleDefinition,
  identityTenancyQueryKeys,
  listRoleDefinitions,
  type CreateRoleDefinitionPayload,
  type RoleDefinition
} from '@/services/identityTenancy';

const columns: ColumnsType<RoleDefinition> = [
  { title: '角色名称', dataIndex: 'name', key: 'name' },
  { title: '层级', dataIndex: 'roleLevel', key: 'roleLevel', render: (value?: string) => value || '—' },
  { title: '权限摘要', dataIndex: 'permissionSummary', key: 'permissionSummary', render: (value?: string) => value || '—' },
  { title: '继承策略', dataIndex: 'inheritancePolicy', key: 'inheritancePolicy', render: (value?: string) => value || '—' },
  { title: '可委派', dataIndex: 'delegable', key: 'delegable', render: (value?: boolean) => (value ? '是' : '否') },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const RoleCatalogPage = () => {
  const permissions = useIdentityTenancyPermissions();
  const [form] = Form.useForm<CreateRoleDefinitionPayload>();
  const createMutation = useMutation({ mutationFn: createRoleDefinition });
  const rolesQuery = useQuery({
    queryKey: identityTenancyQueryKeys.roles(),
    enabled: permissions.canRead,
    queryFn: () => listRoleDefinitions({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无角色目录访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="角色目录"
        description="查看平台级到资源级角色定义、继承策略和是否允许继续委派。"
      />

      {rolesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="角色目录加载失败"
          description={normalizeApiError(rolesQuery.error, '角色目录加载失败，请稍后重试。')}
        />
      ) : null}
      {createMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="角色创建失败"
          description={normalizeApiError(createMutation.error, '角色创建失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="角色定义录入">
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="角色名称" rules={[{ required: true, message: '请输入角色名称' }]}>
            <Input placeholder="例如：组织安全管理员" />
          </Form.Item>
          <Form.Item name="roleLevel" label="角色层级" rules={[{ required: true, message: '请选择角色层级' }]}>
            <Select
              options={[
                { label: '平台级', value: 'platform' },
                { label: '组织级', value: 'organization' },
                { label: '工作空间级', value: 'workspace' },
                { label: '项目级', value: 'project' },
                { label: '资源级', value: 'resource' }
              ]}
            />
          </Form.Item>
          <Form.Item
            name="permissionSummary"
            label="权限摘要"
            rules={[{ required: true, message: '请输入权限摘要' }]}
          >
            <Input placeholder="例如：可查看会话、管理授权、回收风险访问" />
          </Form.Item>
          <Form.Item
            name="inheritancePolicy"
            label="继承策略"
            rules={[{ required: true, message: '请选择继承策略' }]}
          >
            <Select
              options={[
                { label: '严格继承', value: 'strict' },
                { label: '可收敛继承', value: 'bounded' },
                { label: '不继承', value: 'isolated' }
              ]}
            />
          </Form.Item>
          <Form.Item name="description" label="说明">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="delegable" label="允许委派" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Button
            type="primary"
            disabled={!permissions.canManageRole}
            loading={createMutation.isPending}
            onClick={() => {
              void form.validateFields().then((values) => createMutation.mutate(values));
            }}
          >
            保存角色定义
          </Button>
        </Form>
      </Card>

      <Card size="small" title={`角色定义（${rolesQuery.data?.items.length ?? 0}）`}>
        <Table<RoleDefinition>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={rolesQuery.data?.items || []}
          loading={rolesQuery.isLoading || rolesQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>
    </Space>
  );
};
