import { Form, Modal, Select, Typography, message } from 'antd';
import { exportAuditEvents, type AuditEventFilters, type AuditExportFormat } from '@/services/audit';

export type AuditExportModalProps = {
  open: boolean;
  filters: AuditEventFilters;
  onCancel: () => void;
  onSubmitted?: (taskId: string) => void;
};

type ExportFormValues = {
  format: AuditExportFormat;
};

export const AuditExportModal = ({
  open,
  filters,
  onCancel,
  onSubmitted
}: AuditExportModalProps) => {
  const [form] = Form.useForm<ExportFormValues>();

  const onSubmit = async () => {
    const values = await form.validateFields();

    if (!filters.from || !filters.to) {
      message.warning('请先选择导出时间范围');
      return;
    }

    const result = await exportAuditEvents({
      from: filters.from,
      to: filters.to,
      actorUserId: filters.actorUserId,
      clusterId: filters.clusterId,
      result: filters.result,
      eventType: filters.eventType,
      format: values.format
    });

    message.success(
      result.status === 'mocked'
        ? `导出任务已创建（mock）：${result.taskId}`
        : `导出任务已提交：${result.taskId}`
    );

    onSubmitted?.(result.taskId);
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title="导出审计记录"
      open={open}
      onCancel={onCancel}
      onOk={() => void onSubmit()}
      destroyOnHidden
      okText="提交导出"
      cancelText="取消"
    >
      <Typography.Paragraph type="secondary">
        将按当前筛选条件提交导出任务，建议先设置明确时间范围。
      </Typography.Paragraph>
      <Form<ExportFormValues>
        form={form}
        layout="vertical"
        initialValues={{ format: 'csv' }}
      >
        <Form.Item
          label="导出格式"
          name="format"
          rules={[{ required: true, message: '请选择导出格式' }]}
        >
          <Select
            options={[
              { value: 'csv', label: 'CSV（推荐）' },
              { value: 'json', label: 'JSON' }
            ]}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};
