import { Drawer, Form, Input, Select } from 'antd';
import type { CreateTemplateReleasePayload } from '@/services/platformMarketplace';

type TemplateReleaseDrawerProps = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateTemplateReleasePayload) => void;
};

export const TemplateReleaseDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: TemplateReleaseDrawerProps) => {
  const [form] = Form.useForm<CreateTemplateReleasePayload>();

  return (
    <Drawer
      title="发布模板"
      open={open}
      width={460}
      onClose={onClose}
      destroyOnClose
      extra={
        <a
          onClick={() => {
            void form.validateFields().then((values) => onSubmit(values));
          }}
        >
          {submitting ? '提交中...' : '发布'}
        </a>
      }
    >
      <Form layout="vertical" form={form} initialValues={{ targetType: 'workspace' }}>
        <Form.Item name="version" label="目标版本" rules={[{ required: true, message: '请输入版本' }]}>
          <Input placeholder="1.2.0" />
        </Form.Item>
        <Form.Item name="targetType" label="目标范围类型" rules={[{ required: true, message: '请选择范围类型' }]}>
          <Select
            options={[
              { label: '工作空间', value: 'workspace' },
              { label: '项目', value: 'project' },
              { label: '集群', value: 'cluster' }
            ]}
          />
        </Form.Item>
        <Form.Item
          name="targetRef"
          label="目标范围 ID"
          rules={[
            { required: true, message: '请输入目标范围 ID' },
            { pattern: /^\d+$/, message: '目标范围 ID 必须为数字' }
          ]}
        >
          <Input placeholder="例如：12" />
        </Form.Item>
        <Form.Item name="visibilitySummary" label="可见性说明">
          <Input placeholder="scope / target-only" />
        </Form.Item>
        <Form.Item name="releaseNotes" label="发布说明">
          <Input.TextArea rows={4} placeholder="描述变更范围、约束和升级提示" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
