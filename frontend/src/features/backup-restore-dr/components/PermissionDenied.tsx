import { Empty } from 'antd';

type PermissionDeniedProps = {
  description: string;
};

export const PermissionDenied = ({ description }: PermissionDeniedProps) => (
  <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={description} />
);
