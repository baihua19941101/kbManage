import { Button, DatePicker, Form, Input, Space } from 'antd';
import dayjs from 'dayjs';
import type { SilencePayload } from '@/services/observability/silences';

type SilenceWindowFormProps = {
  loading?: boolean;
  onSubmit: (payload: SilencePayload) => void | Promise<void>;
};

type FormValues = {
  name: string;
  reason?: string;
  startsAt: dayjs.Dayjs;
  endsAt: dayjs.Dayjs;
};

export const SilenceWindowForm = ({ loading, onSubmit }: SilenceWindowFormProps) => {
  const [form] = Form.useForm<FormValues>();

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={{
        startsAt: dayjs(),
        endsAt: dayjs().add(30, 'minute')
      }}
      onFinish={(values) =>
        void onSubmit({
          name: values.name,
          reason: values.reason,
          startsAt: values.startsAt.toISOString(),
          endsAt: values.endsAt.toISOString()
        })
      }
    >
      <Form.Item name="name" label="静默名称" rules={[{ required: true, message: '请输入静默名称' }]}>
        <Input placeholder="例如：发布窗口" />
      </Form.Item>
      <Form.Item name="reason" label="静默原因">
        <Input.TextArea rows={2} />
      </Form.Item>
      <Space style={{ width: '100%' }} size="middle" align="start">
        <Form.Item
          name="startsAt"
          label="开始时间"
          rules={[{ required: true, message: '请选择开始时间' }]}
        >
          <DatePicker showTime style={{ minWidth: 220 }} />
        </Form.Item>
        <Form.Item
          name="endsAt"
          label="结束时间"
          rules={[{ required: true, message: '请选择结束时间' }]}
        >
          <DatePicker showTime style={{ minWidth: 220 }} />
        </Form.Item>
      </Space>
      <Form.Item style={{ marginBottom: 0 }}>
        <Button type="primary" htmlType="submit" loading={loading}>
          新增静默窗口
        </Button>
      </Form.Item>
    </Form>
  );
};
