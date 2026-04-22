import { List } from 'antd';
import type { DeliveryChecklistItem } from '@/services/enterprisePolish';

export const DeliveryChecklistBoard = ({ items }: { items: DeliveryChecklistItem[] }) => (
  <List
    dataSource={items}
    renderItem={(item) => (
      <List.Item>
        {item.checkItem} / {item.status || '未知状态'} / {item.owner || '未分配'}
      </List.Item>
    )}
  />
);
