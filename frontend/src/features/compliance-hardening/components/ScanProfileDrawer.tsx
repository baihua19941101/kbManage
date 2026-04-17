import { useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Drawer, Form, Input, Select, Space, message } from 'antd';
import {
  createScanProfile,
  listComplianceBaselines,
  updateScanProfile,
  type ComplianceScopeType,
  type ScanProfile,
  type ScheduleMode
} from '@/services/compliance';
import { normalizeErrorMessage } from '@/app/queryClient';
import { parseKeyValueLines, stringifyKeyValueLines } from '@/features/compliance-hardening/utils';

const scopeOptions: { label: string; value: ComplianceScopeType }[] = [
  { label: '集群', value: 'cluster' },
  { label: '节点', value: 'node' },
  { label: '命名空间', value: 'namespace' },
  { label: '关键资源集', value: 'resource-set' }
];

const scheduleOptions: { label: string; value: ScheduleMode }[] = [
  { label: '按需', value: 'manual' },
  { label: '计划', value: 'scheduled' }
];

type FormValues = {
  name: string;
  baselineId: string;
  scopeType: ComplianceScopeType;
  clusterRefs?: string[];
  namespaceRefs?: string[];
  resourceKinds?: string[];
  nodeSelectorText?: string;
  scheduleMode: ScheduleMode;
  cronExpression?: string;
};

type ScanProfileDrawerProps = {
  open: boolean;
  profile?: ScanProfile;
  defaults?: Partial<FormValues>;
  readonly?: boolean;
  onClose: () => void;
};

export const ScanProfileDrawer = ({
  open,
  profile,
  defaults,
  readonly,
  onClose
}: ScanProfileDrawerProps) => {
  const [form] = Form.useForm<FormValues>();
  const queryClient = useQueryClient();

  const baselinesQuery = useQuery({
    queryKey: ['compliance', 'baselines', 'active-options'],
    queryFn: () => listComplianceBaselines({}),
    enabled: open
  });

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      name: profile?.name || defaults?.name || '',
      baselineId: profile?.baselineId || defaults?.baselineId || '',
      scopeType: profile?.scopeType || defaults?.scopeType || 'cluster',
      clusterRefs: profile?.clusterRefs || defaults?.clusterRefs || [],
      namespaceRefs: profile?.namespaceRefs || defaults?.namespaceRefs || [],
      resourceKinds: profile?.resourceKinds || defaults?.resourceKinds || [],
      nodeSelectorText: stringifyKeyValueLines(profile?.nodeSelectors) || defaults?.nodeSelectorText || '',
      scheduleMode: profile?.scheduleMode || defaults?.scheduleMode || 'manual',
      cronExpression: profile?.cronExpression || defaults?.cronExpression || ''
    });
  }, [defaults, form, open, profile]);

  const mutation = useMutation({
    mutationFn: (values: FormValues) => {
      const payload = {
        name: values.name,
        baselineId: values.baselineId,
        scopeType: values.scopeType,
        clusterRefs: values.clusterRefs?.filter(Boolean),
        namespaceRefs: values.namespaceRefs?.filter(Boolean),
        resourceKinds: values.resourceKinds?.filter(Boolean),
        nodeSelectors: parseKeyValueLines(values.nodeSelectorText),
        scheduleMode: values.scheduleMode,
        cronExpression: values.scheduleMode === 'scheduled' ? values.cronExpression : undefined
      };

      if (profile?.id) {
        return updateScanProfile(profile.id, {
          name: payload.name,
          nodeSelectors: payload.nodeSelectors,
          cronExpression: payload.cronExpression
        });
      }

      return createScanProfile(payload);
    },
    onSuccess: () => {
      message.success(profile?.id ? '扫描配置已更新' : '扫描配置已创建');
      void queryClient.invalidateQueries({ queryKey: ['compliance', 'scan-profiles'] });
      onClose();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '扫描配置保存失败'));
    }
  });

  const scopeType = Form.useWatch('scopeType', form) || 'cluster';
  const scheduleMode = Form.useWatch('scheduleMode', form) || 'manual';

  return (
    <Drawer
      title={profile?.id ? '编辑扫描配置' : '新建扫描配置'}
      width={480}
      open={open}
      onClose={onClose}
      destroyOnClose
    >
      <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
        <Form.Item label="配置名称" name="name" rules={[{ required: true, message: '请输入配置名称' }]}>
          <Input disabled={readonly} placeholder="例如：生产集群节点基线巡检" />
        </Form.Item>
        <Form.Item label="基线标准" name="baselineId" rules={[{ required: true, message: '请选择基线标准' }]}>
          <Select
            disabled={readonly || Boolean(profile?.id)}
            loading={baselinesQuery.isLoading}
            options={(baselinesQuery.data?.items || []).map((item) => ({
              label: `${item.name} (${item.version})`,
              value: item.id
            }))}
          />
        </Form.Item>
        <Form.Item label="扫描范围" name="scopeType" rules={[{ required: true }]}> 
          <Select disabled={readonly || Boolean(profile?.id)} options={scopeOptions} />
        </Form.Item>
        <Form.Item label="集群" name="clusterRefs">
          <Select mode="tags" tokenSeparators={[',']} disabled={readonly} placeholder="输入 cluster id" />
        </Form.Item>
        {scopeType === 'namespace' || scopeType === 'resource-set' ? (
          <Form.Item label="命名空间" name="namespaceRefs">
            <Select mode="tags" tokenSeparators={[',']} disabled={readonly} placeholder="例如 kube-system" />
          </Form.Item>
        ) : null}
        {scopeType === 'resource-set' ? (
          <Form.Item label="关键资源类型" name="resourceKinds">
            <Select
              mode="tags"
              tokenSeparators={[',']}
              disabled={readonly}
              placeholder="例如 Deployment, RoleBinding"
            />
          </Form.Item>
        ) : null}
        {scopeType === 'node' ? (
          <Form.Item label="节点选择器" name="nodeSelectorText">
            <Input.TextArea
              disabled={readonly}
              rows={4}
              placeholder={'一行一个 key=value\n例如\nnode-role.kubernetes.io/control-plane=true'}
            />
          </Form.Item>
        ) : null}
        <Form.Item label="执行模式" name="scheduleMode" rules={[{ required: true }]}> 
          <Select disabled={readonly || Boolean(profile?.id)} options={scheduleOptions} />
        </Form.Item>
        {scheduleMode === 'scheduled' ? (
          <Form.Item label="Cron 表达式" name="cronExpression" rules={[{ required: true, message: '请输入 cron 表达式' }]}>
            <Input disabled={readonly} placeholder="例如：0 3 * * *" />
          </Form.Item>
        ) : null}
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={mutation.isPending} disabled={readonly}>
            保存
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
