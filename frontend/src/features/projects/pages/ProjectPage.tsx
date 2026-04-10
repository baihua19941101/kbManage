import { useMemo, useState } from 'react';
import { Button, Card, Form, Input, Select, Space, Table, Typography, message } from 'antd';

type Project = {
  key: string;
  name: string;
  workspace: string;
  owner: string;
};

type ProjectFormValues = {
  name: string;
  workspace: string;
  owner: string;
};

const mockWorkspaceOptions = [
  { label: 'default', value: 'default' },
  { label: 'dev-team', value: 'dev-team' },
  { label: 'ops-team', value: 'ops-team' }
];

const initialProjects: Project[] = [
  { key: 'prj-kb', name: 'kb-manage', workspace: 'default', owner: 'admin' },
  { key: 'prj-console', name: 'cluster-console', workspace: 'dev-team', owner: 'alice' }
];

export const ProjectPage = () => {
  const [form] = Form.useForm<ProjectFormValues>();
  const [projects, setProjects] = useState<Project[]>(initialProjects);

  const columns = useMemo(
    () => [
      {
        title: '项目名称',
        dataIndex: 'name',
        key: 'name'
      },
      {
        title: '所属工作空间',
        dataIndex: 'workspace',
        key: 'workspace'
      },
      {
        title: 'Owner',
        dataIndex: 'owner',
        key: 'owner'
      }
    ],
    []
  );

  const onCreateProject = (values: ProjectFormValues) => {
    const normalized = values.name.trim();
    const existed = projects.some(
      (item) => item.name === normalized && item.workspace === values.workspace
    );

    if (existed) {
      message.warning('同一工作空间下项目名称已存在');
      return;
    }

    setProjects((prev) => [
      {
        key: `prj-${Date.now()}`,
        name: normalized,
        workspace: values.workspace,
        owner: values.owner.trim()
      },
      ...prev
    ]);

    form.resetFields();
    message.success('项目创建成功（mock）');
  };

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={4} style={{ margin: 0 }}>
        项目管理
      </Typography.Title>

      <Card title="创建项目" size="small">
        <Form<ProjectFormValues>
          form={form}
          layout="vertical"
          onFinish={onCreateProject}
        >
          <Form.Item
            label="项目名称"
            name="name"
            rules={[
              { required: true, message: '请输入项目名称' },
              { min: 2, message: '项目名称至少 2 个字符' }
            ]}
          >
            <Input placeholder="例如：billing-api" maxLength={64} />
          </Form.Item>

          <Form.Item
            label="所属工作空间"
            name="workspace"
            rules={[{ required: true, message: '请选择工作空间' }]}
          >
            <Select options={mockWorkspaceOptions} placeholder="选择工作空间" />
          </Form.Item>

          <Form.Item
            label="Owner"
            name="owner"
            rules={[{ required: true, message: '请输入 owner' }]}
          >
            <Input placeholder="例如：admin" maxLength={64} />
          </Form.Item>

          <Button type="primary" htmlType="submit">
            创建
          </Button>
        </Form>
      </Card>

      <Card title="项目列表" size="small">
        <Table<Project>
          rowKey="key"
          columns={columns}
          dataSource={projects}
          pagination={{ pageSize: 5 }}
        />
      </Card>
    </Space>
  );
};
