import { Drawer, Form, Input, Select } from 'antd';
import type { CreateCatalogSourcePayload } from '@/services/platformMarketplace';

type CatalogSourceDrawerProps = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateCatalogSourcePayload) => void;
};

export const CatalogSourceDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: CatalogSourceDrawerProps) => {
  const [form] = Form.useForm<CreateCatalogSourcePayload>();

  return (
    <Drawer
      title="新增目录来源"
      open={open}
      width={420}
      onClose={onClose}
      destroyOnClose
      extra={
        <a
          onClick={() => {
            void form.validateFields().then((values) => onSubmit(values));
          }}
        >
          {submitting ? '提交中...' : '保存'}
        </a>
      }
    >
      <Form
        layout="vertical"
        form={form}
        initialValues={{ sourceType: 'helm', visibleScope: 'platform' }}
      >
        <Form.Item name="name" label="来源名称" rules={[{ required: true, message: '请输入来源名称' }]}>
          <Input placeholder="例如：平台标准 Helm 目录" />
        </Form.Item>
        <Form.Item name="sourceType" label="来源类型" rules={[{ required: true, message: '请选择来源类型' }]}>
          <Select
            options={[
              { label: 'Helm 仓库', value: 'helm' },
              { label: 'OCI 目录', value: 'oci' },
              { label: 'Git 目录', value: 'git' },
              { label: '内置目录', value: 'builtin' }
            ]}
          />
        </Form.Item>
        <Form.Item name="endpoint" label="来源地址" rules={[{ required: true, message: '请输入来源地址' }]}>
          <Input placeholder="https://charts.example.com" />
        </Form.Item>
        <Form.Item name="visibleScope" label="可见范围">
          <Input placeholder="platform 或 workspace/12" />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
