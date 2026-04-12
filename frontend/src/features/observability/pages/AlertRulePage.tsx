import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Table, Typography, message } from 'antd';
import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { canManageObservability, canReadObservability, useAuthStore } from '@/features/auth/store';
import { AlertRuleForm } from '@/features/observability/components/AlertRuleForm';
import { ApiError, normalizeApiError } from '@/services/api/client';
import {
  createAlertRule,
  deleteAlertRule,
  listAlertRules
} from '@/services/observability/alertRules';

const isAuthorizationError = (error: unknown): boolean =>
  error instanceof ApiError && (error.status === 401 || error.status === 403);

export const AlertRulePage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadObservability(user);
  const canManage = canManageObservability(user);
  const [permissionMessage, setPermissionMessage] = useState<string>();
  const [msgApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const rulesQuery = useQuery({
    queryKey: ['observability', 'alert-rules'],
    queryFn: () => listAlertRules(),
    enabled: canRead
  });

  const createMutation = useMutation({
    mutationFn: createAlertRule,
    onSuccess: async () => {
      setPermissionMessage(undefined);
      await queryClient.invalidateQueries({ queryKey: ['observability', 'alert-rules'] });
      msgApi.success('规则已创建');
    },
    onError: (err) => {
      if (isAuthorizationError(err)) {
        setPermissionMessage(
          normalizeApiError(err, '权限已回收，无法继续创建告警规则。请刷新页面后重试。')
        );
        return;
      }
      msgApi.error(`创建规则失败：${err}`);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteAlertRule,
    onSuccess: async () => {
      setPermissionMessage(undefined);
      await queryClient.invalidateQueries({ queryKey: ['observability', 'alert-rules'] });
      msgApi.success('规则已删除');
    },
    onError: (err) => {
      if (isAuthorizationError(err)) {
        setPermissionMessage(
          normalizeApiError(err, '权限已回收，无法继续删除告警规则。请刷新页面后重试。')
        );
        return;
      }
      msgApi.error(`删除规则失败：${err}`);
    }
  });

  const actionDisabled = !canManage || Boolean(permissionMessage);
  const queryAuthorizationError = isAuthorizationError(rulesQuery.error);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无告警规则访问权限，请联系管理员授予工作空间/项目范围。"
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      {contextHolder}
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            告警规则治理
          </Typography.Title>
          <Typography.Text type="secondary">创建、查看和删除告警规则。</Typography.Text>
        </div>
        <Space wrap>
          <Button onClick={() => void navigate('/observability/alerts')}>告警中心</Button>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
        </Space>
      </Space>

      <Card title="新增规则">
        {!canManage ? (
          <Alert
            type="info"
            showIcon
            style={{ marginBottom: 12 }}
            message="当前为只读模式"
            description="你可以查看现有规则，但无法新增或删除规则。"
          />
        ) : null}
        {permissionMessage ? (
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 12 }}
            message="权限已回收"
            description={permissionMessage}
          />
        ) : null}
        <AlertRuleForm
          loading={createMutation.isPending}
          disabled={actionDisabled}
          onSubmit={(payload) => createMutation.mutate(payload)}
        />
      </Card>

      <Card title="规则列表">
        {rulesQuery.error && !queryAuthorizationError ? (
          <Alert
            type="error"
            showIcon
            message="规则列表加载失败"
            description={normalizeApiError(rulesQuery.error, '规则列表加载失败')}
          />
        ) : null}
        {queryAuthorizationError ? (
          <Alert
            type="warning"
            showIcon
            message="权限已变更"
            description={normalizeApiError(
              rulesQuery.error,
              '当前账号的规则访问权限可能已被回收，请刷新页面或重新登录后重试。'
            )}
          />
        ) : null}
        <Table
          rowKey={(row) => `${row.id}`}
          loading={rulesQuery.isFetching || deleteMutation.isPending}
          dataSource={rulesQuery.data?.items ?? []}
          pagination={{ pageSize: 10 }}
          columns={[
            { title: '名称', dataIndex: 'name', key: 'name' },
            { title: '级别', dataIndex: 'severity', key: 'severity', width: 120 },
            { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
            { title: '表达式', dataIndex: 'conditionExpression', key: 'conditionExpression' },
            {
              title: '操作',
              key: 'actions',
              width: 120,
              render: (_value, record) => (
                <Button
                  danger
                  size="small"
                  disabled={actionDisabled}
                  onClick={() => deleteMutation.mutate(record.id)}
                >
                  删除
                </Button>
              )
            }
          ]}
        />
      </Card>
    </Space>
  );
};
