import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateTenantScopeMappingPayload } from '@/services/identityTenancy';

type TenantScopeDrawerProps = {
  open: boolean;
  submitting: boolean;
  unitName?: string;
  onClose: () => void;
  onSubmit: (payload: CreateTenantScopeMappingPayload) => void;
};

export const TenantScopeDrawer = ({
  open,
  submitting,
  unitName,
  onClose,
  onSubmit
}: TenantScopeDrawerProps) => {
  const [form] = Form.useForm<CreateTenantScopeMappingPayload>();

  return (
    <Drawer
      title={`新增租户边界映射${unitName ? ` · ${unitName}` : ''}`}
      width={520}
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
              void form.validateFields().then((values) => onSubmit(values));
            }}
          >
            保存
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item name="scopeType" label="范围类型" rules={[{ required: true, message: '请选择范围类型' }]}>
          <Select
            options={[
              { label: '工作空间', value: 'workspace' },
              { label: '项目', value: 'project' },
              { label: '资源', value: 'resource' }
            ]}
          />
        </Form.Item>
        <Form.Item name="scopeRef" label="范围引用" rules={[{ required: true, message: '请输入范围引用' }]}>
          <Input placeholder="例如：workspace-retail / project-payment" />
        </Form.Item>
        <Form.Item
          name="inheritanceMode"
          label="继承模式"
          rules={[{ required: true, message: '请选择继承模式' }]}
        >
          <Select
            options={[
              { label: '严格继承', value: 'strict' },
              { label: '可扩展继承', value: 'extendable' },
              { label: '不继承', value: 'isolated' }
            ]}
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
