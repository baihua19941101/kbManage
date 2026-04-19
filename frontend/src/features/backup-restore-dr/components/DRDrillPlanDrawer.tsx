import { Button, Drawer, Form, Input, InputNumber, Space } from 'antd';
import type { CreateDRDrillPlanPayload } from '@/services/backupRestore';

type DRDrillPlanDrawerProps = {
  open: boolean;
  submitting: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateDRDrillPlanPayload) => void;
};

export const DRDrillPlanDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: DRDrillPlanDrawerProps) => {
  const [form] = Form.useForm();

  return (
    <Drawer
      title="新建灾备演练计划"
      width={600}
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
              void form.validateFields().then((values) => {
                onSubmit({
                  name: values.name.trim(),
                  description: values.description?.trim() || undefined,
                  scopeSelection: { summary: values.scopeSummary?.trim() || '关键平台对象' },
                  rpoTargetMinutes: values.rpoTargetMinutes,
                  rtoTargetMinutes: values.rtoTargetMinutes,
                  roleAssignments: values.roleAssignments
                    ?.split('\n')
                    .map((item: string) => item.trim())
                    .filter(Boolean),
                  cutoverProcedure: values.cutoverProcedure
                    .split('\n')
                    .map((item: string) => item.trim())
                    .filter(Boolean),
                  validationChecklist: values.validationChecklist
                    .split('\n')
                    .map((item: string) => item.trim())
                    .filter(Boolean)
                });
              });
            }}
          >
            保存
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item name="name" label="计划名称" rules={[{ required: true, message: '请输入计划名称' }]}>
          <Input placeholder="季度级灾备演练" />
        </Form.Item>
        <Form.Item name="description" label="计划说明">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item name="scopeSummary" label="演练范围">
          <Input.TextArea rows={3} placeholder="例如：平台元数据、RBAC、orders 命名空间" />
        </Form.Item>
        <Form.Item
          name="rpoTargetMinutes"
          label="RPO 目标（分钟）"
          rules={[{ required: true, message: '请输入 RPO 目标' }]}
        >
          <InputNumber min={1} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item
          name="rtoTargetMinutes"
          label="RTO 目标（分钟）"
          rules={[{ required: true, message: '请输入 RTO 目标' }]}
        >
          <InputNumber min={1} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="roleAssignments" label="角色分工">
          <Input.TextArea rows={4} placeholder={'逐行输入，例如：\nSRE：执行切换\n业务负责人：验证订单'} />
        </Form.Item>
        <Form.Item
          name="cutoverProcedure"
          label="切换步骤"
          rules={[{ required: true, message: '请输入切换步骤' }]}
        >
          <Input.TextArea rows={5} placeholder={'逐行输入步骤'} />
        </Form.Item>
        <Form.Item
          name="validationChecklist"
          label="验证清单"
          rules={[{ required: true, message: '请输入验证清单' }]}
        >
          <Input.TextArea rows={5} placeholder={'逐行输入验证项'} />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
