import { useEffect, useState } from 'react';
import { Alert, Form, Modal, Select, Typography, message } from 'antd';
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
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string>();

  useEffect(() => {
    if (!open) {
      setSubmitError(undefined);
      setSubmitting(false);
      form.resetFields();
      return;
    }
    form.setFieldValue('format', 'csv');
  }, [form, open]);

  const mapErrorMessage = (error: unknown): string => {
    if (error instanceof Error && error.message.trim().length > 0) {
      return error.message;
    }
    return '导出任务提交失败，请稍后重试。';
  };

  const onSubmit = async () => {
    try {
      setSubmitError(undefined);
      const values = await form.validateFields();

      if (!filters.from || !filters.to) {
        message.warning('请先选择导出时间范围');
        return;
      }

      setSubmitting(true);
      const result = await exportAuditEvents({
        from: filters.from,
        to: filters.to,
        actorUserId: filters.actorUserId,
        clusterId: filters.clusterId,
        result: filters.result,
        eventType: filters.eventType,
        format: values.format
      });

      message.success(`导出任务已提交：${result.taskId}`);

      onSubmitted?.(result.taskId);
      form.resetFields();
      onCancel();
    } catch (error) {
      setSubmitError(mapErrorMessage(error));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal
      title="导出审计记录"
      open={open}
      onCancel={onCancel}
      onOk={() => void onSubmit()}
      destroyOnHidden
      confirmLoading={submitting}
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
              { value: 'csv', label: 'CSV（推荐）' }
            ]}
          />
        </Form.Item>
      </Form>
      {submitError ? (
        <Alert
          type="error"
          showIcon
          message="导出失败"
          description={`${submitError} 可直接再次点击“提交导出”重试。`}
        />
      ) : null}
    </Modal>
  );
};
