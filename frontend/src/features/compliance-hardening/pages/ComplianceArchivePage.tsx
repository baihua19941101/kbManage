import { useEffect, useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canExportComplianceArchive,
  canReadCompliance,
  useAuthStore
} from '@/features/auth/store';
import { ArchiveExportDrawer } from '@/features/compliance-hardening/components/ArchiveExportDrawer';
import { useArchiveExport } from '@/features/compliance-hardening/hooks/useArchiveExport';
import { archiveStatusColorMap, formatDateTime } from '@/features/compliance-hardening/utils';
import { normalizeApiError } from '@/services/api/client';
import {
  listComplianceArchiveExports,
  type ArchiveExportScope,
  type ArchiveExportStatus,
  type ArchiveExportTask
} from '@/services/compliance';

const FILTER_STORAGE_KEY = 'kbm-compliance-archive-filters';

const exportScopeOptions: Array<{ label: string; value: ArchiveExportScope | '' }> = [
  { label: '全部范围', value: '' },
  { label: '扫描记录', value: 'scans' },
  { label: '失败项', value: 'findings' },
  { label: '趋势', value: 'trends' },
  { label: '审计', value: 'audit' },
  { label: '归档包', value: 'bundle' }
];

const statusOptions: Array<{ label: string; value: ArchiveExportStatus | '' }> = [
  { label: '全部状态', value: '' },
  { label: '待执行', value: 'pending' },
  { label: '执行中', value: 'running' },
  { label: '成功', value: 'succeeded' },
  { label: '失败', value: 'failed' },
  { label: '已过期', value: 'expired' }
];

const columns: ColumnsType<ArchiveExportTask> = [
  {
    title: '导出任务',
    dataIndex: 'id',
    key: 'id'
  },
  {
    title: '范围',
    dataIndex: 'exportScope',
    key: 'exportScope',
    render: (value?: string) => value || '—'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: ArchiveExportTask['status']) =>
      value ? <Tag color={archiveStatusColorMap[value]}>{value}</Tag> : '—'
  },
  {
    title: '申请人',
    dataIndex: 'requestedBy',
    key: 'requestedBy',
    render: (value?: string) => value || '—'
  },
  {
    title: '完成时间',
    dataIndex: 'completedAt',
    key: 'completedAt',
    render: formatDateTime
  },
  {
    title: '产物',
    key: 'artifactRef',
    render: (_, record) =>
      record.artifactRef ? (
        <Typography.Link href={record.artifactRef} target="_blank">
          下载
        </Typography.Link>
      ) : (
        '—'
      )
  }
];

export const ComplianceArchivePage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const canExport = canExportComplianceArchive(user);
  const [exportScope, setExportScope] = useState<ArchiveExportScope | ''>('');
  const [status, setStatus] = useState<ArchiveExportStatus | ''>('');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const createExportMutation = useArchiveExport();

  useEffect(() => {
    if (typeof window === 'undefined') {
      return;
    }
    const raw = window.localStorage.getItem(FILTER_STORAGE_KEY);
    if (!raw) {
      return;
    }
    try {
      const parsed = JSON.parse(raw) as { exportScope?: ArchiveExportScope; status?: ArchiveExportStatus };
      setExportScope(parsed.exportScope || '');
      setStatus(parsed.status || '');
    } catch {
      window.localStorage.removeItem(FILTER_STORAGE_KEY);
    }
  }, []);

  useEffect(() => {
    if (typeof window === 'undefined') {
      return;
    }
    window.localStorage.setItem(FILTER_STORAGE_KEY, JSON.stringify({ exportScope, status }));
  }, [exportScope, status]);

  const queryInput = useMemo(
    () => ({ exportScope: exportScope || undefined, status: status || undefined }),
    [exportScope, status]
  );

  const exportsQuery = useQuery({
    queryKey: ['compliance', 'archive-exports', queryInput],
    enabled: canRead || canExport,
    queryFn: () => listComplianceArchiveExports(queryInput)
  });

  if (!canRead && !canExport) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无归档导出访问权限。" />;
  }

  const tasks = exportsQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          归档导出
        </Typography.Title>
        <Typography.Text type="secondary">
          管理合规扫描、趋势和审计归档导出任务，支持筛选持久化和未授权空态。
        </Typography.Text>
      </div>

      {!canExport ? (
        <Alert type="info" showIcon message="当前不可创建导出任务" description="你可以查看已有归档任务，但无法创建新的导出任务。" />
      ) : null}

      {exportsQuery.error ? (
        <Alert type="error" showIcon message="归档任务加载失败" description={normalizeApiError(exportsQuery.error, '归档任务加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选与操作">
        <Space wrap>
          <Select style={{ width: 180 }} options={exportScopeOptions} value={exportScope} onChange={setExportScope} />
          <Select style={{ width: 180 }} options={statusOptions} value={status} onChange={setStatus} />
          <Button type="primary" disabled={!canExport} onClick={() => setDrawerOpen(true)}>
            新建导出
          </Button>
        </Space>
      </Card>

      <Card size="small" title={`导出任务（${tasks.length}）`}>
        <Table<ArchiveExportTask>
          rowKey="id"
          dataSource={tasks}
          loading={exportsQuery.isLoading || exportsQuery.isFetching}
          columns={columns}
          pagination={{ pageSize: 8 }}
          scroll={{ x: 960 }}
        />
      </Card>

      <ArchiveExportDrawer
        open={drawerOpen}
        readonly={!canExport}
        loading={createExportMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => createExportMutation.mutate(payload, { onSuccess: () => setDrawerOpen(false) })}
      />
    </Space>
  );
};
