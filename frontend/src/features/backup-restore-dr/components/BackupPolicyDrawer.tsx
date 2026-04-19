import { Drawer, Form, Input, Select, Space, Button } from 'antd';
import type { CreateBackupPolicyPayload } from '@/services/backupRestore';

type BackupPolicyDrawerProps = {
  open: boolean;
  submitting: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateBackupPolicyPayload) => void;
};

const scopeOptions = [
  { label: '平台元数据', value: 'platform-metadata' },
  { label: '权限配置', value: 'rbac' },
  { label: '审计记录', value: 'audit' },
  { label: '集群配置', value: 'cluster-config' },
  { label: '关键命名空间', value: 'namespace' }
];

export const BackupPolicyDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: BackupPolicyDrawerProps) => {
  const [form] = Form.useForm<CreateBackupPolicyPayload>();

  return (
    <Drawer
      title="新建备份策略"
      width={520}
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
                  ...values,
                  description: values.description?.trim() || undefined,
                  scopeRef: values.scopeRef?.trim() || undefined,
                  scheduleExpression: values.scheduleExpression?.trim() || undefined
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
        <Form.Item name="name" label="策略名称" rules={[{ required: true, message: '请输入策略名称' }]}>
          <Input placeholder="例如：平台核心元数据每日备份" />
        </Form.Item>
        <Form.Item name="description" label="策略说明">
          <Input.TextArea rows={3} placeholder="描述覆盖范围和负责人" />
        </Form.Item>
        <Form.Item name="scopeType" label="保护范围类型" rules={[{ required: true, message: '请选择范围类型' }]}>
          <Select options={scopeOptions} />
        </Form.Item>
        <Form.Item name="scopeRef" label="范围引用">
          <Input placeholder="例如：workspace-a / cluster-prod / namespace:orders" />
        </Form.Item>
        <Form.Item
          name="executionMode"
          label="执行方式"
          rules={[{ required: true, message: '请选择执行方式' }]}
        >
          <Select
            options={[
              { label: '定时执行', value: 'scheduled' },
              { label: '手动执行', value: 'manual' }
            ]}
          />
        </Form.Item>
        <Form.Item name="scheduleExpression" label="执行计划">
          <Input placeholder="例如：0 2 * * *" />
        </Form.Item>
        <Form.Item
          name="retentionRule"
          label="保留规则"
          rules={[{ required: true, message: '请输入保留规则' }]}
        >
          <Input placeholder="例如：保留 14 天，每周保留 4 个" />
        </Form.Item>
        <Form.Item
          name="consistencyLevel"
          label="一致性级别"
          rules={[{ required: true, message: '请选择一致性级别' }]}
        >
          <Select
            options={[
              { label: '强一致', value: 'strict' },
              { label: '平衡', value: 'balanced' },
              { label: '尽力而为', value: 'best-effort' }
            ]}
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
