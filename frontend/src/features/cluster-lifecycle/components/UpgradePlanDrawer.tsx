import { Button, DatePicker, Drawer, Form, Input, Space } from 'antd';
import type { Dayjs } from 'dayjs';
import type { CreateUpgradePlanRequest } from '@/services/clusterLifecycle';

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateUpgradePlanRequest) => void;
};

type FormValues = {
  toVersion: string;
  reason?: string;
  window?: [Dayjs | null, Dayjs | null];
};

export const UpgradePlanDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm<FormValues>();

  return (
    <Drawer
      open={open}
      width={480}
      title="创建升级计划"
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={(values) =>
          onSubmit({
            toVersion: values.toVersion,
            reason: values.reason,
            windowStart: values.window?.[0]?.toISOString(),
            windowEnd: values.window?.[1]?.toISOString()
          })
        }
      >
        <Form.Item name="toVersion" label="目标版本" rules={[{ required: true, message: '请输入目标版本' }]}>
          <Input placeholder="例如：v1.31.2" />
        </Form.Item>
        <Form.Item name="window" label="升级窗口">
          <DatePicker.RangePicker showTime />
        </Form.Item>
        <Form.Item name="reason" label="变更说明">
          <Input.TextArea rows={3} placeholder="记录升级背景、影响范围和窗口说明。" />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            创建计划
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
