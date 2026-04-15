import { useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Select, Space, Table, Tag, Typography, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { queryKeys } from '@/app/queryClient';
import { normalizeApiError } from '@/services/api/client';
import {
  createPolicyAssignment,
  listPolicyAssignments
} from '@/services/securityPolicy';
import type {
  CreatePolicyAssignmentRequestDTO,
  PolicyAssignmentDTO,
  SecurityPolicyDTO
} from '@/services/api/types';

type PolicyScopeDrawerProps = {
  open: boolean;
  policy?: SecurityPolicyDTO;
  readonly?: boolean;
  onClose: () => void;
  onSuccess?: () => void;
};

type ScopeFormValues = {
  workspaceId?: string;
  projectId?: string;
  clusterRefs?: string;
  namespaceRefs?: string;
  resourceKinds?: string;
  enforcementMode: PolicyAssignmentDTO['enforcementMode'];
  rolloutStage: PolicyAssignmentDTO['rolloutStage'];
};

const splitText = (value?: string): string[] | undefined => {
  const parts =
    value
      ?.split(',')
      .map((item) => item.trim())
      .filter((item) => item.length > 0) ?? [];

  return parts.length > 0 ? parts : undefined;
};

const assignmentColumns: ColumnsType<PolicyAssignmentDTO> = [
  {
    title: '范围',
    key: 'scope',
    render: (_value, record) => {
      const scopes = [
        record.workspaceId ? `workspace:${record.workspaceId}` : '',
        record.projectId ? `project:${record.projectId}` : '',
        record.clusterRefs && record.clusterRefs.length > 0
          ? `clusters:${record.clusterRefs.join(',')}`
          : '',
        record.namespaceRefs && record.namespaceRefs.length > 0
          ? `namespaces:${record.namespaceRefs.join(',')}`
          : ''
      ].filter(Boolean);

      return scopes.length > 0 ? scopes.join(' | ') : '全局';
    }
  },
  {
    title: '资源类型',
    dataIndex: 'resourceKinds',
    key: 'resourceKinds',
    render: (resourceKinds?: string[]) => resourceKinds?.join(', ') || '-'
  },
  {
    title: '模式',
    dataIndex: 'enforcementMode',
    key: 'enforcementMode',
    render: (value: string) => <Tag color={value === 'enforce' ? 'red' : 'blue'}>{value}</Tag>
  },
  {
    title: '阶段',
    dataIndex: 'rolloutStage',
    key: 'rolloutStage'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status'
  }
];

export const PolicyScopeDrawer = ({
  open,
  policy,
  readonly,
  onClose,
  onSuccess
}: PolicyScopeDrawerProps) => {
  const [form] = Form.useForm<ScopeFormValues>();
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      enforcementMode: 'audit',
      rolloutStage: 'pilot'
    });
  }, [form, open]);

  const assignmentsQuery = useQuery({
    queryKey: queryKeys.securityPolicy.assignments(policy?.id),
    enabled: open && Boolean(policy?.id),
    queryFn: () => listPolicyAssignments(policy!.id)
  });

  const createAssignmentMutation = useMutation({
    mutationFn: async (values: ScopeFormValues) => {
      if (!policy?.id) {
        throw new Error('缺少策略 ID，无法分配');
      }

      const payload: CreatePolicyAssignmentRequestDTO = {
        workspaceId: values.workspaceId?.trim() || undefined,
        projectId: values.projectId?.trim() || undefined,
        clusterRefs: splitText(values.clusterRefs),
        namespaceRefs: splitText(values.namespaceRefs),
        resourceKinds: splitText(values.resourceKinds),
        enforcementMode: values.enforcementMode,
        rolloutStage: values.rolloutStage
      };

      return createPolicyAssignment(policy.id, payload);
    },
    onSuccess: () => {
      message.success('策略分配已提交');
      void queryClient.invalidateQueries({ queryKey: queryKeys.securityPolicy.assignments(policy?.id) });
      void queryClient.invalidateQueries({ queryKey: queryKeys.securityPolicy.list() });
      form.resetFields();
      onSuccess?.();
    }
  });

  const handleClose = () => {
    if (createAssignmentMutation.isPending) {
      return;
    }
    form.resetFields();
    onClose();
  };

  return (
    <Drawer
      title={policy ? `分配策略 - ${policy.name}` : '分配策略'}
      open={open}
      width={720}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>关闭</Button>
          <Button
            type="primary"
            disabled={readonly}
            loading={createAssignmentMutation.isPending}
            onClick={() => form.submit()}
          >
            提交分配
          </Button>
        </Space>
      }
    >
      {!policy ? (
        <Typography.Text type="secondary">请选择策略后再进行分配。</Typography.Text>
      ) : (
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          {readonly ? (
            <Alert
              type="info"
              showIcon
              message="当前为只读模式"
              description="你可以查看策略分配结果，但无法提交新的范围分配。"
            />
          ) : null}

          {createAssignmentMutation.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(createAssignmentMutation.error, '策略分配提交失败')}
            />
          ) : null}

          <Form<ScopeFormValues> form={form} layout="vertical" onFinish={(values) => createAssignmentMutation.mutate(values)}>
            <Form.Item label="工作空间 ID" name="workspaceId">
              <Input placeholder="例如：workspace-core" />
            </Form.Item>
            <Form.Item label="项目 ID" name="projectId">
              <Input placeholder="例如：project-payment" />
            </Form.Item>
            <Form.Item label="集群（逗号分隔）" name="clusterRefs">
              <Input placeholder="例如：prod-cn-1,prod-cn-2" />
            </Form.Item>
            <Form.Item label="命名空间（逗号分隔）" name="namespaceRefs">
              <Input placeholder="例如：payments,checkout" />
            </Form.Item>
            <Form.Item label="资源类型（逗号分隔）" name="resourceKinds">
              <Input placeholder="例如：Pod,Deployment" />
            </Form.Item>
            <Form.Item
              label="执行模式"
              name="enforcementMode"
              rules={[{ required: true, message: '请选择执行模式' }]}
            >
              <Select
                options={[
                  { label: 'audit', value: 'audit' },
                  { label: 'alert', value: 'alert' },
                  { label: 'warn', value: 'warn' },
                  { label: 'enforce', value: 'enforce' }
                ]}
              />
            </Form.Item>
            <Form.Item
              label="灰度阶段"
              name="rolloutStage"
              rules={[{ required: true, message: '请选择灰度阶段' }]}
            >
              <Select
                options={[
                  { label: 'pilot', value: 'pilot' },
                  { label: 'canary', value: 'canary' },
                  { label: 'broad', value: 'broad' },
                  { label: 'full', value: 'full' }
                ]}
              />
            </Form.Item>
          </Form>

          <Typography.Title level={5} style={{ margin: 0 }}>
            当前分配
          </Typography.Title>
          {assignmentsQuery.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(assignmentsQuery.error, '分配列表加载失败')}
            />
          ) : null}
          <Table<PolicyAssignmentDTO>
            rowKey={(record) => record.id}
            loading={assignmentsQuery.isLoading || assignmentsQuery.isFetching}
            columns={assignmentColumns}
            dataSource={assignmentsQuery.data?.items || []}
            pagination={{ pageSize: 5 }}
          />
        </Space>
      )}
    </Drawer>
  );
};
