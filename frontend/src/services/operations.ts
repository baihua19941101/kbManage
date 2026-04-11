import { fetchJSON } from '@/services/api/client';
import type { OperationDTO } from '@/services/api/types';

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
  resultMessage?: string;
};

type StoredOperationRecord = OperationRecord;

const STORAGE_KEY = 'kbmanage.operations.records';

const nowISO = () => new Date().toISOString();

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const toNonEmptyText = (value: unknown): string | undefined => {
  if (typeof value === 'string') {
    const trimmed = value.trim();
    return trimmed.length > 0 ? trimmed : undefined;
  }

  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(value);
  }

  return undefined;
};

const toISODate = (value: unknown, fallback: string): string => {
  const text = toNonEmptyText(value);
  if (!text) {
    return fallback;
  }

  const parsed = Date.parse(text);
  if (Number.isNaN(parsed)) {
    return fallback;
  }

  return new Date(parsed).toISOString();
};

const sortRecords = (list: OperationRecord[]) =>
  [...list].sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));

const normalizeOperationType = (
  value: unknown,
  fallback: OperationType = 'restart'
): OperationType => {
  const normalized = toNonEmptyText(value)?.toLowerCase();
  if (normalized === 'scale' || normalized === 'restart' || normalized === 'node-maintenance') {
    return normalized;
  }
  return fallback;
};

const normalizeOperationStatus = (
  value: unknown,
  fallback: OperationStatus = 'pending'
): OperationStatus => {
  const normalized = toNonEmptyText(value)?.toLowerCase();
  if (
    normalized === 'pending' ||
    normalized === 'running' ||
    normalized === 'succeeded' ||
    normalized === 'failed'
  ) {
    return normalized;
  }
  return fallback;
};

const normalizeRiskLevel = (
  value: unknown,
  fallback: OperationRiskLevel = 'medium'
): OperationRiskLevel => {
  const normalized = toNonEmptyText(value)?.toLowerCase();
  if (normalized === 'low' || normalized === 'medium' || normalized === 'high') {
    return normalized;
  }
  if (normalized === 'critical') {
    return 'high';
  }
  return fallback;
};

const parseTargetRef = (
  targetRef: string | undefined,
  fallback: OperationTarget
): OperationTarget => {
  if (!targetRef) {
    return fallback;
  }

  const parts = targetRef.split('/').map((part) => part.trim());
  const source = new Map<string, string>();
  for (const part of parts) {
    const [key, ...rest] = part.split(':');
    if (!key || rest.length === 0) {
      continue;
    }
    const value = rest.join(':').trim();
    if (value.length > 0) {
      source.set(key, value);
    }
  }

  return {
    resourceId: source.get('uid') || fallback.resourceId,
    name: source.get('name') || fallback.name,
    resourceType: source.get('kind') || fallback.resourceType,
    cluster: source.get('cluster') || fallback.cluster,
    namespace: source.get('ns') || fallback.namespace
  };
};

const toOperationId = (value: unknown): string => {
  const text = toNonEmptyText(value);
  if (!text) {
    throw new Error('操作提交成功，但响应未返回有效 operation id。');
  }
  return text;
};

const readStoredOperations = (): StoredOperationRecord[] => {
  if (typeof window === 'undefined') {
    return [];
  }

  const raw = window.localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    return [];
  }

  try {
    const parsed = JSON.parse(raw) as unknown;
    if (!Array.isArray(parsed)) {
      return [];
    }

    const items: StoredOperationRecord[] = [];
    for (const entry of parsed) {
      if (!isRecord(entry)) {
        continue;
      }

      const target = entry.target;
      if (!isRecord(target)) {
        continue;
      }

      const id = toNonEmptyText(entry.id);
      const name = toNonEmptyText(target.name);
      const resourceType = toNonEmptyText(target.resourceType);
      const cluster = toNonEmptyText(target.cluster);
      const namespace = toNonEmptyText(target.namespace);

      if (!id || !name || !resourceType || !cluster || !namespace) {
        continue;
      }

      const createdAt = toISODate(entry.createdAt, nowISO());
      const updatedAt = toISODate(entry.updatedAt, createdAt);

      items.push({
        id,
        type: normalizeOperationType(entry.type),
        target: {
          resourceId: toNonEmptyText(target.resourceId) || '',
          name,
          resourceType,
          cluster,
          namespace
        },
        reason: toNonEmptyText(entry.reason) || '',
        riskLevel: normalizeRiskLevel(entry.riskLevel),
        status: normalizeOperationStatus(entry.status),
        createdAt,
        updatedAt,
        scaleReplicas:
          typeof entry.scaleReplicas === 'number' && Number.isFinite(entry.scaleReplicas)
            ? entry.scaleReplicas
            : undefined,
        resultMessage: toNonEmptyText(entry.resultMessage)
      });
    }

    return sortRecords(items);
  } catch {
    return [];
  }
};

const writeStoredOperations = (items: StoredOperationRecord[]) => {
  if (typeof window === 'undefined') {
    return;
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
};

const upsertStoredOperation = (item: OperationRecord) => {
  const current = readStoredOperations();
  const index = current.findIndex((record) => record.id === item.id);
  if (index >= 0) {
    current[index] = item;
  } else {
    current.unshift(item);
  }
  writeStoredOperations(sortRecords(current));
};

const resolveClusterID = (target: OperationTarget): number | undefined => {
  const candidates = [target.cluster, target.resourceId];
  for (const raw of candidates) {
    const text = toNonEmptyText(raw);
    if (!text) {
      continue;
    }
    if (!/^\d+$/.test(text)) {
      continue;
    }
    const parsed = Number(text);
    if (Number.isInteger(parsed) && parsed > 0) {
      return parsed;
    }
  }

  return undefined;
};

const mapToOperationRecord = (
  dto: OperationDTO,
  fallback: Omit<OperationRecord, 'id'> & { id?: string }
): OperationRecord => {
  const id = toNonEmptyText(dto.id) || fallback.id || '';
  if (!id) {
    throw new Error('操作查询返回了无效 id。');
  }

  const createdAt = toISODate(dto.createdAt, fallback.createdAt);
  const updatedAt = toISODate(dto.updatedAt, createdAt);

  return {
    id,
    type: normalizeOperationType(dto.operationType ?? dto.type, fallback.type),
    target: parseTargetRef(toNonEmptyText(dto.targetRef), fallback.target),
    reason: fallback.reason,
    riskLevel: normalizeRiskLevel(dto.riskLevel, fallback.riskLevel),
    status: normalizeOperationStatus(dto.status, fallback.status),
    createdAt,
    updatedAt,
    scaleReplicas: fallback.scaleReplicas,
    resultMessage: toNonEmptyText(dto.resultMessage) || fallback.resultMessage
  };
};

const buildCreateRequestBody = (payload: CreateOperationPayload) => {
  const clusterID = resolveClusterID(payload.target);
  if (!clusterID) {
    throw new Error('目标资源缺少有效 clusterId，无法提交操作。');
  }

  return {
    idempotencyKey: `web-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 10)}`,
    clusterId: clusterID,
    resourceUid: payload.target.resourceId,
    resourceKind: payload.target.resourceType,
    namespace: payload.target.namespace,
    name: payload.target.name,
    operationType: payload.type,
    riskLevel: payload.riskLevel,
    riskConfirmed: true,
    payload: {
      reason: payload.reason,
      scaleReplicas: payload.type === 'scale' ? payload.scaleReplicas : undefined
    }
  };
};

export const listOperations = async (): Promise<OperationRecord[]> => {
  const stored = readStoredOperations();
  if (stored.length === 0) {
    return [];
  }

  const results = await Promise.allSettled(
    stored.map((item) =>
      fetchJSON<OperationDTO>(`/operations/${encodeURIComponent(item.id)}`, {
        method: 'GET'
      })
    )
  );

  let successCount = 0;
  const records: OperationRecord[] = stored.map((item, index) => {
    const result = results[index];
    if (result.status === 'fulfilled') {
      successCount += 1;
      return mapToOperationRecord(result.value, item);
    }

    return {
      ...item,
      resultMessage:
        result.reason instanceof Error
          ? `状态刷新失败：${result.reason.message}`
          : '状态刷新失败：未知错误'
    };
  });

  writeStoredOperations(sortRecords(records));

  if (successCount === 0) {
    throw new Error('操作状态查询失败，请稍后重试。');
  }

  return sortRecords(records);
};

export const createOperation = async (
  payload: CreateOperationPayload
): Promise<OperationRecord> => {
  if (payload.expectedText !== payload.target.name) {
    throw new Error('二次确认未通过：请输入准确的资源名称。');
  }

  const fallbackCreatedAt = nowISO();
  const fallbackRecord: Omit<OperationRecord, 'id'> = {
    type: payload.type,
    target: payload.target,
    reason: payload.reason,
    riskLevel: payload.riskLevel,
    status: 'pending',
    createdAt: fallbackCreatedAt,
    updatedAt: fallbackCreatedAt,
    scaleReplicas: payload.scaleReplicas
  };

  const response = await fetchJSON<OperationDTO>('/operations', {
    method: 'POST',
    body: JSON.stringify(buildCreateRequestBody(payload))
  });

  const record = mapToOperationRecord(response, {
    id: toOperationId(response.id),
    ...fallbackRecord
  });
  upsertStoredOperation(record);

  return record;
};
