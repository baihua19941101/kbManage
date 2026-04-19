import { Alert } from 'antd';

export const PermissionDenied = ({ description }: { description: string }) => (
  <Alert type="info" showIcon message="访问受限" description={description} />
);
