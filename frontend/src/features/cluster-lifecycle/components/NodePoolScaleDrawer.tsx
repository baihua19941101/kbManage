import { Button, Drawer, Form, Input, InputNumber, Space } from 'antd';
import type { ScaleNodePoolRequest } from '@/services/clusterLifecycle';

type Props = {
  open: boolean;
  submitting?: boolean;
  initialDesiredCount?: number;
  onClose: () => void;
  onSubmit: (payload: ScaleNodePoolRequest) => void;
};

export const NodePoolScaleDrawer = ({
  open,
  submitting,
  initialDesiredCount,
  onClose,
  onSubmit
}: Props) => {
  const [form] = Form.useForm<ScaleNodePoolRequest>();

  return (
    <Drawer
      open={open}
      width={420}
      title="调整节点池容量"
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ desiredCount: initialDesiredCount }}
        onFinish={onSubmit}
      >
        <Form.Item
          name="desiredCount"
          label="目标节点数"
          rules={[{ required: true, message: '请输入目标节点数' }]}
        >
          <InputNumber style={{ width: '100%' }} min={0} />
        </Form.Item>
        <Form.Item name="reason" label="调整原因">
          <Input placeholder="例如：业务扩容窗口" />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            提交调整
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
