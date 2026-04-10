import { Button, Col, Form, Input, Row, Select } from 'antd';

export type ResourceFilterValues = {
  cluster?: string;
  namespace?: string;
  resourceType?: string;
  keyword?: string;
};

type ResourceFiltersProps = {
  values: ResourceFilterValues;
  clusterOptions: string[];
  namespaceOptions: string[];
  resourceTypeOptions: string[];
  onChange: (values: ResourceFilterValues) => void;
  onReset: () => void;
};

export const ResourceFilters = ({
  values,
  clusterOptions,
  namespaceOptions,
  resourceTypeOptions,
  onChange,
  onReset
}: ResourceFiltersProps) => (
  <Form layout="vertical">
    <Row gutter={12}>
      <Col xs={24} sm={12} md={6}>
        <Form.Item label="Cluster">
          <Select
            allowClear
            placeholder="选择集群"
            value={values.cluster}
            options={clusterOptions.map((value) => ({ label: value, value }))}
            onChange={(cluster) => onChange({ ...values, cluster })}
          />
        </Form.Item>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Form.Item label="Namespace">
          <Select
            allowClear
            placeholder="选择命名空间"
            value={values.namespace}
            options={namespaceOptions.map((value) => ({ label: value, value }))}
            onChange={(namespace) => onChange({ ...values, namespace })}
          />
        </Form.Item>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Form.Item label="Resource Type">
          <Select
            allowClear
            placeholder="选择资源类型"
            value={values.resourceType}
            options={resourceTypeOptions.map((value) => ({ label: value, value }))}
            onChange={(resourceType) => onChange({ ...values, resourceType })}
          />
        </Form.Item>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Form.Item label="Keyword">
          <Input
            allowClear
            placeholder="名称/标签关键字"
            value={values.keyword}
            onChange={(event) => onChange({ ...values, keyword: event.target.value })}
          />
        </Form.Item>
      </Col>
    </Row>
    <Row justify="end">
      <Button onClick={onReset}>重置筛选</Button>
    </Row>
  </Form>
);
