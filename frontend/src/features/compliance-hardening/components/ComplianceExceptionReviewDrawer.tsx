import { useEffect } from 'react';
import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type {
  ComplianceExceptionRequest,
  ReviewComplianceExceptionRequest
} from '@/services/compliance';

type FormValues = ReviewComplianceExceptionRequest;

type ComplianceExceptionReviewDrawerProps = {
  open: boolean;
  exception?: ComplianceExceptionRequest;
  readonly?: boolean;
  loading?: boolean;
  onClose: () => void;
  onSubmit: (exceptionId: string, payload: ReviewComplianceExceptionRequest) => void;
};

const decisionOptions: Array<{ label: string; value: ReviewComplianceExceptionRequest['decision'] }> = [
  { label: '批准', value: 'approve' },
  { label: '拒绝', value: 'reject' },
  { label: '撤销', value: 'revoke' }
];

export const ComplianceExceptionReviewDrawer = ({
  open,
  exception,
  readonly,
  loading,
  onClose,
  onSubmit
}: ComplianceExceptionReviewDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }
    form.setFieldsValue({
      decision: 'approve',
      reviewComment: exception?.reviewComment || ''
    });
  }, [exception?.reviewComment, form, open]);

  return (
    <Drawer title="审批合规例外" width={420} open={open} onClose={onClose} destroyOnClose>
      <Form<FormValues>
        form={form}
        layout="vertical"
        onFinish={(values) => {
          if (!exception?.id) {
            return;
          }
          onSubmit(exception.id, values);
        }}
      >
        <Form.Item label="审批结论" name="decision" rules={[{ required: true }]}> 
          <Select disabled={readonly} options={decisionOptions} />
        </Form.Item>
        <Form.Item label="审批说明" name="reviewComment">
          <Input.TextArea disabled={readonly} rows={4} />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" disabled={readonly || !exception?.id} loading={loading}>
            提交
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
