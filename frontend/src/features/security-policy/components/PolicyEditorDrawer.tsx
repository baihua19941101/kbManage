import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Select, Space, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createSecurityPolicy,
  updateSecurityPolicy
} from '@/services/securityPolicy';
import type {
  CreateSecurityPolicyRequestDTO,
  SecurityPolicyDTO,
  UpdateSecurityPolicyRequestDTO
} from '@/services/api/types';

type PolicyEditorDrawerProps = {
  open: boolean;
  policy?: SecurityPolicyDTO;
  onClose: () => void;
  onSuccess?: (policy: SecurityPolicyDTO) => void;
};

type PolicyFormValues = {
  name: string;
  scopeLevel: SecurityPolicyDTO['scopeLevel'];
  category: SecurityPolicyDTO['category'];
  defaultEnforcementMode: SecurityPolicyDTO['defaultEnforcementMode'];
  riskLevel?: SecurityPolicyDTO['riskLevel'];
  status?: SecurityPolicyDTO['status'];
  ruleTemplateText: string;
};

const parseRuleTemplate = (content: string): Record<string, unknown> => {
  const trimmed = content.trim();
  if (!trimmed) {
    return {};
  }

  try {
    const parsed = JSON.parse(trimmed) as unknown;
    if (typeof parsed === 'object' && parsed !== null && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>;
    }
  } catch {
    throw new Error('规则模板必须是 JSON 对象');
  }

  throw new Error('规则模板必须是 JSON 对象');
};

const serializeRuleTemplate = (value?: Record<string, unknown>) => {
  if (!value || Object.keys(value).length === 0) {
    return '{\n  "rules": []\n}';
  }

  return JSON.stringify(value, null, 2);
};

export const PolicyEditorDrawer = ({ open, policy, onClose, onSuccess }: PolicyEditorDrawerProps) => {
  const [form] = Form.useForm<PolicyFormValues>();
  const isEdit = Boolean(policy);

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      name: policy?.name,
      scopeLevel: policy?.scopeLevel ?? 'platform',
      category: policy?.category ?? 'pod-security',
      defaultEnforcementMode: policy?.defaultEnforcementMode ?? 'audit',
      riskLevel: policy?.riskLevel,
      status: policy?.status,
      ruleTemplateText: serializeRuleTemplate(policy?.ruleTemplate)
    });
  }, [form, open, policy]);

  const mutation = useMutation({
    mutationFn: async (values: PolicyFormValues) => {
      const ruleTemplate = parseRuleTemplate(values.ruleTemplateText);

      if (policy) {
        const payload: UpdateSecurityPolicyRequestDTO = {
          name: values.name.trim(),
          defaultEnforcementMode: values.defaultEnforcementMode,
          ruleTemplate,
          status: values.status === 'draft' ? undefined : values.status
        };
        return updateSecurityPolicy(policy.id, payload);
      }

      const payload: CreateSecurityPolicyRequestDTO = {
        name: values.name.trim(),
        scopeLevel: values.scopeLevel,
        category: values.category,
        defaultEnforcementMode: values.defaultEnforcementMode,
        riskLevel: values.riskLevel,
        ruleTemplate
      };

      return createSecurityPolicy(payload);
    },
    onSuccess: (savedPolicy) => {
      message.success(isEdit ? '策略已更新' : '策略已创建');
      onSuccess?.(savedPolicy);
      form.resetFields();
      onClose();
    }
  });

  const handleClose = () => {
    if (mutation.isPending) {
      return;
    }
    form.resetFields();
    onClose();
  };

  return (
    <Drawer
      title={isEdit ? '编辑策略' : '新建策略'}
      open={open}
      width={600}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button type="primary" loading={mutation.isPending} onClick={() => form.submit()}>
            保存策略
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {mutation.error ? (
          <Alert type="error" showIcon message={normalizeApiError(mutation.error, '保存策略失败')} />
        ) : null}

        <Form<PolicyFormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
          <Form.Item label="策略名称" name="name" rules={[{ required: true, message: '请输入策略名称' }]}>
            <Input placeholder="例如：限制高危特权容器" />
          </Form.Item>

          <Form.Item
            label="策略层级"
            name="scopeLevel"
            rules={[{ required: true, message: '请选择策略层级' }]}
          >
            <Select
              disabled={isEdit}
              options={[
                { label: 'platform', value: 'platform' },
                { label: 'workspace', value: 'workspace' },
                { label: 'project', value: 'project' }
              ]}
            />
          </Form.Item>

          <Form.Item label="策略类别" name="category" rules={[{ required: true, message: '请选择策略类别' }]}>
            <Select
              options={[
                { label: 'pod-security', value: 'pod-security' },
                { label: 'image', value: 'image' },
                { label: 'resource', value: 'resource' },
                { label: 'label', value: 'label' },
                { label: 'network', value: 'network' },
                { label: 'admission', value: 'admission' }
              ]}
            />
          </Form.Item>

          <Form.Item
            label="默认执行模式"
            name="defaultEnforcementMode"
            rules={[{ required: true, message: '请选择默认执行模式' }]}
          >
            <Select
              options={[
                { label: 'audit', value: 'audit' },
                { label: 'alert', value: 'alert' },
                { label: 'warn', value: 'warn' },
                { label: 'enforce', value: 'enforce' }
              ]}
            />
          </Form.Item>

          <Form.Item label="风险级别" name="riskLevel">
            <Select
              allowClear
              options={[
                { label: 'low', value: 'low' },
                { label: 'medium', value: 'medium' },
                { label: 'high', value: 'high' },
                { label: 'critical', value: 'critical' }
              ]}
            />
          </Form.Item>

          {isEdit ? (
            <Form.Item label="状态" name="status">
              <Select
                options={[
                  { label: 'active', value: 'active' },
                  { label: 'disabled', value: 'disabled' },
                  { label: 'archived', value: 'archived' }
                ]}
              />
            </Form.Item>
          ) : null}

          <Form.Item
            label="规则模板（JSON）"
            name="ruleTemplateText"
            rules={[{ required: true, message: '请输入规则模板 JSON' }]}
          >
            <Input.TextArea autoSize={{ minRows: 8, maxRows: 14 }} />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};
