import { Table } from 'antd';
import type { DeliveryArtifact } from '@/services/enterprisePolish';

export const DeliveryArtifactCatalog = ({ items }: { items: DeliveryArtifact[] }) => (
  <Table
    rowKey="id"
    pagination={false}
    dataSource={items}
    columns={[
      { title: '标题', dataIndex: 'title' },
      { title: '类型', dataIndex: 'artifactType' },
      { title: '版本范围', dataIndex: 'versionScope' },
      { title: '环境', dataIndex: 'environmentScope' }
    ]}
  />
);
