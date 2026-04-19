import { Drawer, Form, Input, Select, Space, Button } from 'antd';
import type { CreateRestoreJobPayload, RestorePoint } from '@/services/backupRestore';

type RestoreJobDrawerProps = {
  open: boolean;
  submitting: boolean;
  restorePoints: RestorePoint[];
  onClose: () => void;
  onSubmit: (payload: CreateRestoreJobPayload) => void;
};

export const RestoreJobDrawer = ({
  open,
  submitting,
  restorePoints,
  onClose,
  onSubmit
}: RestoreJobDrawerProps) => {
  const [form] = Form.useForm();

  return (
    <Drawer
      title="发起恢复任务"
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
                  restorePointId: values.restorePointId,
                  jobType: values.jobType,
                  sourceEnvironment: values.sourceEnvironment?.trim() || undefined,
                  targetEnvironment: values.targetEnvironment.trim(),
                  scopeSelection: {
                    summary: values.scopeSummary?.trim() || '全量恢复'
                  }
                });
              });
            }}
          >
            提交
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="restorePointId"
          label="恢复点"
          rules={[{ required: true, message: '请选择恢复点' }]}
        >
          <Select
            options={restorePoints.map((item) => ({
              label: `${item.id} / ${item.result || '未知结果'}`,
              value: item.id
            }))}
          />
        </Form.Item>
        <Form.Item name="jobType" label="任务类型" rules={[{ required: true, message: '请选择任务类型' }]}>
          <Select
            options={[
              { label: '原地恢复', value: 'in-place' },
              { label: '跨集群恢复', value: 'cross-cluster' },
              { label: '定向恢复', value: 'selective' }
            ]}
          />
        </Form.Item>
        <Form.Item name="sourceEnvironment" label="源环境">
          <Input placeholder="例如：prod-cn" />
        </Form.Item>
        <Form.Item
          name="targetEnvironment"
          label="目标环境"
          rules={[{ required: true, message: '请输入目标环境' }]}
        >
          <Input placeholder="例如：dr-cn" />
        </Form.Item>
        <Form.Item name="scopeSummary" label="恢复范围说明">
          <Input.TextArea rows={4} placeholder="例如：仅恢复 orders 命名空间和权限配置" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
