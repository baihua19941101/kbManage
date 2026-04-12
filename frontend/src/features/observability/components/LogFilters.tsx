import { useEffect } from 'react';
import { Button, Form, Input, InputNumber, Space } from 'antd';

export type LogFilterValues = {
  clusterId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  workload?: string;
  pod?: string;
  container?: string;
  keyword?: string;
  startAt?: string;
  endAt?: string;
  limit?: number;
};

type LogFiltersProps = {
  initialValues?: LogFilterValues;
  loading?: boolean;
  onSearch: (values: LogFilterValues) => void;
};

export const LogFilters = ({ initialValues, loading, onSearch }: LogFiltersProps) => {
  const [form] = Form.useForm<LogFilterValues>();

  useEffect(() => {
    form.setFieldsValue(initialValues ?? {});
  }, [form, initialValues]);

  const handleFinish = (values: LogFilterValues) => {
    onSearch(values);
  };

  return (
    <Form form={form} layout="vertical" onFinish={handleFinish} initialValues={initialValues}>
      <Space wrap style={{ width: '100%' }} align="start">
        <Form.Item name="clusterId" label="Cluster ID" style={{ width: 180 }}>
          <Input placeholder="cluster-1" />
        </Form.Item>
        <Form.Item name="namespace" label="Namespace" style={{ width: 180 }}>
          <Input placeholder="default" />
        </Form.Item>
        <Form.Item name="resourceKind" label="Kind" style={{ width: 180 }}>
          <Input placeholder="Deployment" />
        </Form.Item>
        <Form.Item name="resourceName" label="Resource" style={{ width: 180 }}>
          <Input placeholder="mock-app" />
        </Form.Item>
      </Space>
      <Space wrap style={{ width: '100%' }} align="start">
        <Form.Item name="workload" label="Workload" style={{ width: 180 }}>
          <Input placeholder="mock-app" />
        </Form.Item>
        <Form.Item name="pod" label="Pod" style={{ width: 180 }}>
          <Input placeholder="mock-app-7b9f" />
        </Form.Item>
        <Form.Item name="container" label="Container" style={{ width: 180 }}>
          <Input placeholder="main" />
        </Form.Item>
        <Form.Item name="keyword" label="关键字" style={{ width: 180 }}>
          <Input placeholder="error" />
        </Form.Item>
      </Space>
      <Space wrap style={{ width: '100%' }} align="start">
        <Form.Item name="startAt" label="开始时间" style={{ width: 220 }}>
          <Input placeholder="2026-04-11T00:00:00Z" />
        </Form.Item>
        <Form.Item name="endAt" label="结束时间" style={{ width: 220 }}>
          <Input placeholder="2026-04-11T01:00:00Z" />
        </Form.Item>
        <Form.Item name="limit" label="条数" style={{ width: 120 }}>
          <InputNumber min={1} max={500} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item label=" ">
          <Button loading={loading} type="primary" htmlType="submit">
            查询
          </Button>
        </Form.Item>
      </Space>
    </Form>
  );
};
