import { Button, Form, Input, Select, Space } from 'antd';

export type RoleBindingPayload = {
  subjectType: 'user' | 'group';
  subjectId: string;
  scopeType: 'workspace' | 'project';
  scopeId: string;
  roleKey: string;
};

type RoleBindingFormProps = {
  loading?: boolean;
  onSubmit?: (payload: RoleBindingPayload) => void;
};

const roleOptions = [
  { label: 'Workspace Owner', value: 'workspace-owner' },
  { label: 'Workspace Viewer', value: 'workspace-viewer' },
  { label: 'Project Owner', value: 'project-owner' },
  { label: 'Project Viewer', value: 'project-viewer' }
];

export const RoleBindingForm = ({ loading = false, onSubmit }: RoleBindingFormProps) => {
  const [form] = Form.useForm<RoleBindingPayload>();

  const handleFinish = (values: RoleBindingPayload) => {
    const payload: RoleBindingPayload = {
      subjectType: values.subjectType,
      subjectId: values.subjectId.trim(),
      scopeType: values.scopeType,
      scopeId: values.scopeId.trim(),
      roleKey: values.roleKey
    };

    onSubmit?.(payload);
    form.resetFields();
  };

  return (
    <Form<RoleBindingPayload> form={form} layout="vertical" onFinish={handleFinish}>
      <Form.Item
        label="主体类型"
        name="subjectType"
        initialValue="user"
        rules={[{ required: true, message: '请选择主体类型' }]}
      >
        <Select
          options={[
            { label: '用户', value: 'user' },
            { label: '用户组', value: 'group' }
          ]}
        />
      </Form.Item>

      <Form.Item
        label="主体 ID"
        name="subjectId"
        rules={[{ required: true, message: '请输入用户标识' }]}
      >
        <Input placeholder="例如：1001" maxLength={64} />
      </Form.Item>

      <Form.Item
        label="角色"
        name="roleKey"
        rules={[{ required: true, message: '请选择角色' }]}
      >
        <Select options={roleOptions} placeholder="选择角色" />
      </Form.Item>

      <Form.Item
        label="范围类型"
        name="scopeType"
        initialValue="workspace"
        rules={[{ required: true, message: '请选择范围类型' }]}
      >
        <Select
          options={[
            { label: '工作空间', value: 'workspace' },
            { label: '项目', value: 'project' }
          ]}
        />
      </Form.Item>

      <Form.Item
        label="范围 ID"
        name="scopeId"
        rules={[{ required: true, message: '请输入 scope' }]}
      >
        <Input placeholder="例如：2001" />
      </Form.Item>

      <Space>
        <Button type="primary" htmlType="submit" loading={loading}>
          提交绑定
        </Button>
        <Button htmlType="button" onClick={() => form.resetFields()}>
          重置
        </Button>
      </Space>
    </Form>
  );
};
