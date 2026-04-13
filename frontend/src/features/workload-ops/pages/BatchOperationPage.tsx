import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Space, Table, Tag, Typography } from 'antd';
import { useSearchParams } from 'react-router-dom';
import { queryKeys } from '@/app/queryClient';
import { normalizeApiError } from '@/services/api/client';
import { getBatchOperation } from '@/services/workloadOps';

export const BatchOperationPage = () => {
  const [searchParams] = useSearchParams();
  const batchId = Number(searchParams.get('batchId') || 0);

  const query = useQuery({
    queryKey: queryKeys.workloadOps.batch(batchId),
    queryFn: () => getBatchOperation(batchId),
    enabled: batchId > 0
  });

  const items = useMemo(() => query.data?.items ?? [], [query.data?.items]);

  if (batchId <= 0) {
    return <Empty description="请提供 batchId 参数" />;
  }

  return (
    <Space direction="vertical" style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        批量任务结果
      </Typography.Title>
      {query.error ? <Alert type="error" showIcon message={normalizeApiError(query.error)} /> : null}
      {query.data ? (
        <Card size="small">
          <Space wrap>
            <Tag>任务ID: {query.data.id}</Tag>
            <Tag color="blue">状态: {query.data.status}</Tag>
            <Tag color="green">成功: {query.data.succeededTargets ?? 0}</Tag>
            <Tag color="red">失败: {query.data.failedTargets ?? 0}</Tag>
          </Space>
        </Card>
      ) : null}
      <Card size="small" title="子项结果">
        <Table
          rowKey={(_, index) => String(index)}
          loading={query.isLoading || query.isFetching}
          dataSource={items}
          pagination={false}
          columns={[
            { title: '资源', dataIndex: 'resourceRef', key: 'resourceRef' },
            { title: '状态', dataIndex: 'status', key: 'status' },
            { title: '结果', dataIndex: 'resultMessage', key: 'resultMessage' },
            { title: '失败原因', dataIndex: 'failureReason', key: 'failureReason' }
          ]}
        />
      </Card>
    </Space>
  );
};
