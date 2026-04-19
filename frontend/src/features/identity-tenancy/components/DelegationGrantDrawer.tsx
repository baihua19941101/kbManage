import { Button, DatePicker, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateDelegationGrantPayload } from '@/services/identityTenancy';

type DelegationGrantDrawerProps = {
  open: boolean;
  submitting: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateDelegationGrantPayload) => void;
};

export const DelegationGrantDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: DelegationGrantDrawerProps) => {
  const [form] = Form.useForm();

  return (
    <Drawer
      title="新建委派关系"
      width={560}
      open={open}
      onClose={onClose}
      destroyOnClose
      extra={
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button
            type="primary"
            loading={submitting}
            onClick={() => {
              void form.validateFields().then((values) =>
                onSubmit({
                  grantorRef: values.grantorRef?.trim(),
                  delegateRef: values.delegateRef?.trim(),
                  allowedRoleLevels: values.allowedRoleLevels,
                  validFrom: values.validFrom.toISOString(),
                  validUntil: values.validUntil.toISOString(),
                  reason: values.reason?.trim() || undefined
                })
              );
            }}
          >
            保存
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item name="grantorRef" label="委派人" rules={[{ required: true, message: '请输入委派人' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="delegateRef" label="被委派人" rules={[{ required: true, message: '请输入被委派人' }]}>
          <Input />
        </Form.Item>
        <Form.Item
          name="allowedRoleLevels"
          label="允许委派层级"
          rules={[{ required: true, message: '请选择允许委派层级' }]}
        >
          <Select
            mode="multiple"
            options={[
              { label: '平台级', value: 'platform' },
              { label: '组织级', value: 'organization' },
              { label: '工作空间级', value: 'workspace' },
              { label: '项目级', value: 'project' },
              { label: '资源级', value: 'resource' }
            ]}
          />
        </Form.Item>
        <Form.Item name="validFrom" label="生效时间" rules={[{ required: true, message: '请选择生效时间' }]}>
          <DatePicker showTime style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="validUntil" label="到期时间" rules={[{ required: true, message: '请选择到期时间' }]}>
          <DatePicker showTime style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="reason" label="原因">
          <Input.TextArea rows={3} placeholder="描述委派背景和限制条件" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
