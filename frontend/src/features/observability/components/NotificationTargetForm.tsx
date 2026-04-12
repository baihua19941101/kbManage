import { Button, Form, Input, Select } from 'antd';
import type { NotificationTargetPayload } from '@/services/observability/notificationTargets';

type NotificationTargetFormProps = {
  loading?: boolean;
  onSubmit: (payload: NotificationTargetPayload) => void | Promise<void>;
};

export const NotificationTargetForm = ({ loading, onSubmit }: NotificationTargetFormProps) => {
  const [form] = Form.useForm<NotificationTargetPayload>();

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={{ targetType: 'webhook', status: 'active' }}
      onFinish={(values) => void onSubmit(values)}
    >
      <Form.Item name="name" label="目标名称" rules={[{ required: true, message: '请输入目标名称' }]}>
        <Input placeholder="例如：OnCall Webhook" />
      </Form.Item>
      <Form.Item name="targetType" label="目标类型" rules={[{ required: true, message: '请选择类型' }]}>
        <Select
          options={[
            { label: 'Webhook', value: 'webhook' },
            { label: 'Email', value: 'email' },
            { label: 'SMS', value: 'sms' }
          ]}
        />
      </Form.Item>
      <Form.Item name="configRef" label="配置引用">
        <Input placeholder="例如：secret://ops/oncall" />
      </Form.Item>
      <Form.Item name="status" label="状态">
        <Select
          options={[
            { label: '启用', value: 'active' },
            { label: '禁用', value: 'disabled' }
          ]}
        />
      </Form.Item>
      <Form.Item style={{ marginBottom: 0 }}>
        <Button type="primary" htmlType="submit" loading={loading}>
          新增通知目标
        </Button>
      </Form.Item>
    </Form>
  );
};
