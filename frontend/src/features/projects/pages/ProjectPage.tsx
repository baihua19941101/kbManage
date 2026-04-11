import { useEffect, useMemo } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Card, Form, Input, Select, Space, Table, Typography, message } from 'antd';
import { normalizeErrorMessage } from '@/app/queryClient';
import {
  createProject,
  listProjectsByWorkspace,
  listWorkspaces
} from '@/services/projects';

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

export const ProjectPage = () => {
  const [form] = Form.useForm<ProjectFormValues>();
  const queryClient = useQueryClient();
  const selectedWorkspaceId = Form.useWatch('workspace', form);

  const { data: workspaces = [], isFetching: isWorkspaceFetching } = useQuery({
    queryKey: ['workspaces'],
    queryFn: listWorkspaces
  });

  useEffect(() => {
    if (selectedWorkspaceId || workspaces.length === 0) {
      return;
    }
    form.setFieldValue('workspace', workspaces[0].id);
  }, [form, selectedWorkspaceId, workspaces]);

  const workspaceOptions = useMemo(
    () => workspaces.map((item) => ({ label: item.name, value: item.id })),
    [workspaces]
  );

  const workspaceNameMap = useMemo(
    () => new Map(workspaces.map((item) => [item.id, item.name])),
    [workspaces]
  );

  const { data: projectItems = [], isFetching: isProjectFetching } = useQuery({
    queryKey: ['workspace-projects', selectedWorkspaceId],
    queryFn: () => listProjectsByWorkspace(String(selectedWorkspaceId)),
    enabled: typeof selectedWorkspaceId === 'string' && selectedWorkspaceId.trim().length > 0
  });

  const createProjectMutation = useMutation({
    mutationFn: (values: ProjectFormValues) =>
      createProject(values.workspace, {
        name: values.name,
        owner: values.owner
      }),
    onSuccess: async (_result, variables) => {
      message.success('项目创建成功');
      form.setFieldValue('name', '');
      form.setFieldValue('owner', '');
      await queryClient.invalidateQueries({
        queryKey: ['workspace-projects', variables.workspace]
      });
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '项目创建失败，请稍后重试'));
    },
    meta: {
      suppressGlobalError: true
    }
  });

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
    createProjectMutation.mutate({
      name: values.name.trim(),
      workspace: values.workspace,
      owner: values.owner.trim()
    });
  };

  const projects = useMemo<Project[]>(
    () =>
      projectItems.map((item, index) => ({
        key: item.id || `${item.workspaceId}-${item.name}-${index}`,
        name: item.name,
        workspace:
          item.workspaceName ||
          workspaceNameMap.get(item.workspaceId) ||
          item.workspaceId ||
          '-',
        owner: item.owner
      })),
    [projectItems, workspaceNameMap]
  );

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
            <Select
              options={workspaceOptions}
              placeholder="选择工作空间"
              loading={isWorkspaceFetching}
            />
          </Form.Item>

          <Form.Item
            label="Owner"
            name="owner"
            rules={[{ required: true, message: '请输入 owner' }]}
          >
            <Input placeholder="例如：admin" maxLength={64} />
          </Form.Item>

          <Button
            type="primary"
            htmlType="submit"
            loading={createProjectMutation.isPending}
            disabled={workspaceOptions.length === 0}
          >
            创建
          </Button>
        </Form>
      </Card>

      <Card title="项目列表" size="small">
        <Table<Project>
          rowKey="key"
          columns={columns}
          dataSource={projects}
          loading={isProjectFetching}
          pagination={{ pageSize: 5 }}
        />
      </Card>
    </Space>
  );
};
