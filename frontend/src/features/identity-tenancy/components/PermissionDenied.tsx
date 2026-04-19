import { Alert } from 'antd';

type PermissionDeniedProps = {
  description: string;
};

export const PermissionDenied = ({ description }: PermissionDeniedProps) => (
  <Alert type="info" showIcon message={description} />
);
