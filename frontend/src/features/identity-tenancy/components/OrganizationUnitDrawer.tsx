import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { CreateOrganizationUnitPayload, OrganizationUnit } from '@/services/identityTenancy';

type OrganizationUnitDrawerProps = {
  open: boolean;
  submitting: boolean;
  units: OrganizationUnit[];
  onClose: () => void;
  onSubmit: (payload: CreateOrganizationUnitPayload) => void;
};

export const OrganizationUnitDrawer = ({
  open,
  submitting,
  units,
  onClose,
  onSubmit
}: OrganizationUnitDrawerProps) => {
  const [form] = Form.useForm<CreateOrganizationUnitPayload>();

  return (
    <Drawer
      title="新建组织单元"
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
              void form.validateFields().then((values) =>
                onSubmit({
                  ...values,
                  description: values.description?.trim() || undefined,
                  parentUnitId: values.parentUnitId || undefined
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
        <Form.Item name="unitType" label="单元类型" rules={[{ required: true, message: '请选择单元类型' }]}>
          <Select
            options={[
              { label: '组织', value: 'organization' },
              { label: '团队', value: 'team' },
              { label: '用户组', value: 'group' }
            ]}
          />
        </Form.Item>
        <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
          <Input placeholder="例如：零售事业部 / 平台工程团队" />
        </Form.Item>
        <Form.Item name="description" label="说明">
          <Input.TextArea rows={3} placeholder="描述责任边界、成员来源和管理责任" />
        </Form.Item>
        <Form.Item name="parentUnitId" label="上级单元">
          <Select
            allowClear
            options={units.map((unit) => ({ label: unit.name, value: unit.id }))}
            placeholder="可选"
          />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
