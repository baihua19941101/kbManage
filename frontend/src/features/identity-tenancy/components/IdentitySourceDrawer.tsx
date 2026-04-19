import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateIdentitySourcePayload } from '@/services/identityTenancy';

type IdentitySourceDrawerProps = {
  open: boolean;
  submitting: boolean;
  onClose: () => void;
  onSubmit: (payload: CreateIdentitySourcePayload) => void;
};

export const IdentitySourceDrawer = ({
  open,
  submitting,
  onClose,
  onSubmit
}: IdentitySourceDrawerProps) => {
  const [form] = Form.useForm<CreateIdentitySourcePayload>();

  return (
    <Drawer
      title="接入身份源"
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
              void form.validateFields().then((values) => onSubmit(values));
            }}
          >
            保存
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <Form.Item name="name" label="身份源名称" rules={[{ required: true, message: '请输入身份源名称' }]}>
          <Input placeholder="例如：企业 OIDC / 总部 LDAP" />
        </Form.Item>
        <Form.Item name="sourceType" label="来源类型" rules={[{ required: true, message: '请选择来源类型' }]}>
          <Select
            options={[
              { label: 'OIDC', value: 'oidc' },
              { label: 'LDAP', value: 'ldap' },
              { label: 'SSO', value: 'sso' },
              { label: '本地账号', value: 'local' }
            ]}
          />
        </Form.Item>
        <Form.Item name="loginMode" label="登录方式" rules={[{ required: true, message: '请选择登录方式' }]}>
          <Select
            options={[
              { label: '本地优先', value: 'local' },
              { label: '外部优先', value: 'external' },
              { label: '并存切换', value: 'mixed' }
            ]}
          />
        </Form.Item>
        <Form.Item name="scopeMode" label="目录范围" rules={[{ required: true, message: '请选择目录范围' }]}>
          <Select
            options={[
              { label: '仅账号同步', value: 'account-only' },
              { label: '账号与组织同步', value: 'organization-sync' },
              { label: '只做认证，不同步目录', value: 'auth-only' }
            ]}
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
