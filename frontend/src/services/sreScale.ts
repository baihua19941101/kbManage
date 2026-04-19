import { fetchJSON } from '@/services/api/client';

type UnknownRecord = Record<string, unknown>;
type QueryValue = string | number | boolean | undefined | null;

const isRecord = (value: unknown): value is UnknownRecord =>
  typeof value === 'object' && value !== null;

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
  return value.map((item) => toText(item)).filter((item): item is string => Boolean(item));
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

const normalizeListResponse = <T>(value: unknown, mapper: (record: UnknownRecord) => T) => {
  if (Array.isArray(value)) {
    return { items: value.filter(isRecord).map(mapper) };
  }
  const record = isRecord(value) ? value : {};
  const items = Array.isArray(record.items) ? record.items : [];
  return { items: items.filter(isRecord).map(mapper) };
};

export type HAPolicy = {
  id: string;
  name: string;
  deploymentMode?: string;
  controlPlaneScope?: string;
  replicaExpectation?: number;
  takeoverStatus?: string;
  status?: string;
  lastRecoveryResult?: string;
};

export type MaintenanceWindow = {
  id: string;
  name: string;
  windowType?: string;
  scope?: string;
  startAt?: string;
  endAt?: string;
  status?: string;
  postCheckStatus?: string;
};

export type PlatformHealthOverview = {
  id: string;
  overallStatus?: string;
  componentHealthSummary?: string;
  dependencyHealthSummary?: string;
  taskBacklogSummary?: string;
  capacityRiskLevel?: string;
  throttlingStatus?: string;
  recoverySummary?: string;
  maintenanceStatus?: string;
  recommendedActions: string[];
};

export type UpgradePrecheckResult = {
  decision?: string;
  compatibilitySummary?: string;
  blockers: string[];
  warnings: string[];
  allowConditions: string[];
};

export type SREUpgradePlan = {
  id: string;
  name: string;
  currentVersion?: string;
  targetVersion?: string;
  executionStage?: string;
  executionProgress?: number;
  acceptanceResult?: string;
  rollbackReadiness?: string;
  status?: string;
};

export type RollbackValidation = {
  id: string;
  validationScope?: string;
  result?: string;
  remainingRisk?: string;
  validatedAt?: string;
};

export type CapacityBaseline = {
  id: string;
  name: string;
  resourceDimension?: string;
  baselineRange?: string;
  growthTrend?: string;
  forecastResult?: string;
  confidenceLevel?: string;
  status?: string;
};

export type ScaleEvidence = {
  id: string;
  evidenceType?: string;
  summary?: string;
  bottleneckSummary?: string;
  forecastSummary?: string;
  confidenceLevel?: string;
  capturedAt?: string;
};

export type RunbookArticle = {
  id: string;
  title: string;
  scenarioType?: string;
  riskLevel?: string;
  checklistSummary?: string;
  recoverySteps?: string;
  verificationSummary?: string;
  status?: string;
};

export type SREAuditEvent = {
  id: string;
  action?: string;
  targetType?: string;
  targetRef?: string;
  outcome?: string;
  occurredAt?: string;
};

const mapHAPolicy = (record: UnknownRecord): HAPolicy => ({
  id: toText(record.id) || 'unknown-ha',
  name: toText(record.name) || '未命名策略',
  deploymentMode: toText(record.deploymentMode),
  controlPlaneScope: toText(record.controlPlaneScope),
  replicaExpectation: toNumber(record.replicaExpectation),
  takeoverStatus: toText(record.takeoverStatus),
  status: toText(record.status),
  lastRecoveryResult: toText(record.lastRecoveryResult)
});

const mapMaintenanceWindow = (record: UnknownRecord): MaintenanceWindow => ({
  id: toText(record.id) || 'unknown-window',
  name: toText(record.name) || '未命名窗口',
  windowType: toText(record.windowType),
  scope: toText(record.scope),
  startAt: toText(record.startAt),
  endAt: toText(record.endAt),
  status: toText(record.status),
  postCheckStatus: toText(record.postCheckStatus)
});

const mapHealthOverview = (record: UnknownRecord): PlatformHealthOverview => ({
  id: toText(record.id) || 'health-overview',
  overallStatus: toText(record.overallStatus),
  componentHealthSummary: toText(record.componentHealthSummary),
  dependencyHealthSummary: toText(record.dependencyHealthSummary),
  taskBacklogSummary: toText(record.taskBacklogSummary),
  capacityRiskLevel: toText(record.capacityRiskLevel),
  throttlingStatus: toText(record.throttlingStatus),
  recoverySummary: toText(record.recoverySummary),
  maintenanceStatus: toText(record.maintenanceStatus),
  recommendedActions:
    toStringArray(record.recommendedActions).length > 0
      ? toStringArray(record.recommendedActions)
      : toText(record.recommendedActions)?.split('；').filter(Boolean) ?? []
});

const mapUpgradePrecheck = (record: UnknownRecord): UpgradePrecheckResult => ({
  decision: toText(record.decision),
  compatibilitySummary: toText(record.compatibilitySummary),
  blockers: toStringArray(record.blockers),
  warnings: toStringArray(record.warnings),
  allowConditions: toStringArray(record.allowConditions)
});

const mapUpgradePlan = (record: UnknownRecord): SREUpgradePlan => ({
  id: toText(record.id) || 'unknown-upgrade',
  name: toText(record.name) || '未命名升级',
  currentVersion: toText(record.currentVersion),
  targetVersion: toText(record.targetVersion),
  executionStage: toText(record.executionStage),
  executionProgress: toNumber(record.executionProgress),
  acceptanceResult: toText(record.acceptanceResult),
  rollbackReadiness: toText(record.rollbackReadiness),
  status: toText(record.status)
});

const mapRollbackValidation = (record: UnknownRecord): RollbackValidation => ({
  id: toText(record.id) || 'unknown-rollback',
  validationScope: toText(record.validationScope),
  result: toText(record.result),
  remainingRisk: toText(record.remainingRisk),
  validatedAt: toText(record.validatedAt)
});

const mapCapacityBaseline = (record: UnknownRecord): CapacityBaseline => ({
  id: toText(record.id) || 'unknown-capacity',
  name: toText(record.name) || '未命名基线',
  resourceDimension: toText(record.resourceDimension),
  baselineRange: toText(record.baselineRange),
  growthTrend: toText(record.growthTrend),
  forecastResult: toText(record.forecastResult),
  confidenceLevel: toText(record.confidenceLevel),
  status: toText(record.status)
});

const mapScaleEvidence = (record: UnknownRecord): ScaleEvidence => ({
  id: toText(record.id) || 'unknown-evidence',
  evidenceType: toText(record.evidenceType),
  summary: toText(record.summary),
  bottleneckSummary: toText(record.bottleneckSummary),
  forecastSummary: toText(record.forecastSummary),
  confidenceLevel: toText(record.confidenceLevel),
  capturedAt: toText(record.capturedAt)
});

const mapRunbook = (record: UnknownRecord): RunbookArticle => ({
  id: toText(record.id) || 'unknown-runbook',
  title: toText(record.title) || '未命名手册',
  scenarioType: toText(record.scenarioType),
  riskLevel: toText(record.riskLevel),
  checklistSummary: toText(record.checklistSummary),
  recoverySteps: toText(record.recoverySteps),
  verificationSummary: toText(record.verificationSummary),
  status: toText(record.status)
});

const mapAuditEvent = (record: UnknownRecord): SREAuditEvent => ({
  id: toText(record.id) || 'unknown-audit',
  action: toText(record.action),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  outcome: toText(record.outcome),
  occurredAt: toText(record.occurredAt)
});

export const listHAPolicies = async (query: Record<string, QueryValue> = {}) =>
  normalizeListResponse(await fetchJSON<unknown>(withQuery('/sre/ha-policies', query)), mapHAPolicy);

export const createHAPolicy = async (payload: Record<string, unknown>) =>
  mapHAPolicy(
    await fetchJSON<unknown>('/sre/ha-policies', { method: 'POST', body: JSON.stringify(payload) }).then((v) =>
      isRecord(v) ? v : {}
    )
  );

export const listMaintenanceWindows = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/sre/maintenance-windows'), mapMaintenanceWindow);

export const createMaintenanceWindow = async (payload: Record<string, unknown>) =>
  mapMaintenanceWindow(
    await fetchJSON<unknown>('/sre/maintenance-windows', { method: 'POST', body: JSON.stringify(payload) }).then((v) =>
      isRecord(v) ? v : {}
    )
  );

export const getHealthOverview = async (query: Record<string, QueryValue> = {}) =>
  mapHealthOverview(await fetchJSON<unknown>(withQuery('/sre/health/overview', query)).then((v) => (isRecord(v) ? v : {})));

export const runUpgradePrecheck = async (payload: Record<string, unknown>) =>
  mapUpgradePrecheck(
    await fetchJSON<unknown>('/sre/upgrades/prechecks', { method: 'POST', body: JSON.stringify(payload) }).then((v) =>
      isRecord(v) ? v : {}
    )
  );

export const listUpgradePlans = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/sre/upgrades'), mapUpgradePlan);

export const createUpgradePlan = async (payload: Record<string, unknown>) =>
  mapUpgradePlan(
    await fetchJSON<unknown>('/sre/upgrades', { method: 'POST', body: JSON.stringify(payload) }).then((v) =>
      isRecord(v) ? v : {}
    )
  );

export const createRollbackValidation = async (upgradeId: string, payload: Record<string, unknown>) =>
  mapRollbackValidation(
    await fetchJSON<unknown>(`/sre/upgrades/${upgradeId}/rollback-validations`, {
      method: 'POST',
      body: JSON.stringify(payload)
    }).then((v) => (isRecord(v) ? v : {}))
  );

export const listCapacityBaselines = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/sre/capacity/baselines'), mapCapacityBaseline);

export const listScaleEvidence = async (query: Record<string, QueryValue> = {}) =>
  normalizeListResponse(await fetchJSON<unknown>(withQuery('/sre/scale-evidence', query)), mapScaleEvidence);

export const listRunbooks = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/sre/runbooks'), mapRunbook);

export const listSREAuditEvents = async (query: Record<string, QueryValue> = {}) =>
  normalizeListResponse(await fetchJSON<unknown>(withQuery('/audit/sre/events', query)), mapAuditEvent);
