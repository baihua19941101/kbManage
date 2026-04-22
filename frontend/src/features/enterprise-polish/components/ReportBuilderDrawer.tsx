import { Button, Drawer, Form, Input, Select } from 'antd';
import { useState } from 'react';
import type { CreateGovernanceReportPayload } from '@/services/enterprisePolish';

type Props = { onSubmit: (payload: CreateGovernanceReportPayload) => Promise<void> | void };

export const ReportBuilderDrawer = ({ onSubmit }: Props) => {
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm<CreateGovernanceReportPayload>();
  return (
    <>
      <Button type="primary" onClick={() => setOpen(true)}>
        生成报表
      </Button>
      <Drawer open={open} title="生成治理报表" onClose={() => setOpen(false)} destroyOnClose>
        <Form
          form={form}
          layout="vertical"
          initialValues={{ workspaceId: 1, reportType: 'management', audienceType: 'leadership' }}
          onFinish={async (values) => {
            await onSubmit(values);
            setOpen(false);
            form.resetFields();
          }}
        >
          <Form.Item name="title" label="标题" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="reportType" label="报表类型" rules={[{ required: true }]}>
            <Select options={[{ value: 'management' }, { value: 'audit' }, { value: 'delivery' }]} />
          </Form.Item>
          <Form.Item name="audienceType" label="面向对象" rules={[{ required: true }]}>
            <Select options={[{ value: 'leadership' }, { value: 'auditor' }, { value: 'customer' }]} />
          </Form.Item>
          <Form.Item name="timeRange" label="时间范围">
            <Input />
          </Form.Item>
          <Button htmlType="submit" type="primary">
            提交
          </Button>
        </Form>
      </Drawer>
    </>
  );
};
