import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateTemplateRequest } from '@/services/clusterLifecycle';

type Props = {
  open: boolean;
  driverOptions: Array<{ label: string; value: string }>;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateTemplateRequest) => void;
};

export const ClusterTemplateFormDrawer = ({
  open,
  driverOptions,
  submitting,
  onClose,
  onSubmit
}: Props) => {
  const [form] = Form.useForm<CreateTemplateRequest>();

  return (
    <Drawer
      open={open}
      width={520}
      title="新建集群模板"
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Form form={form} layout="vertical" onFinish={onSubmit}>
        <Form.Item name="name" label="模板名称" rules={[{ required: true, message: '请输入模板名称' }]}>
          <Input placeholder="prod-standard-ha" />
        </Form.Item>
        <Form.Item
          name="infrastructureType"
          label="基础设施类型"
          rules={[{ required: true, message: '请输入基础设施类型' }]}
        >
          <Input placeholder="cloud / virtualized / baremetal" />
        </Form.Item>
        <Form.Item name="driverKey" label="驱动键" rules={[{ required: true, message: '请选择驱动' }]}>
          <Select options={driverOptions} showSearch />
        </Form.Item>
        <Form.Item name="driverVersionRange" label="驱动版本范围">
          <Input placeholder=">=1.0.0 <2.0.0" />
        </Form.Item>
        <Form.Item name="requiredCapabilities" label="必需能力">
          <Select
            mode="tags"
            placeholder="例如：network, security, observability"
            tokenSeparators={[',']}
          />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            保存模板
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
