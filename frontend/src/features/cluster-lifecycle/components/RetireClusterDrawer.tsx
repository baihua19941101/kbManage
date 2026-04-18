import { Alert, Button, Checkbox, Drawer, Form, Input, Space } from 'antd';
import type {
  DisableClusterRequest,
  RetireClusterRequest
} from '@/services/clusterLifecycle';

type Props = {
  open: boolean;
  mode: 'disable' | 'retire';
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: DisableClusterRequest | RetireClusterRequest) => void;
};

export const RetireClusterDrawer = ({
  open,
  mode,
  submitting,
  onClose,
  onSubmit
}: Props) => {
  const [form] = Form.useForm<RetireClusterRequest>();
  const isRetire = mode === 'retire';

  return (
    <Drawer
      open={open}
      width={500}
      title={isRetire ? '发起退役流程' : '发起停用流程'}
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Alert
          type={isRetire ? 'warning' : 'info'}
          showIcon
          message={isRetire ? '退役会锁定后续生命周期动作。' : '停用会阻止新的变更动作。'}
        />
        <Form form={form} layout="vertical" onFinish={onSubmit}>
          <Form.Item name="reason" label="原因说明" rules={[{ required: true, message: '请输入原因' }]}>
            <Input.TextArea rows={4} placeholder="说明停用/退役原因、范围与影响。" />
          </Form.Item>
          {isRetire ? (
            <Form.Item name="evidenceNote" label="证据备注">
              <Input.TextArea rows={3} placeholder="记录资源清理、交接或遗留风险说明。" />
            </Form.Item>
          ) : null}
          <Form.Item
            name="confirmation"
            valuePropName="checked"
            rules={[{ validator: (_, value) => (value ? Promise.resolve() : Promise.reject(new Error('请确认已完成复核')))}]}
          >
            <Checkbox>我已确认风险、审计和影响范围</Checkbox>
          </Form.Item>
          <Space>
            <Button onClick={onClose}>取消</Button>
            <Button danger={isRetire} type="primary" htmlType="submit" loading={submitting}>
              {isRetire ? '确认退役' : '确认停用'}
            </Button>
          </Space>
        </Form>
      </Space>
    </Drawer>
  );
};
