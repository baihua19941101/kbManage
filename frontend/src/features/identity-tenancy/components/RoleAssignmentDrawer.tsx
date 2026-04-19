import { Button, DatePicker, Drawer, Form, Input, Select, Space } from 'antd';
import type {
  CreateRoleAssignmentPayload,
  RoleDefinition
} from '@/services/identityTenancy';

type RoleAssignmentDrawerProps = {
  open: boolean;
  submitting: boolean;
  roleOptions: RoleDefinition[];
  onClose: () => void;
  onSubmit: (payload: CreateRoleAssignmentPayload) => void;
};

export const RoleAssignmentDrawer = ({
  open,
  submitting,
  roleOptions,
  onClose,
  onSubmit
}: RoleAssignmentDrawerProps) => {
  const [form] = Form.useForm();

  return (
    <Drawer
      title="新建授权"
      width={560}
      open={open}
      onClose={onClose}
      destroyOnClose
      extra={
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button
            type="primary"
            loading={submitting}
            onClick={() => {
              void form.validateFields().then((values) =>
                onSubmit({
                  subjectType: values.subjectType,
                  subjectRef: values.subjectRef?.trim(),
                  roleDefinitionId: values.roleDefinitionId,
                  scopeType: values.scopeType,
                  scopeRef: values.scopeRef?.trim(),
                  validUntil: values.validUntil?.toISOString()
                })
              );
            }}
          >
            保存
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item name="subjectType" label="主体类型" rules={[{ required: true, message: '请选择主体类型' }]}>
          <Select
            options={[
              { label: '用户', value: 'user' },
              { label: '团队', value: 'team' },
              { label: '用户组', value: 'group' }
            ]}
          />
        </Form.Item>
        <Form.Item name="subjectRef" label="主体标识" rules={[{ required: true, message: '请输入主体标识' }]}>
          <Input placeholder="例如：alice / team-platform / group-ops" />
        </Form.Item>
        <Form.Item
          name="roleDefinitionId"
          label="角色"
          rules={[{ required: true, message: '请选择角色' }]}
        >
          <Select options={roleOptions.map((role) => ({ label: role.name, value: role.id }))} />
        </Form.Item>
        <Form.Item name="scopeType" label="授权层级" rules={[{ required: true, message: '请选择授权层级' }]}>
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
        <Form.Item name="scopeRef" label="范围对象" rules={[{ required: true, message: '请输入范围对象' }]}>
          <Input placeholder="例如：platform / org-retail / project-orders" />
        </Form.Item>
        <Form.Item name="validUntil" label="到期时间">
          <DatePicker showTime style={{ width: '100%' }} />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
