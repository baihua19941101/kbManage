import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Statistic, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canManageComplianceBaseline,
  canReadCompliance,
  useAuthStore
} from '@/features/auth/store';
import { BaselineFormDrawer } from '@/features/compliance-hardening/components/BaselineFormDrawer';
import { formatDateTime, recordStatusColorMap } from '@/features/compliance-hardening/utils';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import {
  listComplianceBaselines,
  type ComplianceBaseline,
  type ComplianceRecordStatus,
  type ComplianceStandardType
} from '@/services/compliance';

const standardOptions: Array<{ label: string; value: ComplianceStandardType | '' }> = [
  { label: '全部标准', value: '' },
  { label: 'CIS', value: 'cis' },
  { label: 'STIG', value: 'stig' },
  { label: '平台基线', value: 'platform-baseline' }
];

const statusOptions: Array<{ label: string; value: ComplianceRecordStatus | '' }> = [
  { label: '全部状态', value: '' },
  { label: '草稿', value: 'draft' },
  { label: '启用', value: 'active' },
  { label: '停用', value: 'disabled' },
  { label: '归档', value: 'archived' }
];

const columns = (onEdit: (item: ComplianceBaseline) => void, readonly: boolean): ColumnsType<ComplianceBaseline> => [
  {
    title: '基线',
    key: 'name',
    render: (_, record) => (
      <Space direction="vertical" size={0}>
        <Typography.Text strong>{record.name}</Typography.Text>
        <Typography.Text type="secondary">{record.description || '未填写说明'}</Typography.Text>
      </Space>
    )
  },
  { title: '标准', dataIndex: 'standardType', key: 'standardType' },
  { title: '版本', dataIndex: 'version', key: 'version' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: ComplianceBaseline['status']) =>
      value ? <Tag color={recordStatusColorMap[value]}>{value}</Tag> : '—'
  },
  { title: '规则数', dataIndex: 'ruleCount', key: 'ruleCount', render: (value?: number) => value ?? '—' },
  {
    title: '操作',
    key: 'actions',
    render: (_, record) => (
      <Button type="link" disabled={readonly} onClick={() => onEdit(record)}>
        编辑
      </Button>
    )
  }
];

export const ComplianceBaselinePage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const canManage = canManageComplianceBaseline(user);
  const [standardType, setStandardType] = useState<ComplianceStandardType | ''>('');
  const [status, setStatus] = useState<ComplianceRecordStatus | ''>('');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingBaseline, setEditingBaseline] = useState<ComplianceBaseline>();

  const listQuery = useMemo(
    () => ({
      standardType: standardType || undefined,
      status: status || undefined
    }),
    [standardType, status]
  );

  const baselinesQuery = useQuery({
    queryKey: ['compliance', 'baselines', listQuery],
    enabled: canRead,
    queryFn: () => listComplianceBaselines(listQuery)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无合规与加固访问权限。" />;
  }

  const baselines = baselinesQuery.data?.items || [];
  const authorizationChanged = isAuthorizationError(baselinesQuery.error);
  const activeCount = baselines.filter((item) => item.status === 'active').length;
  const latestVersion = baselines[0]?.version;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          合规基线管理
        </Typography.Title>
        <Typography.Text type="secondary">
          统一维护 CIS、STIG 和平台基线版本，为扫描配置和历史比对提供稳定快照口径。
        </Typography.Text>
      </div>

      {!canManage ? (
        <Alert type="info" showIcon message="当前为只读模式" description="你可以查看基线，但无法创建或更新。" />
      ) : null}

      {authorizationChanged ? (
        <Alert type="warning" showIcon message="权限已变更" description="当前会话可能已失去合规读取权限，请刷新或重新登录。" />
      ) : null}

      {baselinesQuery.error && !authorizationChanged ? (
        <Alert type="error" showIcon message="基线列表加载失败" description={normalizeApiError(baselinesQuery.error, '基线列表加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选与概览" extra={<Typography.Text type="secondary">更新时间：{formatDateTime(new Date().toISOString())}</Typography.Text>}>
        <Space wrap size={24} style={{ width: '100%', justifyContent: 'space-between' }}>
          <Space wrap>
            <Select style={{ width: 180 }} options={standardOptions} value={standardType} onChange={setStandardType} />
            <Select style={{ width: 180 }} options={statusOptions} value={status} onChange={setStatus} />
            <Button
              type="primary"
              disabled={!canManage || authorizationChanged}
              onClick={() => {
                setEditingBaseline(undefined);
                setDrawerOpen(true);
              }}
            >
              新建基线
            </Button>
          </Space>
          <Space size={24} wrap>
            <Statistic title="基线总数" value={baselines.length} />
            <Statistic title="启用中" value={activeCount} />
            <Statistic title="当前首屏版本" value={latestVersion || '—'} />
          </Space>
        </Space>
      </Card>

      <Card size="small" title={`基线列表（${baselines.length}）`}>
        <Table<ComplianceBaseline>
          rowKey="id"
          loading={baselinesQuery.isLoading || baselinesQuery.isFetching}
          dataSource={baselines}
          columns={columns(
            (item) => {
              setEditingBaseline(item);
              setDrawerOpen(true);
            },
            !canManage || authorizationChanged
          )}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <BaselineFormDrawer
        open={drawerOpen}
        baseline={editingBaseline}
        readonly={!canManage || authorizationChanged}
        onClose={() => setDrawerOpen(false)}
      />
    </Space>
  );
};
