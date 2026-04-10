import { Alert, Button, Checkbox, Drawer, Form, Input, InputNumber, Select, Space, Typography } from 'antd';
import { useEffect } from 'react';
import type { CreateOperationPayload, OperationRiskLevel, OperationType } from '@/services/operations';

type OperationConfirmDrawerProps = {
  open: boolean;
  actionType: OperationType;
  resourceName: string;
  onClose: () => void;
  onConfirm: (payload: Omit<CreateOperationPayload, 'target'>) => Promise<void>;
  submitting?: boolean;
};

type ConfirmFormValues = {
  riskLevel: OperationRiskLevel;
  reason: string;
  expectedText: string;
  acknowledged: boolean;
  scaleReplicas?: number;
};

const actionMeta: Record<
  OperationType,
  {
    title: string;
    riskNotice: string;
    defaultRisk: OperationRiskLevel;
    requireScale?: boolean;
  }
> = {
  scale: {
    title: '发起扩缩容',
    riskNotice: '扩缩容会触发副本变更，可能引起容量波动或短时抖动。',
    defaultRisk: 'medium',
    requireScale: true
  },
  restart: {
    title: '发起重启',
    riskNotice: '重启会中断当前实例，需确认业务具备可用副本或熔断策略。',
    defaultRisk: 'high'
  },
  'node-maintenance': {
    title: '发起节点维护',
    riskNotice: '节点维护可能导致 Pod 驱逐与重调度，注意维护窗口和容量冗余。',
    defaultRisk: 'high'
  }
};

export const OperationConfirmDrawer = ({
  open,
  actionType,
  resourceName,
  onClose,
  onConfirm,
  submitting
}: OperationConfirmDrawerProps) => {
  const [form] = Form.useForm<ConfirmFormValues>();
  const meta = actionMeta[actionType];

  useEffect(() => {
    if (open) {
      form.setFieldsValue({
        riskLevel: meta.defaultRisk,
        reason: '',
        expectedText: '',
        acknowledged: false,
        scaleReplicas: actionType === 'scale' ? 3 : undefined
      });
    }
  }, [actionType, form, meta.defaultRisk, open]);

  const onSubmit = async () => {
    const values = await form.validateFields();
    await onConfirm({
      type: actionType,
      riskLevel: values.riskLevel,
      reason: values.reason,
      expectedText: values.expectedText,
      scaleReplicas: actionType === 'scale' ? values.scaleReplicas : undefined
    });
    form.resetFields();
    onClose();
  };

  return (
    <Drawer
      title={meta.title}
      width={500}
      open={open}
      onClose={onClose}
      destroyOnClose
      extra={
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" loading={submitting} onClick={() => void onSubmit()}>
            确认执行
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <Alert type="warning" showIcon message="高风险操作确认" description={meta.riskNotice} />
        <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
          为防止误操作，请填写变更原因，并输入目标资源名称 <Typography.Text code>{resourceName}</Typography.Text> 完成二次确认。
        </Typography.Paragraph>
        <Form layout="vertical" form={form}>
          <Form.Item label="风险等级" name="riskLevel" rules={[{ required: true, message: '请选择风险等级' }]}>
            <Select
              options={[
                { label: 'low', value: 'low' },
                { label: 'medium', value: 'medium' },
                { label: 'high', value: 'high' }
              ]}
            />
          </Form.Item>
          {meta.requireScale ? (
            <Form.Item
              label="目标副本数"
              name="scaleReplicas"
              rules={[{ required: true, message: '请输入副本数' }]}
            >
              <InputNumber min={1} max={500} precision={0} style={{ width: '100%' }} />
            </Form.Item>
          ) : null}
          <Form.Item label="变更原因" name="reason" rules={[{ required: true, message: '请输入变更原因' }]}>
            <Input.TextArea rows={3} placeholder="请说明本次操作原因与影响范围" />
          </Form.Item>
          <Form.Item
            label="二次确认"
            name="expectedText"
            rules={[
              { required: true, message: '请输入资源名称完成确认' },
              {
                validator: (_, value) =>
                  value === resourceName
                    ? Promise.resolve()
                    : Promise.reject(new Error('输入内容与资源名称不一致'))
              }
            ]}
          >
            <Input placeholder={`请输入 ${resourceName}`} />
          </Form.Item>
          <Form.Item
            name="acknowledged"
            valuePropName="checked"
            rules={[
              {
                validator: (_, value) =>
                  value ? Promise.resolve() : Promise.reject(new Error('请勾选风险确认'))
              }
            ]}
          >
            <Checkbox>我已确认操作风险并知晓回滚方案</Checkbox>
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};
