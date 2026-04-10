import { Button, Form, Input, Select, Space } from 'antd';

export type RoleBindingPayload = {
  user: string;
  role: string;
  scope: string;
};

type RoleBindingFormProps = {
  loading?: boolean;
  onSubmit?: (payload: RoleBindingPayload) => void;
};

const roleOptions = [
  { label: '平台管理员', value: 'platform-admin' },
  { label: '工作空间管理员', value: 'workspace-admin' },
  { label: '开发者', value: 'developer' },
  { label: '只读用户', value: 'viewer' }
];

export const RoleBindingForm = ({ loading = false, onSubmit }: RoleBindingFormProps) => {
  const [form] = Form.useForm<RoleBindingPayload>();

  const handleFinish = (values: RoleBindingPayload) => {
    const payload: RoleBindingPayload = {
      user: values.user.trim(),
      role: values.role,
      scope: values.scope.trim()
    };

    onSubmit?.(payload);
    form.resetFields();
  };

  return (
    <Form<RoleBindingPayload> form={form} layout="vertical" onFinish={handleFinish}>
      <Form.Item
        label="用户"
        name="user"
        rules={[{ required: true, message: '请输入用户标识' }]}
      >
        <Input placeholder="例如：alice" maxLength={64} />
      </Form.Item>

      <Form.Item
        label="角色"
        name="role"
        rules={[{ required: true, message: '请选择角色' }]}
      >
        <Select options={roleOptions} placeholder="选择角色" />
      </Form.Item>

      <Form.Item
        label="授权范围"
        name="scope"
        rules={[{ required: true, message: '请输入 scope' }]}
      >
        <Input placeholder="例如：workspace:dev-team 或 project:billing-api" />
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
