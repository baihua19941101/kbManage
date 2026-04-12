import { Button, Form, Input, Select, Space } from 'antd';
import type { AlertRulePayload } from '@/services/observability/alertRules';

type AlertRuleFormProps = {
  loading?: boolean;
  disabled?: boolean;
  onSubmit: (payload: AlertRulePayload) => void | Promise<void>;
};

export const AlertRuleForm = ({ loading, disabled, onSubmit }: AlertRuleFormProps) => {
  const [form] = Form.useForm<AlertRulePayload>();

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={{ severity: 'warning', status: 'enabled' }}
      onFinish={(values) => void onSubmit(values)}
    >
      <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }]}>
        <Input placeholder="例如：CPU 使用率过高" disabled={disabled} />
      </Form.Item>
      <Form.Item
        name="conditionExpression"
        label="触发表达式"
        rules={[{ required: true, message: '请输入触发表达式' }]}
      >
        <Input placeholder="例如：cpu_usage > 80" disabled={disabled} />
      </Form.Item>
      <Space style={{ width: '100%' }} size="middle" align="start">
        <Form.Item name="severity" label="级别" style={{ minWidth: 180 }}>
          <Select
            disabled={disabled}
            options={[
              { label: 'Info', value: 'info' },
              { label: 'Warning', value: 'warning' },
              { label: 'Critical', value: 'critical' }
            ]}
          />
        </Form.Item>
        <Form.Item name="status" label="状态" style={{ minWidth: 180 }}>
          <Select
            disabled={disabled}
            options={[
              { label: '启用', value: 'enabled' },
              { label: '禁用', value: 'disabled' }
            ]}
          />
        </Form.Item>
      </Space>
      <Form.Item name="description" label="描述">
        <Input.TextArea rows={2} disabled={disabled} />
      </Form.Item>
      <Form.Item style={{ marginBottom: 0 }}>
        <Button type="primary" htmlType="submit" loading={loading} disabled={disabled}>
          新增规则
        </Button>
      </Form.Item>
    </Form>
  );
};
