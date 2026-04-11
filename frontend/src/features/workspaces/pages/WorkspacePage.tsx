import { useMemo } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Form, Input, Space, Table, Typography, message } from 'antd';
import { normalizeErrorMessage } from '@/app/queryClient';
import {
  createWorkspace,
  listWorkspaces,
  type Workspace
} from '@/services/workspaces';

const WORKSPACES_QUERY_KEY = ['workspaces'] as const;

type WorkspaceFormValues = {
  name: string;
  description?: string;
};

export const WorkspacePage = () => {
  const [form] = Form.useForm<WorkspaceFormValues>();
  const queryClient = useQueryClient();

  const { data, isFetching, error } = useQuery({
    queryKey: WORKSPACES_QUERY_KEY,
    queryFn: listWorkspaces,
    meta: {
      suppressGlobalError: true
    }
  });

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

  const createWorkspaceMutation = useMutation({
    mutationFn: createWorkspace,
    onSuccess: async (workspace) => {
      form.resetFields();
      message.success(`工作空间 ${workspace.name} 创建成功`);
      await queryClient.invalidateQueries({
        queryKey: WORKSPACES_QUERY_KEY
      });
    },
    onError: (err) => {
      message.error(normalizeErrorMessage(err, '工作空间创建失败，请稍后重试'));
    }
  });

  const onCreateWorkspace = (values: WorkspaceFormValues) => {
    createWorkspaceMutation.mutate({
      name: values.name,
      description: values.description
    });
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
          <Button type="primary" htmlType="submit" loading={createWorkspaceMutation.isPending}>
            创建
          </Button>
        </Form>
      </Card>

      <Card title="工作空间列表" size="small">
        {error ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
            message="工作空间列表加载失败"
            description={normalizeErrorMessage(error)}
          />
        ) : null}
        <Table<Workspace>
          rowKey="id"
          columns={columns}
          dataSource={data || []}
          loading={isFetching}
          pagination={{ pageSize: 5 }}
        />
      </Card>
    </Space>
  );
};
