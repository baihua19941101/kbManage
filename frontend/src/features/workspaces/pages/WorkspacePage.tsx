import { useMemo, useState } from 'react';
import { Button, Card, Form, Input, Space, Table, Typography, message } from 'antd';

type Workspace = {
  key: string;
  name: string;
  description: string;
};

type WorkspaceFormValues = {
  name: string;
  description?: string;
};

const initialWorkspaces: Workspace[] = [
  { key: 'ws-default', name: 'default', description: '默认工作空间' },
  { key: 'ws-dev', name: 'dev-team', description: '研发团队工作空间' }
];

export const WorkspacePage = () => {
  const [form] = Form.useForm<WorkspaceFormValues>();
  const [workspaces, setWorkspaces] = useState<Workspace[]>(initialWorkspaces);

  const columns = useMemo(
    () => [
      {
        title: '名称',
        dataIndex: 'name',
        key: 'name'
      },
      {
        title: '描述',
        dataIndex: 'description',
        key: 'description'
      }
    ],
    []
  );

  const onCreateWorkspace = (values: WorkspaceFormValues) => {
    const nextName = values.name.trim();
    const existed = workspaces.some((item) => item.name === nextName);
    if (existed) {
      message.warning('工作空间名称已存在');
      return;
    }

    const newItem: Workspace = {
      key: `ws-${Date.now()}`,
      name: nextName,
      description: values.description?.trim() || '-'
    };

    setWorkspaces((prev) => [newItem, ...prev]);
    form.resetFields();
    message.success('工作空间创建成功（mock）');
  };

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={4} style={{ margin: 0 }}>
        工作空间管理
      </Typography.Title>

      <Card title="创建工作空间" size="small">
        <Form<WorkspaceFormValues>
          form={form}
          layout="vertical"
          onFinish={onCreateWorkspace}
        >
          <Form.Item
            label="工作空间名称"
            name="name"
            rules={[
              { required: true, message: '请输入工作空间名称' },
              { min: 2, message: '名称至少 2 个字符' }
            ]}
          >
            <Input placeholder="例如：prod-team" maxLength={64} />
          </Form.Item>
          <Form.Item label="描述" name="description">
            <Input.TextArea placeholder="可选描述" rows={3} maxLength={200} />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            创建
          </Button>
        </Form>
      </Card>

      <Card title="工作空间列表" size="small">
        <Table<Workspace>
          rowKey="key"
          columns={columns}
          dataSource={workspaces}
          pagination={{ pageSize: 5 }}
        />
      </Card>
    </Space>
  );
};
