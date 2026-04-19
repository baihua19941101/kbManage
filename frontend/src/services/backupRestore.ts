import { buildScopeQueryKey, fetchJSON } from '@/services/api/client';

type UnknownRecord = Record<string, unknown>;
type QueryValue = string | number | boolean | undefined | null;

const isRecord = (value: unknown): value is UnknownRecord =>
  typeof value === 'object' && value !== null;

const pick = (record: UnknownRecord, keys: string[]): unknown => {
  for (const key of keys) {
    if (key in record) {
      return record[key];
    }
  }
  return undefined;
};

const toText = (value: unknown): string | undefined => {
  if (typeof value === 'string') {
    const trimmed = value.trim();
    return trimmed.length > 0 ? trimmed : undefined;
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(value);
  }
  return undefined;
};

const toNumber = (value: unknown): number | undefined => {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === 'string') {
    const parsed = Number(value.trim());
    return Number.isFinite(parsed) ? parsed : undefined;
  }
  return undefined;
};

const toStringArray = (value: unknown): string[] => {
  if (!Array.isArray(value)) {
    return [];
  }
  return value
    .map((item) => toText(item))
    .filter((item): item is string => Boolean(item));
};

const toQueryString = (query: Record<string, QueryValue>) => {
  const params = new URLSearchParams();
  Object.entries(query).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') {
      return;
    }
    params.set(key, String(value));
  });
  return params.toString();
};

const withQuery = (path: string, query: Record<string, QueryValue>) => {
  const queryString = toQueryString(query);
  return queryString ? `${path}?${queryString}` : path;
};

const normalizeListResponse = <T>(
  value: unknown,
  mapper: (record: UnknownRecord, index: number) => T
): { items: T[] } => {
  if (Array.isArray(value)) {
    return { items: value.filter(isRecord).map(mapper) };
  }
  const record = isRecord(value) ? value : {};
  const items = Array.isArray(record.items)
    ? record.items
    : Array.isArray(record.Items)
      ? record.Items
      : [];
  return { items: items.filter(isRecord).map(mapper) };
};

export type BackupPolicyStatus = 'active' | 'paused' | 'draft' | 'failed' | string;
export type RestorePointResult = 'succeeded' | 'partial' | 'failed' | 'running' | string;
export type RestoreJobType =
  | 'in-place'
  | 'cross-cluster'
  | 'selective'
  | 'migration'
  | string;
export type RestoreJobStatus =
  | 'pending'
  | 'running'
  | 'succeeded'
  | 'partial'
  | 'failed'
  | 'blocked'
  | string;
export type DrillStatus =
  | 'draft'
  | 'scheduled'
  | 'running'
  | 'completed'
  | 'failed'
  | string;

export type ScopeSelection = Record<string, unknown>;

export type BackupPolicy = {
  id: string;
  name: string;
  description?: string;
  scopeType?: string;
  scopeRef?: string;
  executionMode?: string;
  scheduleExpression?: string;
  retentionRule?: string;
  consistencyLevel?: string;
  status?: BackupPolicyStatus;
  ownerUserId?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateBackupPolicyPayload = {
  name: string;
  description?: string;
  scopeType: string;
  scopeRef?: string;
  executionMode: string;
  scheduleExpression?: string;
  retentionRule: string;
  consistencyLevel: string;
};

export type RestorePoint = {
  id: string;
  policyId?: string;
  scopeSnapshot?: ScopeSelection;
  backupStartedAt?: string;
  backupCompletedAt?: string;
  durationSeconds?: number;
  result?: RestorePointResult;
  consistencySummary?: string;
  failureReason?: string;
  storageRef?: string;
  expiresAt?: string;
  createdBy?: string;
};

export type RestoreJob = {
  id: string;
  restorePointId?: string;
  jobType?: RestoreJobType;
  sourceEnvironment?: string;
  targetEnvironment?: string;
  scopeSelection?: ScopeSelection;
  conflictSummary?: string;
  consistencyNotice?: string;
  status?: RestoreJobStatus;
  resultSummary?: string;
  failureReason?: string;
  requestedBy?: string;
  startedAt?: string;
  completedAt?: string;
};

export type CreateRestoreJobPayload = {
  restorePointId: string;
  jobType: RestoreJobType;
  sourceEnvironment?: string;
  targetEnvironment: string;
  scopeSelection: ScopeSelection;
};

export type PrecheckResult = {
  status?: string;
  blockers: string[];
  warnings: string[];
  consistencyNotice?: string;
};

export type MigrationPlan = {
  id: string;
  name: string;
  sourceClusterId?: string;
  targetClusterId?: string;
  scopeSelection?: ScopeSelection;
  mappingRules?: ScopeSelection;
  cutoverSteps: string[];
  status?: string;
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateMigrationPlanPayload = {
  name: string;
  sourceClusterId: string;
  targetClusterId: string;
  scopeSelection: ScopeSelection;
  mappingRules?: ScopeSelection;
  cutoverSteps?: string[];
};

export type DRDrillPlan = {
  id: string;
  name: string;
  description?: string;
  scopeSelection?: ScopeSelection;
  rpoTargetMinutes?: number;
  rtoTargetMinutes?: number;
  roleAssignments: string[];
  cutoverProcedure: string[];
  validationChecklist: string[];
  status?: DrillStatus;
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateDRDrillPlanPayload = {
  name: string;
  description?: string;
  scopeSelection: ScopeSelection;
  rpoTargetMinutes: number;
  rtoTargetMinutes: number;
  roleAssignments?: string[];
  cutoverProcedure: string[];
  validationChecklist: string[];
};

export type DRDrillRecord = {
  id: string;
  planId?: string;
  startedAt?: string;
  completedAt?: string;
  actualRpoMinutes?: number;
  actualRtoMinutes?: number;
  status?: DrillStatus;
  stepResults: string[];
  validationResults: string[];
  incidentNotes?: string;
  executedBy?: string;
};

export type DRDrillReport = {
  id: string;
  drillRecordId?: string;
  goalAssessment?: string;
  gapSummary?: string;
  issuesFound: string[];
  improvementActions: string[];
  publishedAt?: string;
  publishedBy?: string;
};

export type BackupAuditEvent = {
  id: string;
  action?: string;
  actorUserId?: string;
  targetType?: string;
  targetRef?: string;
  scopeSnapshot?: ScopeSelection;
  outcome?: string;
  detailSnapshot?: ScopeSelection;
  occurredAt?: string;
};

export type BackupPolicyListQuery = {
  scopeType?: string;
  status?: string;
};

export type RestorePointListQuery = {
  policyId?: string;
  result?: string;
};

export type RestoreJobListQuery = {
  jobType?: string;
  status?: string;
};

export type BackupAuditListQuery = {
  action?: string;
  outcome?: string;
  targetType?: string;
};

export const backupRestoreQueryKeys = {
  all: ['backup-restore'] as const,
  policies: (scope?: string) => ['backup-restore', 'policies', scope ?? 'all'] as const,
  restorePoints: (scope?: string) =>
    ['backup-restore', 'restore-points', scope ?? 'all'] as const,
  restorePointDetail: (restorePointId?: string) =>
    ['backup-restore', 'restore-point-detail', restorePointId ?? 'unknown'] as const,
  restoreJobs: (scope?: string) =>
    ['backup-restore', 'restore-jobs', scope ?? 'all'] as const,
  drillPlans: (scope?: string) => ['backup-restore', 'drill-plans', scope ?? 'all'] as const,
  drillRecord: (recordId?: string) =>
    ['backup-restore', 'drill-record', recordId ?? 'unknown'] as const,
  audit: (scope?: string) => ['backup-restore', 'audit', scope ?? 'all'] as const
};

const mapBackupPolicy = (record: UnknownRecord): BackupPolicy => ({
  id: toText(pick(record, ['id', 'policyId'])) || 'unknown-policy',
  name: toText(pick(record, ['name', 'displayName'])) || '未命名策略',
  description: toText(record.description),
  scopeType: toText(record.scopeType),
  scopeRef: toText(record.scopeRef),
  executionMode: toText(record.executionMode),
  scheduleExpression: toText(record.scheduleExpression),
  retentionRule: toText(record.retentionRule),
  consistencyLevel: toText(record.consistencyLevel),
  status: toText(record.status),
  ownerUserId: toText(record.ownerUserId),
  createdAt: toText(record.createdAt),
  updatedAt: toText(record.updatedAt)
});

const mapRestorePoint = (record: UnknownRecord): RestorePoint => ({
  id: toText(pick(record, ['id', 'restorePointId'])) || 'unknown-restore-point',
  policyId: toText(record.policyId),
  scopeSnapshot: isRecord(record.scopeSnapshot) ? record.scopeSnapshot : undefined,
  backupStartedAt: toText(record.backupStartedAt),
  backupCompletedAt: toText(record.backupCompletedAt),
  durationSeconds: toNumber(record.durationSeconds),
  result: toText(record.result),
  consistencySummary: toText(record.consistencySummary),
  failureReason: toText(record.failureReason),
  storageRef: toText(record.storageRef),
  expiresAt: toText(record.expiresAt),
  createdBy: toText(record.createdBy)
});

const mapRestoreJob = (record: UnknownRecord): RestoreJob => ({
  id: toText(pick(record, ['id', 'jobId'])) || 'unknown-restore-job',
  restorePointId: toText(record.restorePointId),
  jobType: toText(record.jobType),
  sourceEnvironment: toText(record.sourceEnvironment),
  targetEnvironment: toText(record.targetEnvironment),
  scopeSelection: isRecord(record.scopeSelection) ? record.scopeSelection : undefined,
  conflictSummary: toText(record.conflictSummary),
  consistencyNotice: toText(record.consistencyNotice),
  status: toText(record.status),
  resultSummary: toText(record.resultSummary),
  failureReason: toText(record.failureReason),
  requestedBy: toText(record.requestedBy),
  startedAt: toText(record.startedAt),
  completedAt: toText(record.completedAt)
});

const mapPrecheckResult = (value: unknown): PrecheckResult => {
  const record = isRecord(value) ? value : {};
  return {
    status: toText(record.status),
    blockers: toStringArray(record.blockers),
    warnings: toStringArray(record.warnings),
    consistencyNotice: toText(record.consistencyNotice)
  };
};

const mapMigrationPlan = (record: UnknownRecord): MigrationPlan => ({
  id: toText(pick(record, ['id', 'migrationPlanId'])) || 'unknown-migration-plan',
  name: toText(record.name) || '未命名迁移计划',
  sourceClusterId: toText(record.sourceClusterId),
  targetClusterId: toText(record.targetClusterId),
  scopeSelection: isRecord(record.scopeSelection) ? record.scopeSelection : undefined,
  mappingRules: isRecord(record.mappingRules) ? record.mappingRules : undefined,
  cutoverSteps: toStringArray(record.cutoverSteps),
  status: toText(record.status),
  createdBy: toText(record.createdBy),
  createdAt: toText(record.createdAt),
  updatedAt: toText(record.updatedAt)
});

const mapDrillPlan = (record: UnknownRecord): DRDrillPlan => ({
  id: toText(pick(record, ['id', 'planId'])) || 'unknown-drill-plan',
  name: toText(record.name) || '未命名演练计划',
  description: toText(record.description),
  scopeSelection: isRecord(record.scopeSelection) ? record.scopeSelection : undefined,
  rpoTargetMinutes: toNumber(record.rpoTargetMinutes),
  rtoTargetMinutes: toNumber(record.rtoTargetMinutes),
  roleAssignments: toStringArray(record.roleAssignments),
  cutoverProcedure: toStringArray(record.cutoverProcedure),
  validationChecklist: toStringArray(record.validationChecklist),
  status: toText(record.status),
  createdBy: toText(record.createdBy),
  createdAt: toText(record.createdAt),
  updatedAt: toText(record.updatedAt)
});

const mapDrillRecord = (record: UnknownRecord): DRDrillRecord => ({
  id: toText(pick(record, ['id', 'recordId'])) || 'unknown-drill-record',
  planId: toText(record.planId),
  startedAt: toText(record.startedAt),
  completedAt: toText(record.completedAt),
  actualRpoMinutes: toNumber(record.actualRpoMinutes),
  actualRtoMinutes: toNumber(record.actualRtoMinutes),
  status: toText(record.status),
  stepResults: toStringArray(record.stepResults),
  validationResults: toStringArray(record.validationResults),
  incidentNotes: toText(record.incidentNotes),
  executedBy: toText(record.executedBy)
});

const mapDrillReport = (record: UnknownRecord): DRDrillReport => ({
  id: toText(pick(record, ['id', 'reportId'])) || 'unknown-drill-report',
  drillRecordId: toText(record.drillRecordId),
  goalAssessment: toText(record.goalAssessment),
  gapSummary: toText(record.gapSummary),
  issuesFound: toStringArray(record.issuesFound),
  improvementActions: toStringArray(record.improvementActions),
  publishedAt: toText(record.publishedAt),
  publishedBy: toText(record.publishedBy)
});

const mapAuditEvent = (record: UnknownRecord): BackupAuditEvent => ({
  id: toText(pick(record, ['id', 'eventId'])) || 'unknown-audit-event',
  action: toText(record.action),
  actorUserId: toText(record.actorUserId),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  scopeSnapshot: isRecord(record.scopeSnapshot) ? record.scopeSnapshot : undefined,
  outcome: toText(record.outcome),
  detailSnapshot: isRecord(record.detailSnapshot) ? record.detailSnapshot : undefined,
  occurredAt: toText(record.occurredAt)
});

export const listBackupPolicies = async (query: BackupPolicyListQuery) => {
  const response = await fetchJSON<unknown>(withQuery('/backup-restore/policies', query));
  return normalizeListResponse(response, mapBackupPolicy);
};

export const createBackupPolicy = async (payload: CreateBackupPolicyPayload) => {
  const response = await fetchJSON<unknown>('/backup-restore/policies', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapBackupPolicy(isRecord(response) ? response : {});
};

export const runBackupPolicy = async (policyId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/policies/${policyId}/run`, {
    method: 'POST'
  });
  return mapRestorePoint(isRecord(response) ? response : {});
};

export const listRestorePoints = async (query: RestorePointListQuery) => {
  const response = await fetchJSON<unknown>(withQuery('/backup-restore/restore-points', query));
  return normalizeListResponse(response, mapRestorePoint);
};

export const getRestorePointDetail = async (restorePointId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/restore-points/${restorePointId}`);
  return mapRestorePoint(isRecord(response) ? response : {});
};

export const listRestoreJobs = async (query: RestoreJobListQuery) => {
  const response = await fetchJSON<unknown>(withQuery('/backup-restore/restore-jobs', query));
  return normalizeListResponse(response, mapRestoreJob);
};

export const createRestoreJob = async (payload: CreateRestoreJobPayload) => {
  const response = await fetchJSON<unknown>('/backup-restore/restore-jobs', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapRestoreJob(isRecord(response) ? response : {});
};

export const validateRestoreJob = async (jobId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/restore-jobs/${jobId}/validate`, {
    method: 'POST'
  });
  return mapPrecheckResult(response);
};

export const createMigrationPlan = async (payload: CreateMigrationPlanPayload) => {
  const response = await fetchJSON<unknown>('/backup-restore/migrations', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapMigrationPlan(isRecord(response) ? response : {});
};

export const listDrillPlans = async () => {
  const response = await fetchJSON<unknown>('/backup-restore/drills/plans');
  return normalizeListResponse(response, mapDrillPlan);
};

export const createDrillPlan = async (payload: CreateDRDrillPlanPayload) => {
  const response = await fetchJSON<unknown>('/backup-restore/drills/plans', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapDrillPlan(isRecord(response) ? response : {});
};

export const runDrillPlan = async (planId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/drills/plans/${planId}/run`, {
    method: 'POST'
  });
  return mapDrillRecord(isRecord(response) ? response : {});
};

export const getDrillRecordDetail = async (recordId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/drills/records/${recordId}`);
  return mapDrillRecord(isRecord(response) ? response : {});
};

export const createDrillReport = async (recordId: string) => {
  const response = await fetchJSON<unknown>(`/backup-restore/drills/records/${recordId}/report`, {
    method: 'POST'
  });
  return mapDrillReport(isRecord(response) ? response : {});
};

export const listBackupRestoreAuditEvents = async (query: BackupAuditListQuery) => {
  const response = await fetchJSON<unknown>(withQuery('/audit/backup-restore/events', query));
  return normalizeListResponse(response, mapAuditEvent);
};

export const backupRestoreScopeSummary = (value: ScopeSelection | undefined) => {
  if (!value || Object.keys(value).length === 0) {
    return '未指定范围';
  }
  return Object.entries(value)
    .map(([key, item]) => {
      if (Array.isArray(item)) {
        return `${key}:${item.length}`;
      }
      if (typeof item === 'string') {
        return `${key}:${item}`;
      }
      if (isRecord(item)) {
        return `${key}:${Object.keys(item).length}`;
      }
      return `${key}`;
    })
    .join(' / ');
};

export const restoreJobQueryScope = (query: RestoreJobListQuery) =>
  buildScopeQueryKey([query.jobType, query.status]);

export const policyQueryScope = (query: BackupPolicyListQuery) =>
  buildScopeQueryKey([query.scopeType, query.status]);

export const restorePointQueryScope = (query: RestorePointListQuery) =>
  buildScopeQueryKey([query.policyId, query.result]);
