import { Drawer, Form, Input, InputNumber } from 'antd';

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: Record<string, unknown>) => void;
};

export const UpgradePlanDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm();
  return (
    <Drawer
      title="创建升级计划"
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
          {submitting ? '提交中...' : '创建'}
        </a>
      }
    >
      <Form form={form} layout="vertical" initialValues={{ workspaceId: 1, rolloutStrategy: 'rolling' }}>
        <Form.Item name="workspaceId" label="工作空间 ID" rules={[{ required: true }]}>
          <InputNumber style={{ width: '100%' }} min={1} />
        </Form.Item>
        <Form.Item name="name" label="计划名称" rules={[{ required: true }]}>
          <Input placeholder="platform-1.31-upgrade" />
        </Form.Item>
        <Form.Item name="currentVersion" label="当前版本" rules={[{ required: true }]}>
          <Input placeholder="1.30.0" />
        </Form.Item>
        <Form.Item name="targetVersion" label="目标版本" rules={[{ required: true }]}>
          <Input placeholder="1.31.0" />
        </Form.Item>
        <Form.Item name="rolloutStrategy" label="升级策略" rules={[{ required: true }]}>
          <Input placeholder="rolling" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
