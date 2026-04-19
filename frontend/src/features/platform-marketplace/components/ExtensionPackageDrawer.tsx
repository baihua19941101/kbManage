import { Drawer, Form, Input } from 'antd';
import type { CreateExtensionPackagePayload } from '@/services/platformMarketplace';

type ExtensionPackageDrawerProps = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateExtensionPackagePayload) => void;
};

export const ExtensionPackageDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: ExtensionPackageDrawerProps) => {
  const [form] = Form.useForm<CreateExtensionPackagePayload>();

  return (
    <Drawer
      title="注册扩展"
      open={open}
      width={460}
      onClose={onClose}
      destroyOnClose
      extra={
        <a
          onClick={() => {
            void form.validateFields().then((values) => onSubmit(values));
          }}
        >
          {submitting ? '提交中...' : '注册'}
        </a>
      }
    >
      <Form layout="vertical" form={form} initialValues={{ visibilityScope: 'platform-admin' }}>
        <Form.Item name="name" label="扩展名称" rules={[{ required: true, message: '请输入扩展名称' }]}>
          <Input placeholder="service-mesh-observer" />
        </Form.Item>
        <Form.Item name="version" label="版本" rules={[{ required: true, message: '请输入扩展版本' }]}>
          <Input placeholder="2.1.0" />
        </Form.Item>
        <Form.Item
          name="visibilityScope"
          label="可见范围"
          rules={[{ required: true, message: '请输入可见范围' }]}
        >
          <Input placeholder="workspace/platform" />
        </Form.Item>
        <Form.Item
          name="permissionSummary"
          label="权限声明"
          rules={[{ required: true, message: '请输入权限声明' }]}
        >
          <Input.TextArea rows={4} placeholder="说明扩展需要的平台权限、目标资源访问范围和影响。" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
