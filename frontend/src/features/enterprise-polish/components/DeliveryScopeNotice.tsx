import { Alert } from 'antd';

export const DeliveryScopeNotice = ({ text }: { text: string }) => (
  <Alert type="warning" showIcon message={text} />
);
