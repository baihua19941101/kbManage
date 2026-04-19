import { Drawer, Form, Input } from 'antd';

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: Record<string, unknown>) => void;
};

export const RollbackValidationDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm();
  return (
    <Drawer
      title="登记回退验证"
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
          {submitting ? '提交中...' : '提交'}
        </a>
      }
    >
      <Form form={form} layout="vertical" initialValues={{ validationScope: 'platform', result: 'passed' }}>
        <Form.Item name="validationScope" label="验证范围" rules={[{ required: true }]}>
          <Input placeholder="platform" />
        </Form.Item>
        <Form.Item name="result" label="验证结果" rules={[{ required: true }]}>
          <Input placeholder="passed" />
        </Form.Item>
        <Form.Item name="remainingRisk" label="剩余风险">
          <Input.TextArea rows={4} placeholder="记录尚未覆盖的风险点" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
