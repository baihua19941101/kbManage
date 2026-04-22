import { List } from 'antd';
import type { GovernanceActionItem } from '@/services/enterprisePolish';

export const ActionItemList = ({ items }: { items: GovernanceActionItem[] }) => (
  <List
    dataSource={items}
    renderItem={(item) => (
      <List.Item>
        {item.title} / {item.priority || '未知优先级'} / {item.status || '未知状态'}
      </List.Item>
    )}
  />
);
