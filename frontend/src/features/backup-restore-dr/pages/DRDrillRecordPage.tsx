import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Input, List, Space, Typography } from 'antd';
import { useParams } from 'react-router-dom';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { StatusTag } from '@/features/backup-restore-dr/components/status';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  backupRestoreQueryKeys,
  getDrillRecordDetail
} from '@/services/backupRestore';

export const DRDrillRecordPage = () => {
  const permissions = useBackupRestorePermissions();
  const params = useParams<{ recordId?: string }>();
  const initialRecordId = params.recordId || 'record-001';
  const [recordIdInput, setRecordIdInput] = useState(initialRecordId);
  const [recordId, setRecordId] = useState(initialRecordId);
  const recordQuery = useQuery({
    queryKey: backupRestoreQueryKeys.drillRecord(recordId),
    enabled: permissions.canRead && recordId.trim().length > 0,
    queryFn: () => getDrillRecordDetail(recordId)
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无灾备演练记录访问权限。" />;
  }

  const record = recordQuery.data;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="灾备演练记录"
        description="按记录编号查看演练过程、步骤结果、验证清单完成度和异常记录。"
      />

      <Card size="small" title="记录查询">
        <Space wrap>
          <Input
            value={recordIdInput}
            onChange={(event) => setRecordIdInput(event.target.value)}
            placeholder="输入演练记录 ID"
            style={{ width: 260 }}
          />
          <Button type="primary" onClick={() => setRecordId(recordIdInput.trim())}>
            查询记录
          </Button>
        </Space>
      </Card>

      {recordQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="演练记录加载失败"
          description={normalizeApiError(recordQuery.error, '演练记录加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="记录概览" loading={recordQuery.isLoading || recordQuery.isFetching}>
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          <Typography.Text>
            状态：<StatusTag value={record?.status} />
          </Typography.Text>
          <Typography.Text>
            实际 RPO / RTO：{record?.actualRpoMinutes ?? '—'} 分钟 / {record?.actualRtoMinutes ?? '—'} 分钟
          </Typography.Text>
          <Typography.Text>执行人：{record?.executedBy || '—'}</Typography.Text>
          <Typography.Text>异常记录：{record?.incidentNotes || '无'}</Typography.Text>
          <div>
            <Typography.Text strong>步骤结果</Typography.Text>
            <List
              size="small"
              dataSource={record?.stepResults || []}
              locale={{ emptyText: '暂无步骤结果' }}
              renderItem={(item) => <List.Item>{item}</List.Item>}
            />
          </div>
          <div>
            <Typography.Text strong>验证结果</Typography.Text>
            <List
              size="small"
              dataSource={record?.validationResults || []}
              locale={{ emptyText: '暂无验证结果' }}
              renderItem={(item) => <List.Item>{item}</List.Item>}
            />
          </div>
        </Space>
      </Card>
    </Space>
  );
};
