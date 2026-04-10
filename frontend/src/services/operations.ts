export type OperationType = 'scale' | 'restart' | 'node-maintenance';
export type OperationStatus = 'pending' | 'running' | 'succeeded' | 'failed';

export type OperationRiskLevel = 'low' | 'medium' | 'high';

export type OperationTarget = {
  resourceId: string;
  name: string;
  resourceType: string;
  cluster: string;
  namespace: string;
};

export type CreateOperationPayload = {
  type: OperationType;
  target: OperationTarget;
  reason: string;
  riskLevel: OperationRiskLevel;
  expectedText: string;
  scaleReplicas?: number;
};

export type OperationRecord = {
  id: string;
  type: OperationType;
  target: OperationTarget;
  reason: string;
  riskLevel: OperationRiskLevel;
  status: OperationStatus;
  createdAt: string;
  updatedAt: string;
  scaleReplicas?: number;
};

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const nowISO = () => new Date().toISOString();

const genId = () =>
  `op-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

const records: OperationRecord[] = [
  {
    id: 'op-seed-1',
    type: 'restart',
    target: {
      resourceId: 'res-1',
      name: 'payment-api',
      resourceType: 'Deployment',
      cluster: 'prod-cn',
      namespace: 'payments'
    },
    reason: '发布后实例异常，执行重启',
    riskLevel: 'medium',
    status: 'succeeded',
    createdAt: '2026-04-09T02:00:00.000Z',
    updatedAt: '2026-04-09T02:01:00.000Z'
  },
  {
    id: 'op-seed-2',
    type: 'node-maintenance',
    target: {
      resourceId: 'res-2',
      name: 'edge-gateway',
      resourceType: 'Service',
      cluster: 'prod-cn',
      namespace: 'gateway'
    },
    reason: '节点内核升级窗口',
    riskLevel: 'high',
    status: 'running',
    createdAt: '2026-04-09T03:20:00.000Z',
    updatedAt: '2026-04-09T03:21:00.000Z'
  }
];

const sortRecords = (list: OperationRecord[]) =>
  [...list].sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));

const updateStatus = (id: string, status: OperationStatus) => {
  const target = records.find((item) => item.id === id);
  if (!target) {
    return;
  }

  target.status = status;
  target.updatedAt = nowISO();
};

const scheduleMockProgress = (id: string) => {
  window.setTimeout(() => {
    updateStatus(id, 'running');
  }, 1200);

  window.setTimeout(() => {
    const finalStatus: OperationStatus = Math.random() > 0.15 ? 'succeeded' : 'failed';
    updateStatus(id, finalStatus);
  }, 3600);
};

export const listOperations = async (): Promise<OperationRecord[]> => {
  await wait(180);
  return sortRecords(records);
};

export const createOperation = async (
  payload: CreateOperationPayload
): Promise<OperationRecord> => {
  await wait(220);

  if (payload.expectedText !== payload.target.name) {
    throw new Error('二次确认未通过：请输入准确的资源名称。');
  }

  const operation: OperationRecord = {
    id: genId(),
    type: payload.type,
    target: payload.target,
    reason: payload.reason,
    riskLevel: payload.riskLevel,
    status: 'pending',
    createdAt: nowISO(),
    updatedAt: nowISO(),
    scaleReplicas: payload.scaleReplicas
  };

  records.unshift(operation);
  scheduleMockProgress(operation.id);

  return operation;
};
