import { useEffect } from 'react';
import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { ArchiveExportScope, CreateArchiveExportTaskRequest } from '@/services/compliance';
import { toDatetimeLocal, toIsoDateTime } from '@/features/compliance-hardening/utils';

type FormValues = {
  exportScope: ArchiveExportScope;
  baselineId?: string;
  timeFrom?: string;
  timeTo?: string;
  filtersText?: string;
};

type ArchiveExportDrawerProps = {
  open: boolean;
  readonly?: boolean;
  loading?: boolean;
  initialValues?: Partial<CreateArchiveExportTaskRequest>;
  onClose: () => void;
  onSubmit: (payload: CreateArchiveExportTaskRequest) => void;
};

const exportScopeOptions: Array<{ label: string; value: ArchiveExportScope }> = [
  { label: '扫描记录', value: 'scans' },
  { label: '失败项', value: 'findings' },
  { label: '趋势', value: 'trends' },
  { label: '审计', value: 'audit' },
  { label: '归档包', value: 'bundle' }
];

const parseFilters = (value?: string) => {
  if (!value?.trim()) {
    return undefined;
  }
  try {
    return JSON.parse(value) as Record<string, unknown>;
  } catch {
    return undefined;
  }
};

export const ArchiveExportDrawer = ({
  open,
  readonly,
  loading,
  initialValues,
  onClose,
  onSubmit
}: ArchiveExportDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }
    form.setFieldsValue({
      exportScope: initialValues?.exportScope || 'bundle',
      baselineId: initialValues?.baselineId || '',
      timeFrom: toDatetimeLocal(initialValues?.timeFrom),
      timeTo: toDatetimeLocal(initialValues?.timeTo),
      filtersText: initialValues?.filters ? JSON.stringify(initialValues.filters, null, 2) : ''
    });
  }, [form, initialValues, open]);

  return (
    <Drawer title="创建归档导出" width={460} open={open} onClose={onClose} destroyOnClose>
      <Form<FormValues>
        form={form}
        layout="vertical"
        onFinish={(values) => {
          onSubmit({
            exportScope: values.exportScope,
            baselineId: values.baselineId || undefined,
            timeFrom: toIsoDateTime(values.timeFrom),
            timeTo: toIsoDateTime(values.timeTo),
            filters: parseFilters(values.filtersText)
          });
        }}
      >
        <Form.Item label="导出范围" name="exportScope" rules={[{ required: true }]}> 
          <Select disabled={readonly} options={exportScopeOptions} />
        </Form.Item>
        <Form.Item label="基线 ID" name="baselineId">
          <Input disabled={readonly} placeholder="可选，按基线聚焦导出" />
        </Form.Item>
        <Form.Item label="开始时间" name="timeFrom">
          <Input disabled={readonly} type="datetime-local" />
        </Form.Item>
        <Form.Item label="结束时间" name="timeTo">
          <Input disabled={readonly} type="datetime-local" />
        </Form.Item>
        <Form.Item label="附加筛选(JSON)" name="filtersText">
          <Input.TextArea disabled={readonly} rows={6} placeholder='{"status":"failed"}' />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" disabled={readonly} loading={loading}>
            创建导出
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
