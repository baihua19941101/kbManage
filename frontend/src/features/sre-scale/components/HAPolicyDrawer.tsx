import { Drawer, Form, Input, InputNumber } from 'antd';

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: Record<string, unknown>) => void;
};

export const HAPolicyDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm();
  return (
    <Drawer
      title="新增高可用策略"
      open={open}
      width={420}
      onClose={onClose}
      destroyOnClose
      extra={
        <a
          onClick={() => {
            void form.validateFields().then((values) => onSubmit(values));
          }}
        >
          {submitting ? '提交中...' : '保存'}
        </a>
      }
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ workspaceId: 1, deploymentMode: 'active-active', replicaExpectation: 3 }}
      >
        <Form.Item name="workspaceId" label="工作空间 ID" rules={[{ required: true }]}>
          <InputNumber style={{ width: '100%' }} min={1} />
        </Form.Item>
        <Form.Item name="name" label="策略名称" rules={[{ required: true }]}>
          <Input placeholder="platform-control-plane-ha" />
        </Form.Item>
        <Form.Item name="deploymentMode" label="部署模式" rules={[{ required: true }]}>
          <Input placeholder="active-active" />
        </Form.Item>
        <Form.Item name="replicaExpectation" label="期望副本数" rules={[{ required: true }]}>
          <InputNumber style={{ width: '100%' }} min={2} />
        </Form.Item>
        <Form.Item name="failoverTriggerPolicy" label="切换门槛">
          <Input.TextArea rows={4} placeholder="节点不可用 30s 且主依赖探测失败时切换" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
