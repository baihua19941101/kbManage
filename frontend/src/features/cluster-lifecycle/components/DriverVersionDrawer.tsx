import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateDriverRequest } from '@/services/clusterLifecycle';

const providerOptions = [
  { label: '公有云', value: 'cloud' },
  { label: '虚拟化', value: 'virtualized' },
  { label: '裸金属', value: 'baremetal' },
  { label: '托管 K8s', value: 'managed-kubernetes' }
];

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateDriverRequest) => void;
};

export const DriverVersionDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm<CreateDriverRequest>();

  return (
    <Drawer
      open={open}
      width={480}
      title="登记驱动版本"
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Form form={form} layout="vertical" onFinish={onSubmit}>
        <Form.Item name="displayName" label="展示名称">
          <Input placeholder="例如：RKE2 on vSphere" />
        </Form.Item>
        <Form.Item name="driverKey" label="驱动键" rules={[{ required: true, message: '请输入驱动键' }]}>
          <Input placeholder="rke2-vsphere" />
        </Form.Item>
        <Form.Item name="version" label="驱动版本" rules={[{ required: true, message: '请输入版本' }]}>
          <Input placeholder="1.0.0" />
        </Form.Item>
        <Form.Item
          name="providerType"
          label="基础设施类型"
          rules={[{ required: true, message: '请选择基础设施类型' }]}
        >
          <Select options={providerOptions} />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            保存驱动版本
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
