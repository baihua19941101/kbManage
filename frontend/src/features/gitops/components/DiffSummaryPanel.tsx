import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Space, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { queryKeys } from '@/app/queryClient';
import { normalizeApiError } from '@/services/api/client';
import { getGitOpsDeliveryUnitDiff, type GitOpsDeliveryDiff, type ResourceId } from '@/services/gitops';

type DiffSummaryPanelProps = {
  unitId: ResourceId;
  environment?: string;
  targetGroupId?: number;
};

type DiffItem = NonNullable<GitOpsDeliveryDiff['items']>[number];

const diffTypeColorMap: Record<string, string> = {
  added: 'green',
  modified: 'blue',
  removed: 'red',
  unavailable: 'default'
};

const columns: ColumnsType<DiffItem> = [
  {
    title: '对象',
    dataIndex: 'objectRef',
    key: 'objectRef',
    render: (value?: string) => value || '-'
  },
  {
    title: '环境',
    dataIndex: 'environment',
    key: 'environment',
    render: (value?: string) => value || '-'
  },
  {
    title: '差异类型',
    dataIndex: 'diffType',
    key: 'diffType',
    render: (value?: string) => (
      <Tag color={diffTypeColorMap[value || 'unavailable'] || 'default'}>{value || 'unknown'}</Tag>
    )
  },
  {
    title: '期望摘要',
    dataIndex: 'desiredSummary',
    key: 'desiredSummary',
    render: (value?: string) => value || '-'
  },
  {
    title: '实际摘要',
    dataIndex: 'liveSummary',
    key: 'liveSummary',
    render: (value?: string) => value || '-'
  }
];

export const DiffSummaryPanel = ({ unitId, environment, targetGroupId }: DiffSummaryPanelProps) => {
  const diffQuery = useQuery({
    queryKey: queryKeys.gitops.deliveryUnits(`diff:${unitId}:${environment || 'all'}:${targetGroupId || 'all'}`),
    queryFn: () => getGitOpsDeliveryUnitDiff(unitId, { environment, targetGroupId }),
    meta: { suppressGlobalError: true }
  });

  const summary = diffQuery.data?.summary;
  const items = diffQuery.data?.items || [];

  return (
    <Card size="small" title="差异/漂移摘要" loading={diffQuery.isLoading || diffQuery.isFetching}>
      {diffQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="差异数据加载失败"
          description={normalizeApiError(diffQuery.error, '差异数据加载失败，请稍后重试。')}
          style={{ marginBottom: 12 }}
        />
      ) : null}

      <Space wrap style={{ marginBottom: 12 }}>
        <Tag color="green">新增：{summary?.added ?? 0}</Tag>
        <Tag color="blue">变更：{summary?.modified ?? 0}</Tag>
        <Tag color="red">删除：{summary?.removed ?? 0}</Tag>
        <Tag color="default">不可用：{summary?.unavailable ?? 0}</Tag>
        <Tag>总计：{items.length}</Tag>
      </Space>

      {items.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无差异" />
      ) : (
        <Table<DiffItem>
          rowKey={(record, index) => `${record.objectRef || 'object'}-${index}`}
          columns={columns}
          dataSource={items}
          pagination={false}
        />
      )}
    </Card>
  );
};
