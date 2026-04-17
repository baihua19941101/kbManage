import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Input, Select, Space, Statistic, Table, Tag, Typography, message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  canExecuteComplianceScan,
  canManageComplianceBaseline,
  canReadCompliance,
  useAuthStore
} from '@/features/auth/store';
import { ScanProfileDrawer } from '@/features/compliance-hardening/components/ScanProfileDrawer';
import {
  coverageStatusColorMap,
  formatDateTime,
  recordStatusColorMap,
  scanStatusColorMap
} from '@/features/compliance-hardening/utils';
import { normalizeErrorMessage } from '@/app/queryClient';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import {
  executeScanProfile,
  listScanExecutions,
  listScanProfiles,
  type ComplianceScopeType,
  type ScanExecution,
  type ScanProfile
} from '@/services/compliance';

const scopeOptions: Array<{ label: string; value: ComplianceScopeType | '' }> = [
  { label: '全部范围', value: '' },
  { label: '集群', value: 'cluster' },
  { label: '节点', value: 'node' },
  { label: '命名空间', value: 'namespace' },
  { label: '关键资源集', value: 'resource-set' }
];

const profileColumns = (
  onEdit: (profile: ScanProfile) => void,
  onExecute: (profile: ScanProfile) => void,
  readonly: boolean,
  scanReadonly: boolean
): ColumnsType<ScanProfile> => [
  {
    title: '配置',
    key: 'profile',
    render: (_, record) => (
      <Space direction="vertical" size={0}>
        <Typography.Text strong>{record.name}</Typography.Text>
        <Typography.Text type="secondary">
          {record.scopeType} / {record.scheduleMode}
        </Typography.Text>
      </Space>
    )
  },
  {
    title: '范围摘要',
    key: 'scope',
    render: (_, record) => {
      if (record.scopeType === 'node') {
        const selectors = Object.entries(record.nodeSelectors || {})
          .map(([key, value]) => `${key}=${value}`)
          .join(', ');
        return selectors || '全部节点';
      }
      return [record.clusterRefs?.join(','), record.namespaceRefs?.join(','), record.resourceKinds?.join(',')]
        .filter(Boolean)
        .join(' / ') || '—';
    }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: ScanProfile['status']) =>
      value ? <Tag color={recordStatusColorMap[value === 'paused' ? 'disabled' : value === 'archived' ? 'archived' : value === 'draft' ? 'draft' : 'active']}>{value}</Tag> : '—'
  },
  {
    title: '最近执行',
    dataIndex: 'lastRunAt',
    key: 'lastRunAt',
    render: formatDateTime
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Space>
        <Button type="link" disabled={readonly} onClick={() => onEdit(record)}>
          编辑
        </Button>
        <Button type="link" disabled={scanReadonly} onClick={() => onExecute(record)}>
          立即扫描
        </Button>
      </Space>
    )
  }
];

const executionColumns = (navigate: ReturnType<typeof useNavigate>): ColumnsType<ScanExecution> => [
  {
    title: '执行 ID',
    dataIndex: 'id',
    key: 'id',
    render: (value: string) => <Typography.Text code>{value}</Typography.Text>
  },
  {
    title: '基线快照',
    key: 'baselineSnapshot',
    render: (_, record) => record.baselineSnapshot ? `${record.baselineSnapshot.name || '未命名'} / ${record.baselineSnapshot.version || '未标注版本'}` : '—'
  },
  {
    title: '得分',
    dataIndex: 'score',
    key: 'score',
    render: (value?: number) => (typeof value === 'number' ? `${value.toFixed(1)}%` : '—')
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: ScanExecution['status']) => value ? <Tag color={scanStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '覆盖率',
    dataIndex: 'coverageStatus',
    key: 'coverageStatus',
    render: (value?: ScanExecution['coverageStatus']) => value ? <Tag color={coverageStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '开始时间',
    dataIndex: 'startedAt',
    key: 'startedAt',
    render: formatDateTime
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" onClick={() => void navigate(`/compliance-hardening/findings/${record.id}?scanExecutionId=${record.id}`)}>
        查看失败项
      </Button>
    )
  }
];

export const ScanCenterPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const canManageProfile = canManageComplianceBaseline(user);
  const canExecute = canExecuteComplianceScan(user);
  const [scopeType, setScopeType] = useState<ComplianceScopeType | ''>('');
  const [profileDrawerOpen, setProfileDrawerOpen] = useState(false);
  const [editingProfile, setEditingProfile] = useState<ScanProfile>();
  const [executeReason, setExecuteReason] = useState('');

  const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const contextDefaults = useMemo(() => {
    const clusterId = searchParams.get('clusterId') || undefined;
    const namespace = searchParams.get('namespace') || undefined;
    const resourceKind = searchParams.get('resourceKind') || undefined;
    return {
      clusterRefs: clusterId ? [clusterId] : undefined,
      namespaceRefs: namespace ? [namespace] : undefined,
      resourceKinds: resourceKind ? [resourceKind] : undefined,
      scopeType: namespace ? ('namespace' as const) : ('cluster' as const)
    };
  }, [searchParams]);

  const listQuery = useMemo(
    () => ({
      scopeType: scopeType || undefined
    }),
    [scopeType]
  );

  const profilesQuery = useQuery({
    queryKey: ['compliance', 'scan-profiles', listQuery],
    enabled: canRead,
    queryFn: () => listScanProfiles(listQuery)
  });

  const scansQuery = useQuery({
    queryKey: ['compliance', 'scans', listQuery],
    enabled: canRead,
    queryFn: () => listScanExecutions({})
  });

  const executeMutation = useMutation({
    mutationFn: (profile: ScanProfile) => executeScanProfile(profile.id, { reason: executeReason || undefined }),
    onSuccess: () => {
      message.success('扫描已触发');
      void queryClient.invalidateQueries({ queryKey: ['compliance', 'scans'] });
      void queryClient.invalidateQueries({ queryKey: ['compliance', 'scan-profiles'] });
      setExecuteReason('');
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '扫描触发失败'));
    }
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无合规扫描访问权限。" />;
  }

  const profiles = profilesQuery.data?.items || [];
  const scans = scansQuery.data?.items || [];
  const authorizationChanged = isAuthorizationError(profilesQuery.error) || isAuthorizationError(scansQuery.error);

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          扫描中心
        </Typography.Title>
        <Typography.Text type="secondary">
          维护扫描配置并触发按需扫描，查看基线快照、覆盖率和失败项明细。
        </Typography.Text>
      </div>

      {searchParams.get('clusterId') ? (
        <Alert
          type="info"
          showIcon
          message={`已带入上下文：cluster=${searchParams.get('clusterId')}`}
          description="该筛选通常来自集群或资源详情，可直接创建针对该上下文的扫描配置。"
        />
      ) : null}

      {!canManageProfile ? (
        <Alert type="info" showIcon message="当前无法维护扫描配置" description="你可以查看扫描结果，但不能创建或编辑扫描配置。" />
      ) : null}

      {authorizationChanged ? (
        <Alert type="warning" showIcon message="权限已变更" description="当前账号可能已失去读取或执行扫描权限，请刷新或重新登录。" />
      ) : null}

      {(profilesQuery.error || scansQuery.error) && !authorizationChanged ? (
        <Alert
          type="error"
          showIcon
          message="扫描中心加载失败"
          description={normalizeApiError(profilesQuery.error || scansQuery.error, '扫描中心加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="筛选与操作">
        <Space wrap size={16} style={{ width: '100%', justifyContent: 'space-between' }}>
          <Space wrap>
            <Select style={{ width: 180 }} options={scopeOptions} value={scopeType} onChange={setScopeType} />
            <Input style={{ width: 240 }} placeholder="本次扫描原因（可选）" value={executeReason} onChange={(event) => setExecuteReason(event.target.value)} />
            <Button
              type="primary"
              disabled={!canManageProfile || authorizationChanged}
              onClick={() => {
                setEditingProfile(undefined);
                setProfileDrawerOpen(true);
              }}
            >
              新建扫描配置
            </Button>
          </Space>
          <Space size={24} wrap>
            <Statistic title="扫描配置" value={profiles.length} />
            <Statistic title="最近执行" value={scans.length} />
            <Statistic title="成功执行" value={scans.filter((item) => item.status === 'succeeded').length} />
          </Space>
        </Space>
      </Card>

      <Card size="small" title={`扫描配置（${profiles.length}）`}>
        <Table<ScanProfile>
          rowKey="id"
          dataSource={profiles}
          columns={profileColumns(
            (profile) => {
              setEditingProfile(profile);
              setProfileDrawerOpen(true);
            },
            (profile) => executeMutation.mutate(profile),
            !canManageProfile || authorizationChanged,
            !canExecute || authorizationChanged || executeMutation.isPending
          )}
          loading={profilesQuery.isLoading || profilesQuery.isFetching}
          pagination={{ pageSize: 6 }}
          scroll={{ x: 960 }}
        />
      </Card>

      <Card size="small" title={`最近扫描（${scans.length}）`}>
        <Table<ScanExecution>
          rowKey="id"
          dataSource={scans}
          columns={executionColumns(navigate)}
          loading={scansQuery.isLoading || scansQuery.isFetching}
          pagination={{ pageSize: 6 }}
          scroll={{ x: 1040 }}
        />
      </Card>

      <ScanProfileDrawer
        open={profileDrawerOpen}
        profile={editingProfile}
        defaults={editingProfile ? undefined : contextDefaults}
        readonly={!canManageProfile || authorizationChanged}
        onClose={() => setProfileDrawerOpen(false)}
      />
    </Space>
  );
};
