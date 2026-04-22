import { Table } from 'antd';
import type { PermissionTrail } from '@/services/enterprisePolish';

export const PermissionTrailTable = ({ items }: { items: PermissionTrail[] }) => (
  <Table
    rowKey="id"
    pagination={false}
    dataSource={items}
    columns={[
      { title: '主体', dataIndex: 'subjectRef' },
      { title: '变更类型', dataIndex: 'changeType' },
      { title: '变更前', dataIndex: 'beforeState' },
      { title: '变更后', dataIndex: 'afterState' },
      { title: '证据完整度', dataIndex: 'evidenceCompleteness' }
    ]}
  />
);
