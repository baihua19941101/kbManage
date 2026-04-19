import { Button, Drawer, Form, Input, Space } from 'antd';
import type { CreateMigrationPlanPayload } from '@/services/backupRestore';

type MigrationPlanDrawerProps = {
  open: boolean;
  submitting: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateMigrationPlanPayload) => void;
};

export const MigrationPlanDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: MigrationPlanDrawerProps) => {
  const [form] = Form.useForm();

  return (
    <Drawer
      title="新建迁移计划"
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
              void form.validateFields().then((values) => {
                onSubmit({
                  name: values.name.trim(),
                  sourceClusterId: values.sourceClusterId.trim(),
                  targetClusterId: values.targetClusterId.trim(),
                  scopeSelection: { summary: values.scopeSummary?.trim() || '全量迁移' },
                  mappingRules: { summary: values.mappingRules?.trim() || '按默认映射' },
                  cutoverSteps: values.cutoverSteps
                    ?.split('\n')
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
          <Input placeholder="例如：上海主集群到灾备集群迁移" />
        </Form.Item>
        <Form.Item
          name="sourceClusterId"
          label="源集群"
          rules={[{ required: true, message: '请输入源集群 ID' }]}
        >
          <Input placeholder="cluster-prod-cn" />
        </Form.Item>
        <Form.Item
          name="targetClusterId"
          label="目标集群"
          rules={[{ required: true, message: '请输入目标集群 ID' }]}
        >
          <Input placeholder="cluster-dr-cn" />
        </Form.Item>
        <Form.Item name="scopeSummary" label="迁移范围">
          <Input.TextArea rows={3} placeholder="例如：订单、支付命名空间与平台 RBAC" />
        </Form.Item>
        <Form.Item name="mappingRules" label="映射规则">
          <Input.TextArea rows={3} placeholder="例如：命名空间前缀保留，存储类映射到 fast-ssd" />
        </Form.Item>
        <Form.Item name="cutoverSteps" label="切换步骤">
          <Input.TextArea rows={5} placeholder={'逐行输入步骤，例如：\n冻结写流量\n执行增量恢复\n切换入口'} />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
