import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Input, Space, Typography, message } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { canManageObservability, canReadObservability, useAuthStore } from '@/features/auth/store';
import { AlertDetailDrawer } from '@/features/observability/components/AlertDetailDrawer';
import { AlertTable } from '@/features/observability/components/AlertTable';
import { ApiError, normalizeApiError } from '@/services/api/client';
import type { ObservabilityAlertDTO } from '@/services/api/types';
import {
  acknowledgeAlert,
  createAlertHandlingRecord,
  listAlerts
} from '@/services/observability/alerts';

const isAuthorizationError = (error: unknown): boolean =>
  error instanceof ApiError && (error.status === 401 || error.status === 403);

export const AlertCenterPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadObservability(user);
  const canManage = canManageObservability(user);
  const [selectedAlert, setSelectedAlert] = useState<ObservabilityAlertDTO>();
  const [handlingNote, setHandlingNote] = useState('');
  const [permissionMessage, setPermissionMessage] = useState<string>();
  const [msgApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const status = searchParams.get('status') ?? undefined;

  const alertsQuery = useQuery({
    queryKey: ['observability', 'alerts', { status }],
    queryFn: () => listAlerts({ status: status as never }),
    enabled: canRead
  });

  const acknowledgeMutation = useMutation({
    mutationFn: (alertId: string | number) => acknowledgeAlert(alertId, handlingNote || 'acknowledged'),
    onSuccess: async () => {
      setPermissionMessage(undefined);
      await queryClient.invalidateQueries({ queryKey: ['observability', 'alerts'] });
      msgApi.success('告警已确认');
    },
    onError: (err) => {
      if (isAuthorizationError(err)) {
        setPermissionMessage(
          normalizeApiError(err, '权限已回收，无法继续确认告警。请刷新页面后重试。')
        );
        return;
      }
      msgApi.error(`确认失败：${err}`);
    }
  });

  const handlingMutation = useMutation({
    mutationFn: (alertId: string | number) =>
      createAlertHandlingRecord(alertId, {
        actionType: 'note',
        content: handlingNote || 'manual handling note'
      }),
    onSuccess: () => {
      setPermissionMessage(undefined);
      msgApi.success('处理记录已写入');
    },
    onError: (err) => {
      if (isAuthorizationError(err)) {
        setPermissionMessage(
          normalizeApiError(err, '权限已回收，无法继续写入处理记录。请刷新页面后重试。')
        );
        return;
      }
      msgApi.error(`写入记录失败：${err}`);
    }
  });

  const actionDisabled = !canManage || Boolean(permissionMessage);
  const queryAuthorizationError = isAuthorizationError(alertsQuery.error);

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无告警中心访问权限，请联系管理员授予工作空间/项目范围。"
      />
    );
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      {contextHolder}
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            告警中心
          </Typography.Title>
          <Typography.Text type="secondary">
            统一查看告警状态并执行确认与处理记录。
          </Typography.Text>
        </div>
        <Space wrap>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
          <Button onClick={() => void navigate('/observability/alert-rules')}>规则治理</Button>
          <Button onClick={() => void navigate('/observability/silences')}>静默窗口</Button>
        </Space>
      </Space>

      <Card>
        {!canManage ? (
          <Alert
            type="info"
            showIcon
            style={{ marginBottom: 12 }}
            message="当前为只读模式"
            description="你可以查看告警与上下文，但无法确认告警或写入处理记录。"
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
        <Space wrap>
          <Input.TextArea
            style={{ minWidth: 360 }}
            rows={2}
            placeholder="输入处理说明（用于确认/记录）"
            value={handlingNote}
            disabled={actionDisabled}
            onChange={(e) => setHandlingNote(e.target.value)}
          />
        </Space>
      </Card>

      <Card>
        {alertsQuery.error && !queryAuthorizationError ? (
          <Alert
            type="error"
            showIcon
            style={{ marginBottom: 12 }}
            message="告警列表加载失败"
            description={normalizeApiError(alertsQuery.error, '告警列表加载失败')}
          />
        ) : null}
        {queryAuthorizationError ? (
          <Alert
            type="warning"
            showIcon
            style={{ marginBottom: 12 }}
            message="权限已变更"
            description={normalizeApiError(
              alertsQuery.error,
              '当前账号的告警访问权限可能已被回收，请刷新页面或重新登录后重试。'
            )}
          />
        ) : null}
        <AlertTable
          loading={alertsQuery.isFetching || acknowledgeMutation.isPending}
          items={alertsQuery.data?.items ?? []}
          onViewDetail={(item) => setSelectedAlert(item)}
          acknowledgeDisabled={actionDisabled}
          onAcknowledge={
            actionDisabled
              ? undefined
              : (item) => {
                  acknowledgeMutation.mutate(item.id);
                  handlingMutation.mutate(item.id);
                }
          }
        />
      </Card>

      <AlertDetailDrawer
        open={Boolean(selectedAlert)}
        alert={selectedAlert}
        onClose={() => setSelectedAlert(undefined)}
      />
    </Space>
  );
};
