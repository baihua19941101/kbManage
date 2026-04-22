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
    if (value === undefined || value === null || value === '') return;
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

export type PermissionTrail = {
  id: string;
  subjectType?: string;
  subjectRef?: string;
  changeType?: string;
  beforeState?: string;
  afterState?: string;
  authorizationBasis?: string;
  evidenceCompleteness?: string;
  changedAt?: string;
};

export type KeyOperationTrace = {
  id: string;
  actorRef?: string;
  operationType?: string;
  targetType?: string;
  targetRef?: string;
  contextSummary?: string;
  riskLevel?: string;
  outcome?: string;
  occurredAt?: string;
};

export type GovernanceCoverageSnapshot = {
  id: string;
  coverageDomain?: string;
  coverageRate?: number;
  statusBreakdown?: string[];
  missingReasonSummary?: string;
  confidenceLevel?: string;
  trendSummary?: string;
  snapshotAt?: string;
};

export type GovernanceActionItem = {
  id: string;
  sourceType?: string;
  sourceRef?: string;
  title: string;
  priority?: string;
  owner?: string;
  status?: string;
  resolutionSummary?: string;
};

export type GovernanceReportPackage = {
  id: string;
  reportType?: string;
  title: string;
  audienceType?: string;
  timeRange?: string;
  summarySection?: string;
  detailSection?: string;
  attachmentCatalog: string[];
  status?: string;
};

export type ExportRecord = {
  id: string;
  exportType?: string;
  audienceScope?: string;
  contentLevel?: string;
  result?: string;
  exportedAt?: string;
};

export type DeliveryArtifact = {
  id: string;
  artifactType?: string;
  title: string;
  versionScope?: string;
  environmentScope?: string;
  ownerRole?: string;
  status?: string;
  applicabilityNote?: string;
};

export type DeliveryReadinessBundle = {
  id: string;
  name: string;
  targetEnvironment?: string;
  targetAudience?: string;
  artifactSummary?: string;
  checklistStatus?: string;
  missingItems: string[];
  readinessConclusion?: string;
};

export type DeliveryChecklistItem = {
  id: string;
  checkItem: string;
  category?: string;
  owner?: string;
  evidenceRequirement?: string;
  status?: string;
  remark?: string;
  completedAt?: string;
};

export type EnterpriseAuditEvent = {
  id: string;
  action?: string;
  targetType?: string;
  targetRef?: string;
  outcome?: string;
  occurredAt?: string;
};

export type CreateGovernanceReportPayload = {
  workspaceId: number;
  reportType: string;
  title: string;
  audienceType: string;
  timeRange?: string;
  visibilityPolicy?: string;
};

export type CreateExportRecordPayload = {
  audienceScope: string;
  contentLevel: string;
  exportType: string;
};

const mapPermissionTrail = (record: UnknownRecord): PermissionTrail => ({
  id: toText(record.id) || 'unknown-trail',
  subjectType: toText(record.subjectType),
  subjectRef: toText(record.subjectRef),
  changeType: toText(record.changeType),
  beforeState: toText(record.beforeState),
  afterState: toText(record.afterState),
  authorizationBasis: toText(record.authorizationBasis),
  evidenceCompleteness: toText(record.evidenceCompleteness),
  changedAt: toText(record.changedAt)
});

const mapKeyOperation = (record: UnknownRecord): KeyOperationTrace => ({
  id: toText(record.id) || 'unknown-operation',
  actorRef: toText(record.actorRef),
  operationType: toText(record.operationType),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  contextSummary: toText(record.contextSummary),
  riskLevel: toText(record.riskLevel),
  outcome: toText(record.outcome),
  occurredAt: toText(record.occurredAt)
});

const mapCoverage = (record: UnknownRecord): GovernanceCoverageSnapshot => ({
  id: toText(record.id) || 'unknown-coverage',
  coverageDomain: toText(record.coverageDomain),
  coverageRate: toNumber(record.coverageRate),
  statusBreakdown: toStringArray(record.statusBreakdown),
  missingReasonSummary: toText(record.missingReasonSummary),
  confidenceLevel: toText(record.confidenceLevel),
  trendSummary: toText(record.trendSummary),
  snapshotAt: toText(record.snapshotAt)
});

const mapActionItem = (record: UnknownRecord): GovernanceActionItem => ({
  id: toText(record.id) || 'unknown-action',
  sourceType: toText(record.sourceType),
  sourceRef: toText(record.sourceRef),
  title: toText(record.title) || '未命名待办',
  priority: toText(record.priority),
  owner: toText(record.owner),
  status: toText(record.status),
  resolutionSummary: toText(record.resolutionSummary)
});

const mapReport = (record: UnknownRecord): GovernanceReportPackage => ({
  id: toText(record.id) || 'unknown-report',
  reportType: toText(record.reportType),
  title: toText(record.title) || '未命名报表',
  audienceType: toText(record.audienceType),
  timeRange: toText(record.timeRange),
  summarySection: toText(record.summarySection),
  detailSection: toText(record.detailSection),
  attachmentCatalog: toStringArray(record.attachmentCatalog),
  status: toText(record.status)
});

const mapExportRecord = (record: UnknownRecord): ExportRecord => ({
  id: toText(record.id) || 'unknown-export',
  exportType: toText(record.exportType),
  audienceScope: toText(record.audienceScope),
  contentLevel: toText(record.contentLevel),
  result: toText(record.result),
  exportedAt: toText(record.exportedAt)
});

const mapDeliveryArtifact = (record: UnknownRecord): DeliveryArtifact => ({
  id: toText(record.id) || 'unknown-artifact',
  artifactType: toText(record.artifactType),
  title: toText(record.title) || '未命名材料',
  versionScope: toText(record.versionScope),
  environmentScope: toText(record.environmentScope),
  ownerRole: toText(record.ownerRole),
  status: toText(record.status),
  applicabilityNote: toText(record.applicabilityNote)
});

const mapBundle = (record: UnknownRecord): DeliveryReadinessBundle => ({
  id: toText(record.id) || 'unknown-bundle',
  name: toText(record.name) || '未命名交付包',
  targetEnvironment: toText(record.targetEnvironment),
  targetAudience: toText(record.targetAudience),
  artifactSummary: toText(record.artifactSummary),
  checklistStatus: toText(record.checklistStatus),
  missingItems: toStringArray(record.missingItems),
  readinessConclusion: toText(record.readinessConclusion)
});

const mapChecklist = (record: UnknownRecord): DeliveryChecklistItem => ({
  id: toText(record.id) || 'unknown-check',
  checkItem: toText(record.checkItem) || '未命名检查项',
  category: toText(record.category),
  owner: toText(record.owner),
  evidenceRequirement: toText(record.evidenceRequirement),
  status: toText(record.status),
  remark: toText(record.remark),
  completedAt: toText(record.completedAt)
});

const mapAuditEvent = (record: UnknownRecord): EnterpriseAuditEvent => ({
  id: toText(record.id) || 'unknown-audit',
  action: toText(record.action),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  outcome: toText(record.outcome),
  occurredAt: toText(record.occurredAt)
});

export const listPermissionTrails = async () =>
  normalizeListResponse(
    await fetchJSON<unknown>('/enterprise/audit/permission-trails'),
    mapPermissionTrail
  );

export const listKeyOperations = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/audit/key-operations'), mapKeyOperation);

export const listGovernanceCoverage = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/governance/coverage'), mapCoverage);

export const listGovernanceActionItems = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/governance/action-items'), mapActionItem);

export const listGovernanceReports = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/reports'), mapReport);

export const createGovernanceReport = async (payload: CreateGovernanceReportPayload) =>
  mapReport(
    (await fetchJSON<unknown>('/enterprise/reports', {
      method: 'POST',
      body: JSON.stringify(payload)
    })) as UnknownRecord
  );

export const createExportRecord = async (reportId: string, payload: CreateExportRecordPayload) =>
  mapExportRecord(
    (await fetchJSON<unknown>(`/enterprise/reports/${reportId}/exports`, {
      method: 'POST',
      body: JSON.stringify(payload)
    })) as UnknownRecord
  );

export const listDeliveryArtifacts = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/delivery/artifacts'), mapDeliveryArtifact);

export const listDeliveryBundles = async () =>
  normalizeListResponse(await fetchJSON<unknown>('/enterprise/delivery/bundles'), mapBundle);

export const listDeliveryChecklist = async (bundleId: string) =>
  normalizeListResponse(
    await fetchJSON<unknown>(withQuery(`/enterprise/delivery/bundles/${bundleId}/checklists`, {})),
    mapChecklist
  );

export const listEnterpriseAuditEvents = async (query: Record<string, QueryValue>) =>
  normalizeListResponse(
    await fetchJSON<unknown>(withQuery('/audit/enterprise/events', query)),
    mapAuditEvent
  );
