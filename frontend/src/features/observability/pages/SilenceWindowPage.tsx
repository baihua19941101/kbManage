import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table, Typography, message } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { NotificationTargetForm } from '@/features/observability/components/NotificationTargetForm';
import { SilenceWindowForm } from '@/features/observability/components/SilenceWindowForm';
import {
  createNotificationTarget,
  listNotificationTargets
} from '@/services/observability/notificationTargets';
import { cancelSilence, createSilence, listSilences } from '@/services/observability/silences';

export const SilenceWindowPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [msgApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const silencesQuery = useQuery({
    queryKey: ['observability', 'silences'],
    queryFn: () => listSilences()
  });
  const targetsQuery = useQuery({
    queryKey: ['observability', 'notification-targets'],
    queryFn: listNotificationTargets
  });

  const createSilenceMutation = useMutation({
    mutationFn: createSilence,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['observability', 'silences'] });
      msgApi.success('静默窗口已创建');
    },
    onError: (err) => msgApi.error(`创建静默窗口失败：${err}`)
  });

  const cancelSilenceMutation = useMutation({
    mutationFn: cancelSilence,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['observability', 'silences'] });
      msgApi.success('静默窗口已取消');
    },
    onError: (err) => msgApi.error(`取消静默窗口失败：${err}`)
  });

  const createTargetMutation = useMutation({
    mutationFn: createNotificationTarget,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['observability', 'notification-targets'] });
      msgApi.success('通知目标已创建');
    },
    onError: (err) => msgApi.error(`创建通知目标失败：${err}`)
  });

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      {contextHolder}
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            通知目标与静默窗口
          </Typography.Title>
          <Typography.Text type="secondary">
            管理值班通知目标与告警静默窗口策略。
          </Typography.Text>
        </div>
        <Space wrap>
          <Button onClick={() => void navigate('/observability/alerts')}>告警中心</Button>
          <Button onClick={() => void navigate('/observability/alert-rules')}>规则治理</Button>
          <Button onClick={() => void navigate(`/observability${location.search}`)}>总览</Button>
        </Space>
      </Space>

      <Card title="新增通知目标">
        <NotificationTargetForm
          loading={createTargetMutation.isPending}
          onSubmit={(payload) => createTargetMutation.mutate(payload)}
        />
      </Card>

      <Card title="通知目标列表">
        {targetsQuery.error ? (
          <Alert type="error" showIcon message="通知目标加载失败" description={`${targetsQuery.error}`} />
        ) : null}
        <Table
          rowKey={(row) => `${row.id}`}
          dataSource={targetsQuery.data?.items ?? []}
          loading={targetsQuery.isFetching}
          pagination={{ pageSize: 8 }}
          columns={[
            { title: '名称', dataIndex: 'name', key: 'name' },
            { title: '类型', dataIndex: 'targetType', key: 'targetType', width: 120 },
            { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
            { title: '配置引用', dataIndex: 'configRef', key: 'configRef' }
          ]}
        />
      </Card>

      <Card title="新增静默窗口">
        <SilenceWindowForm
          loading={createSilenceMutation.isPending}
          onSubmit={(payload) => createSilenceMutation.mutate(payload)}
        />
      </Card>

      <Card title="静默窗口列表">
        {silencesQuery.error ? (
          <Alert type="error" showIcon message="静默窗口加载失败" description={`${silencesQuery.error}`} />
        ) : null}
        <Table
          rowKey={(row) => `${row.id}`}
          dataSource={silencesQuery.data?.items ?? []}
          loading={silencesQuery.isFetching || cancelSilenceMutation.isPending}
          pagination={{ pageSize: 8 }}
          columns={[
            { title: '名称', dataIndex: 'name', key: 'name' },
            { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
            { title: '开始时间', dataIndex: 'startsAt', key: 'startsAt', width: 220 },
            { title: '结束时间', dataIndex: 'endsAt', key: 'endsAt', width: 220 },
            { title: '原因', dataIndex: 'reason', key: 'reason' },
            {
              title: '操作',
              key: 'actions',
              width: 120,
              render: (_value, record) => (
                <Button size="small" danger onClick={() => cancelSilenceMutation.mutate(record.id)}>
                  取消
                </Button>
              )
            }
          ]}
        />
      </Card>
    </Space>
  );
};
