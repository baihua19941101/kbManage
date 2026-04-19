import { useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Segmented, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/identity-tenancy/components/PageHeader';
import { PermissionDenied } from '@/features/identity-tenancy/components/PermissionDenied';
import { SessionRiskDrawer } from '@/features/identity-tenancy/components/SessionRiskDrawer';
import { StatusTag } from '@/features/identity-tenancy/components/status';
import { useIdentityTenancyPermissions } from '@/features/identity-tenancy/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  identityTenancyQueryKeys,
  listSessionRecords,
  revokeSessionRecord,
  sessionQueryScope,
  updatePreferredLoginMode,
  type LoginMode,
  type SessionRecord
} from '@/services/identityTenancy';

const loginModeOptions: Array<{ label: string; value: LoginMode }> = [
  { label: '本地优先', value: 'local' },
  { label: '外部优先', value: 'external' },
  { label: '并存切换', value: 'mixed' }
];

export const SessionGovernancePage = () => {
  const [selectedRecord, setSelectedRecord] = useState<SessionRecord | undefined>();
  const [loginMode, setLoginMode] = useState<LoginMode>('mixed');
  const permissions = useIdentityTenancyPermissions();
  const loginModeMutation = useMutation({ mutationFn: updatePreferredLoginMode });
  const revokeMutation = useMutation({ mutationFn: revokeSessionRecord });
  const sessionsQuery = useQuery({
    queryKey: identityTenancyQueryKeys.sessions(sessionQueryScope({})),
    enabled: permissions.canRead,
    queryFn: () => listSessionRecords({})
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无会话治理访问权限。" />;
  }

  const columns: ColumnsType<SessionRecord> = [
    { title: '用户', key: 'user', render: (_, record) => record.username || record.userId || '—' },
    { title: '身份来源', dataIndex: 'identitySourceId', key: 'identitySourceId', render: (value?: string) => value || '本地' },
    { title: '登录方式', dataIndex: 'loginMethod', key: 'loginMethod', render: (value?: string) => value || '—' },
    { title: '会话状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> },
    { title: '风险等级', dataIndex: 'riskLevel', key: 'riskLevel', render: (value?: string) => <StatusTag value={value} /> },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => setSelectedRecord(record)}>
            查看详情
          </Button>
          <Button
            size="small"
            danger
            disabled={!permissions.canGovernSession}
            onClick={() => revokeMutation.mutate(record.id)}
          >
            回收会话
          </Button>
        </Space>
      )
    }
  ];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="会话治理"
        description="查看外部身份会话、本地兜底会话与权限回收影响。"
        actions={
          <Segmented
            value={loginMode}
            options={loginModeOptions}
            onChange={(value) => {
              const mode = String(value) as LoginMode;
              setLoginMode(mode);
              loginModeMutation.mutate(mode);
            }}
          />
        }
      />

      <Alert
        type="info"
        showIcon
        message="登录方式切换"
        description="切换登录方式不会删除本地账号，但会影响默认登录入口和会话治理策略。"
      />

      {sessionsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="会话视图加载失败"
          description={normalizeApiError(sessionsQuery.error, '会话视图加载失败，请稍后重试。')}
        />
      ) : null}
      {loginModeMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="登录方式切换失败"
          description={normalizeApiError(loginModeMutation.error, '登录方式切换失败，请稍后重试。')}
        />
      ) : null}
      {revokeMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="会话回收失败"
          description={normalizeApiError(revokeMutation.error, '会话回收失败，请稍后重试。')}
        />
      ) : null}
      {!permissions.canGovernSession ? (
        <Alert type="warning" showIcon message="当前账号只有查看权限，无法执行会话回收。" />
      ) : null}

      <Card size="small" title={`会话列表（${sessionsQuery.data?.items.length ?? 0}）`}>
        <Table<SessionRecord>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={sessionsQuery.data?.items || []}
          loading={sessionsQuery.isLoading || sessionsQuery.isFetching}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <SessionRiskDrawer
        open={Boolean(selectedRecord)}
        record={selectedRecord}
        onClose={() => setSelectedRecord(undefined)}
      />
    </Space>
  );
};
