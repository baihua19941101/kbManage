import { useEffect } from 'react';
import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type {
  CreateRemediationTaskRequest,
  RemediationTask,
  RemediationTaskStatus,
  RiskLevel,
  UpdateRemediationTaskRequest
} from '@/services/compliance';
import { toDatetimeLocal, toIsoDateTime } from '@/features/compliance-hardening/utils';

type FormValues = {
  title: string;
  owner: string;
  priority: RiskLevel;
  dueAt?: string;
  status?: RemediationTaskStatus;
  resolutionSummary?: string;
};

type RemediationTaskDrawerProps = {
  open: boolean;
  task?: RemediationTask;
  findingId?: string;
  readonly?: boolean;
  loading?: boolean;
  onClose: () => void;
  onCreate: (findingId: string, payload: CreateRemediationTaskRequest) => void;
  onUpdate: (taskId: string, payload: UpdateRemediationTaskRequest) => void;
};

const priorityOptions: Array<{ label: string; value: RiskLevel }> = [
  { label: '低', value: 'low' },
  { label: '中', value: 'medium' },
  { label: '高', value: 'high' },
  { label: '严重', value: 'critical' }
];

const statusOptions: Array<{ label: string; value: RemediationTaskStatus }> = [
  { label: '待处理', value: 'todo' },
  { label: '处理中', value: 'in_progress' },
  { label: '阻塞', value: 'blocked' },
  { label: '完成', value: 'done' },
  { label: '取消', value: 'canceled' }
];

export const RemediationTaskDrawer = ({
  open,
  task,
  findingId,
  readonly,
  loading,
  onClose,
  onCreate,
  onUpdate
}: RemediationTaskDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }
    form.setFieldsValue({
      title: task?.title || '',
      owner: task?.owner || '',
      priority: task?.priority || 'medium',
      dueAt: toDatetimeLocal(task?.dueAt),
      status: task?.status || 'todo',
      resolutionSummary: task?.resolutionSummary || ''
    });
  }, [form, open, task]);

  return (
    <Drawer
      title={task?.id ? '更新整改任务' : '创建整改任务'}
      width={460}
      open={open}
      onClose={onClose}
      destroyOnClose
    >
      <Form<FormValues>
        form={form}
        layout="vertical"
        onFinish={(values) => {
          if (task?.id) {
            onUpdate(task.id, {
              status: values.status,
              resolutionSummary: values.resolutionSummary || undefined
            });
            return;
          }
          if (!findingId) {
            return;
          }
          onCreate(findingId, {
            title: values.title,
            owner: values.owner,
            priority: values.priority,
            dueAt: toIsoDateTime(values.dueAt),
            summary: values.resolutionSummary || undefined
          });
        }}
      >
        <Form.Item label="任务标题" name="title" rules={[{ required: true, message: '请输入任务标题' }]}>
          <Input disabled={readonly || Boolean(task?.id)} />
        </Form.Item>
        <Form.Item label="负责人" name="owner" rules={[{ required: true, message: '请输入负责人' }]}>
          <Input disabled={readonly || Boolean(task?.id)} />
        </Form.Item>
        <Form.Item label="优先级" name="priority" rules={[{ required: true }]}>
          <Select disabled={readonly || Boolean(task?.id)} options={priorityOptions} />
        </Form.Item>
        <Form.Item label="到期时间" name="dueAt">
          <Input type="datetime-local" disabled={readonly || Boolean(task?.id)} />
        </Form.Item>
        {task?.id ? (
          <Form.Item label="状态" name="status">
            <Select disabled={readonly} options={statusOptions} />
          </Form.Item>
        ) : null}
        <Form.Item label={task?.id ? '处理说明' : '补充说明'} name="resolutionSummary">
          <Input.TextArea disabled={readonly} rows={4} />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" disabled={readonly} loading={loading}>
            保存
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
