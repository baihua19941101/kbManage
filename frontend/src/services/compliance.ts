import { ApiError, fetchJSON } from '@/services/api/client';

export type ComplianceStandardType = 'cis' | 'stig' | 'platform-baseline';
export type ComplianceRecordStatus = 'draft' | 'active' | 'disabled' | 'archived';
export type ComplianceScopeType = 'cluster' | 'node' | 'namespace' | 'resource-set';
export type ScheduleMode = 'manual' | 'scheduled';
export type ScanTriggerSource = 'manual' | 'schedule' | 'recheck';
export type ScanExecutionStatus =
  | 'pending'
  | 'running'
  | 'partially_succeeded'
  | 'succeeded'
  | 'failed'
  | 'canceled';
export type CoverageStatus = 'full' | 'partial' | 'unavailable';
export type FindingResult = 'pass' | 'fail' | 'warn' | 'skipped' | 'error';
export type RiskLevel = 'low' | 'medium' | 'high' | 'critical';
export type FindingRemediationStatus =
  | 'open'
  | 'in_progress'
  | 'exception_active'
  | 'ready_for_recheck'
  | 'closed';
export type EvidenceConfidence = 'high' | 'medium' | 'low';
export type RedactionStatus = 'raw' | 'masked';
export type RemediationTaskStatus = 'todo' | 'in_progress' | 'blocked' | 'done' | 'canceled';
export type ComplianceExceptionStatus =
  | 'pending'
  | 'approved'
  | 'rejected'
  | 'active'
  | 'expired'
  | 'revoked';
export type RecheckStatus = 'pending' | 'running' | 'passed' | 'failed' | 'canceled';
export type RecheckTriggerSource = 'manual' | 'remediation_done' | 'exception_expired';
export type OverviewGroupBy = 'cluster' | 'workspace' | 'project' | 'baseline';
export type TrendScopeType = 'cluster' | 'workspace' | 'project';
export type ArchiveExportScope = 'scans' | 'findings' | 'trends' | 'audit' | 'bundle';
export type ArchiveExportStatus = 'pending' | 'running' | 'succeeded' | 'failed' | 'expired';

export type Pagination<T> = {
  items: T[];
};

export type ComplianceBaseline = {
  id: string;
  name: string;
  standardType: ComplianceStandardType;
  version: string;
  status?: ComplianceRecordStatus;
  ruleCount?: number;
  description?: string;
};

export type CreateComplianceBaselineRequest = {
  name: string;
  standardType: ComplianceStandardType;
  version: string;
  description?: string;
};

export type UpdateComplianceBaselineRequest = Partial<
  Pick<ComplianceBaseline, 'name' | 'description' | 'status'>
>;

export type ScanProfile = {
  id: string;
  name: string;
  baselineId: string;
  workspaceId?: string;
  projectId?: string;
  scopeType: ComplianceScopeType;
  clusterRefs?: string[];
  nodeSelectors?: Record<string, string>;
  namespaceRefs?: string[];
  resourceKinds?: string[];
  scheduleMode: ScheduleMode;
  cronExpression?: string;
  status?: 'draft' | 'active' | 'paused' | 'archived';
  lastRunAt?: string;
};

export type CreateScanProfileRequest = {
  name: string;
  baselineId: string;
  workspaceId?: string;
  projectId?: string;
  scopeType: ComplianceScopeType;
  clusterRefs?: string[];
  nodeSelectors?: Record<string, string>;
  namespaceRefs?: string[];
  resourceKinds?: string[];
  scheduleMode: ScheduleMode;
  cronExpression?: string;
};

export type UpdateScanProfileRequest = Partial<
  Pick<ScanProfile, 'name' | 'status' | 'nodeSelectors' | 'cronExpression'>
>;

export type BaselineSnapshot = {
  baselineId?: string;
  name?: string;
  standardType?: ComplianceStandardType;
  version?: string;
};

export type ScanExecution = {
  id: string;
  profileId: string;
  baselineSnapshot?: BaselineSnapshot;
  triggerSource?: ScanTriggerSource;
  status?: ScanExecutionStatus;
  coverageStatus?: CoverageStatus;
  score?: number;
  passCount?: number;
  failCount?: number;
  warningCount?: number;
  startedAt?: string;
  completedAt?: string;
  errorSummary?: string;
};

export type ExecuteScanRequest = {
  reason?: string;
};

export type ComplianceFinding = {
  id: string;
  scanExecutionId?: string;
  controlId?: string;
  controlTitle?: string;
  result?: FindingResult;
  riskLevel?: RiskLevel;
  clusterId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  remediationStatus?: FindingRemediationStatus;
  summary?: string;
};

export type EvidenceRecord = {
  id: string;
  evidenceType?: string;
  sourceRef?: string;
  collectedAt?: string;
  confidence?: EvidenceConfidence;
  summary?: string;
  artifactRef?: string;
  redactionStatus?: RedactionStatus;
};

export type RemediationTask = {
  id: string;
  findingId?: string;
  title?: string;
  owner?: string;
  priority?: RiskLevel;
  status?: RemediationTaskStatus;
  dueAt?: string;
  resolutionSummary?: string;
};

export type ComplianceExceptionRequest = {
  id: string;
  findingId?: string;
  status?: ComplianceExceptionStatus;
  reason?: string;
  startsAt?: string;
  expiresAt?: string;
  reviewComment?: string;
};

export type RecheckTask = {
  id: string;
  findingId?: string;
  triggerSource?: RecheckTriggerSource;
  status?: RecheckStatus;
  resultScanExecutionId?: string;
  summary?: string;
};

export type ComplianceFindingDetail = ComplianceFinding & {
  evidences?: EvidenceRecord[];
  remediationTasks?: RemediationTask[];
  exceptions?: ComplianceExceptionRequest[];
  rechecks?: RecheckTask[];
};

export type ComplianceOverviewGroup = {
  groupKey?: string;
  scoreAvg?: number;
  coverageRate?: number;
  openFindingsCount?: number;
  highRiskOpenCount?: number;
};

export type ComplianceOverview = {
  coverageRate?: number;
  openFindingsCount?: number;
  highRiskOpenCount?: number;
  remediationCompletionRate?: number;
  groups?: ComplianceOverviewGroup[];
};

export type TrendComparisonBasis = {
  baselineId?: string;
  baselineVersions?: string[];
  mixedBaselineFlag?: boolean;
};

export type ComplianceTrendPoint = {
  windowStart?: string;
  windowEnd?: string;
  scoreAvg?: number;
  coverageRate?: number;
  remediationCompletionRate?: number;
  highRiskOpenCount?: number;
  baselineVersion?: string;
};

export type ComplianceTrendResponse = {
  points?: ComplianceTrendPoint[];
  comparisonBasis?: TrendComparisonBasis;
};

export type ArchiveExportTask = {
  id: string;
  exportScope?: ArchiveExportScope;
  status?: ArchiveExportStatus;
  artifactRef?: string;
  requestedBy?: string;
  startedAt?: string;
  completedAt?: string;
  failureReason?: string;
};

export type ComplianceAuditEvent = {
  action?: string;
  operatorId?: string;
  outcome?: string;
  occurredAt?: string;
  details?: Record<string, unknown>;
};

export type BaselineListQuery = {
  standardType?: ComplianceStandardType;
  status?: ComplianceRecordStatus;
};

export type ScanProfileListQuery = {
  workspaceId?: string;
  projectId?: string;
  scopeType?: ComplianceScopeType;
  scheduleMode?: ScheduleMode;
  status?: ScanProfile['status'];
};

export type ScanExecutionListQuery = {
  workspaceId?: string;
  projectId?: string;
  profileId?: string;
  status?: ScanExecutionStatus;
  triggerSource?: ScanTriggerSource;
  timeFrom?: string;
  timeTo?: string;
};

export type ComplianceFindingListQuery = {
  baselineId?: string;
  clusterId?: string;
  namespace?: string;
  remediationStatus?: FindingRemediationStatus;
  riskLevel?: RiskLevel;
  timeFrom?: string;
  timeTo?: string;
  workspaceId?: string;
  projectId?: string;
  scanId?: string;
  result?: FindingResult;
};

export type RemediationTaskListQuery = {
  workspaceId?: string;
  projectId?: string;
  owner?: string;
  status?: RemediationTaskStatus;
  priority?: RiskLevel;
  timeFrom?: string;
  timeTo?: string;
};

export type ComplianceExceptionListQuery = {
  workspaceId?: string;
  projectId?: string;
  status?: ComplianceExceptionStatus;
  baselineId?: string;
};

export type RecheckTaskListQuery = {
  workspaceId?: string;
  projectId?: string;
  status?: RecheckStatus;
  triggerSource?: RecheckTriggerSource;
  timeFrom?: string;
  timeTo?: string;
};

export type ComplianceOverviewQuery = {
  workspaceId?: string;
  projectId?: string;
  groupBy?: OverviewGroupBy;
  timeFrom?: string;
  timeTo?: string;
};

export type ComplianceTrendQuery = {
  workspaceId?: string;
  projectId?: string;
  baselineId?: string;
  scopeType?: TrendScopeType;
  scopeRef?: string;
  timeFrom?: string;
  timeTo?: string;
};

export type ArchiveExportListQuery = {
  workspaceId?: string;
  projectId?: string;
  exportScope?: ArchiveExportScope;
  status?: ArchiveExportStatus;
  timeFrom?: string;
  timeTo?: string;
};

export type ComplianceAuditQuery = {
  workspaceId?: string;
  projectId?: string;
  baselineId?: string;
  action?: string;
  outcome?: string;
  timeFrom?: string;
  timeTo?: string;
};

export type CreateRemediationTaskRequest = {
  title: string;
  owner: string;
  priority: RiskLevel;
  dueAt?: string;
  summary?: string;
};

export type UpdateRemediationTaskRequest = {
  status?: RemediationTaskStatus;
  resolutionSummary?: string;
};

export type CreateComplianceExceptionRequest = {
  reason: string;
  startsAt: string;
  expiresAt: string;
};

export type ReviewComplianceExceptionRequest = {
  decision: 'approve' | 'reject' | 'revoke';
  reviewComment?: string;
};

export type CreateRecheckTaskRequest = {
  reason?: string;
};

export type CreateArchiveExportTaskRequest = {
  workspaceId?: string;
  projectId?: string;
  baselineId?: string;
  exportScope: ArchiveExportScope;
  timeFrom?: string;
  timeTo?: string;
  filters?: Record<string, unknown>;
};

const toQueryString = (query: Record<string, unknown>) => {
  const params = new URLSearchParams();

  Object.entries(query).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') {
      return;
    }
    params.set(key, String(value));
  });

  const result = params.toString();
  return result ? `?${result}` : '';
};

const withQuery = (path: string, query: Record<string, unknown>) => `${path}${toQueryString(query)}`;

const asRecord = (value: unknown): Record<string, unknown> =>
  typeof value === 'object' && value !== null ? (value as Record<string, unknown>) : {};

const asArray = (value: unknown): unknown[] => (Array.isArray(value) ? value : []);

const readString = (record: Record<string, unknown>, key: string): string | undefined => {
  const value = record[key];
  return typeof value === 'string' && value.trim() ? value.trim() : undefined;
};

const readNumber = (record: Record<string, unknown>, key: string): number | undefined => {
  const value = record[key];
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === 'string' && value.trim()) {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : undefined;
  }
  return undefined;
};

const readObject = (record: Record<string, unknown>, key: string): Record<string, unknown> | undefined => {
  const value = record[key];
  return typeof value === 'object' && value !== null && !Array.isArray(value)
    ? (value as Record<string, unknown>)
    : undefined;
};

const readStringArray = (record: Record<string, unknown>, key: string): string[] | undefined => {
  const value = record[key];
  if (!Array.isArray(value)) {
    return undefined;
  }
  return value.filter((item): item is string => typeof item === 'string' && item.trim().length > 0);
};

const mapBaseline = (input: unknown): ComplianceBaseline => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    name: readString(record, 'name') || '',
    standardType: (readString(record, 'standardType') as ComplianceStandardType) || 'cis',
    version: readString(record, 'version') || '',
    status: readString(record, 'status') as ComplianceRecordStatus | undefined,
    ruleCount: readNumber(record, 'ruleCount'),
    description: readString(record, 'description')
  };
};

const mapNodeSelectors = (input: unknown): Record<string, string> | undefined => {
  const record = typeof input === 'object' && input !== null && !Array.isArray(input)
    ? (input as Record<string, unknown>)
    : undefined;
  if (!record) {
    return undefined;
  }
  const entries = Object.entries(record).filter(
    (entry): entry is [string, string] => typeof entry[1] === 'string'
  );
  return entries.length > 0 ? Object.fromEntries(entries) : undefined;
};

const mapScanProfile = (input: unknown): ScanProfile => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    name: readString(record, 'name') || '',
    baselineId: readString(record, 'baselineId') || '',
    workspaceId: readString(record, 'workspaceId'),
    projectId: readString(record, 'projectId'),
    scopeType: (readString(record, 'scopeType') as ComplianceScopeType) || 'cluster',
    clusterRefs: readStringArray(record, 'clusterRefs'),
    nodeSelectors: mapNodeSelectors(record.nodeSelectors),
    namespaceRefs: readStringArray(record, 'namespaceRefs'),
    resourceKinds: readStringArray(record, 'resourceKinds'),
    scheduleMode: (readString(record, 'scheduleMode') as ScheduleMode) || 'manual',
    cronExpression: readString(record, 'cronExpression'),
    status: readString(record, 'status') as ScanProfile['status'] | undefined,
    lastRunAt: readString(record, 'lastRunAt')
  };
};

const mapBaselineSnapshot = (input: unknown): BaselineSnapshot | undefined => {
  const record = asRecord(input);
  if (Object.keys(record).length === 0) {
    return undefined;
  }
  return {
    baselineId: readString(record, 'baselineId'),
    name: readString(record, 'name'),
    standardType: readString(record, 'standardType') as ComplianceStandardType | undefined,
    version: readString(record, 'version')
  };
};

const mapScanExecution = (input: unknown): ScanExecution => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    profileId: readString(record, 'profileId') || '',
    baselineSnapshot: mapBaselineSnapshot(record.baselineSnapshot),
    triggerSource: readString(record, 'triggerSource') as ScanTriggerSource | undefined,
    status: readString(record, 'status') as ScanExecutionStatus | undefined,
    coverageStatus: readString(record, 'coverageStatus') as CoverageStatus | undefined,
    score: readNumber(record, 'score'),
    passCount: readNumber(record, 'passCount'),
    failCount: readNumber(record, 'failCount'),
    warningCount: readNumber(record, 'warningCount'),
    startedAt: readString(record, 'startedAt'),
    completedAt: readString(record, 'completedAt'),
    errorSummary: readString(record, 'errorSummary')
  };
};

const mapFinding = (input: unknown): ComplianceFinding => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    scanExecutionId: readString(record, 'scanExecutionId'),
    controlId: readString(record, 'controlId'),
    controlTitle: readString(record, 'controlTitle'),
    result: readString(record, 'result') as FindingResult | undefined,
    riskLevel: readString(record, 'riskLevel') as RiskLevel | undefined,
    clusterId: readString(record, 'clusterId'),
    namespace: readString(record, 'namespace'),
    resourceKind: readString(record, 'resourceKind'),
    resourceName: readString(record, 'resourceName'),
    remediationStatus: readString(record, 'remediationStatus') as FindingRemediationStatus | undefined,
    summary: readString(record, 'summary')
  };
};

const mapEvidence = (input: unknown): EvidenceRecord => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    evidenceType: readString(record, 'evidenceType'),
    sourceRef: readString(record, 'sourceRef'),
    collectedAt: readString(record, 'collectedAt'),
    confidence: readString(record, 'confidence') as EvidenceConfidence | undefined,
    summary: readString(record, 'summary'),
    artifactRef: readString(record, 'artifactRef'),
    redactionStatus: readString(record, 'redactionStatus') as RedactionStatus | undefined
  };
};

const mapRemediationTask = (input: unknown): RemediationTask => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    findingId: readString(record, 'findingId'),
    title: readString(record, 'title'),
    owner: readString(record, 'owner'),
    priority: readString(record, 'priority') as RiskLevel | undefined,
    status: readString(record, 'status') as RemediationTaskStatus | undefined,
    dueAt: readString(record, 'dueAt'),
    resolutionSummary: readString(record, 'resolutionSummary')
  };
};

const mapException = (input: unknown): ComplianceExceptionRequest => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    findingId: readString(record, 'findingId'),
    status: readString(record, 'status') as ComplianceExceptionStatus | undefined,
    reason: readString(record, 'reason'),
    startsAt: readString(record, 'startsAt'),
    expiresAt: readString(record, 'expiresAt'),
    reviewComment: readString(record, 'reviewComment')
  };
};

const mapRecheck = (input: unknown): RecheckTask => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    findingId: readString(record, 'findingId'),
    triggerSource: readString(record, 'triggerSource') as RecheckTriggerSource | undefined,
    status: readString(record, 'status') as RecheckStatus | undefined,
    resultScanExecutionId: readString(record, 'resultScanExecutionId'),
    summary: readString(record, 'summary')
  };
};

const mapFindingDetail = (input: unknown): ComplianceFindingDetail => {
  const record = asRecord(input);
  return {
    ...mapFinding(record),
    evidences: asArray(record.evidences).map(mapEvidence),
    remediationTasks: asArray(record.remediationTasks).map(mapRemediationTask),
    exceptions: asArray(record.exceptions).map(mapException),
    rechecks: asArray(record.rechecks).map(mapRecheck)
  };
};

const mapOverviewGroup = (input: unknown): ComplianceOverviewGroup => {
  const record = asRecord(input);
  return {
    groupKey: readString(record, 'groupKey'),
    scoreAvg: readNumber(record, 'scoreAvg'),
    coverageRate: readNumber(record, 'coverageRate'),
    openFindingsCount: readNumber(record, 'openFindingsCount'),
    highRiskOpenCount: readNumber(record, 'highRiskOpenCount')
  };
};

const mapOverview = (input: unknown): ComplianceOverview => {
  const record = asRecord(input);
  return {
    coverageRate: readNumber(record, 'coverageRate'),
    openFindingsCount: readNumber(record, 'openFindingsCount'),
    highRiskOpenCount: readNumber(record, 'highRiskOpenCount'),
    remediationCompletionRate: readNumber(record, 'remediationCompletionRate'),
    groups: asArray(record.groups).map(mapOverviewGroup)
  };
};

const mapTrendPoint = (input: unknown): ComplianceTrendPoint => {
  const record = asRecord(input);
  return {
    windowStart: readString(record, 'windowStart'),
    windowEnd: readString(record, 'windowEnd'),
    scoreAvg: readNumber(record, 'scoreAvg'),
    coverageRate: readNumber(record, 'coverageRate'),
    remediationCompletionRate: readNumber(record, 'remediationCompletionRate'),
    highRiskOpenCount: readNumber(record, 'highRiskOpenCount'),
    baselineVersion: readString(record, 'baselineVersion')
  };
};

const mapTrendComparison = (input: unknown): TrendComparisonBasis | undefined => {
  const record = asRecord(input);
  if (Object.keys(record).length === 0) {
    return undefined;
  }
  return {
    baselineId: readString(record, 'baselineId'),
    baselineVersions: readStringArray(record, 'baselineVersions'),
    mixedBaselineFlag: Boolean(record.mixedBaselineFlag)
  };
};

const mapArchiveExport = (input: unknown): ArchiveExportTask => {
  const record = asRecord(input);
  return {
    id: readString(record, 'id') || '',
    exportScope: readString(record, 'exportScope') as ArchiveExportScope | undefined,
    status: readString(record, 'status') as ArchiveExportStatus | undefined,
    artifactRef: readString(record, 'artifactRef'),
    requestedBy: readString(record, 'requestedBy'),
    startedAt: readString(record, 'startedAt'),
    completedAt: readString(record, 'completedAt'),
    failureReason: readString(record, 'failureReason')
  };
};

const mapAuditEvent = (input: unknown): ComplianceAuditEvent => {
  const record = asRecord(input);
  return {
    action: readString(record, 'action'),
    operatorId: readString(record, 'operatorId'),
    outcome: readString(record, 'outcome'),
    occurredAt: readString(record, 'occurredAt'),
    details: readObject(record, 'details')
  };
};

const mapItems = <T>(payload: unknown, mapper: (value: unknown) => T): Pagination<T> => {
  const record = asRecord(payload);
  const items = Array.isArray(record.items) ? record.items : [];
  return { items: items.map(mapper) };
};

const ensureId = (id: string, label: string) => {
  const normalized = id.trim();
  if (!normalized) {
    throw new ApiError(400, `${label}不能为空`, { url: '/compliance' });
  }
  return normalized;
};

export const listComplianceBaselines = async (query: BaselineListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/baselines', query));
  return mapItems(response, mapBaseline);
};

export const createComplianceBaseline = async (payload: CreateComplianceBaselineRequest) =>
  mapBaseline(
    await fetchJSON<unknown>('/compliance/baselines', {
      method: 'POST',
      body: JSON.stringify(payload)
    })
  );

export const updateComplianceBaseline = async (
  baselineId: string,
  payload: UpdateComplianceBaselineRequest
) =>
  mapBaseline(
    await fetchJSON<unknown>(`/compliance/baselines/${encodeURIComponent(ensureId(baselineId, 'baselineId'))}`, {
      method: 'PATCH',
      body: JSON.stringify(payload)
    })
  );

export const listScanProfiles = async (query: ScanProfileListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/scan-profiles', query));
  return mapItems(response, mapScanProfile);
};

export const getScanProfile = async (profileId: string) =>
  mapScanProfile(
    await fetchJSON<unknown>(`/compliance/scan-profiles/${encodeURIComponent(ensureId(profileId, 'profileId'))}`)
  );

export const createScanProfile = async (payload: CreateScanProfileRequest) =>
  mapScanProfile(
    await fetchJSON<unknown>('/compliance/scan-profiles', {
      method: 'POST',
      body: JSON.stringify(payload)
    })
  );

export const updateScanProfile = async (
  profileId: string,
  payload: UpdateScanProfileRequest
) =>
  mapScanProfile(
    await fetchJSON<unknown>(`/compliance/scan-profiles/${encodeURIComponent(ensureId(profileId, 'profileId'))}`, {
      method: 'PATCH',
      body: JSON.stringify(payload)
    })
  );

export const executeScanProfile = async (
  profileId: string,
  payload: ExecuteScanRequest = {}
) =>
  mapScanExecution(
    await fetchJSON<unknown>(
      `/compliance/scan-profiles/${encodeURIComponent(ensureId(profileId, 'profileId'))}/execute`,
      {
        method: 'POST',
        body: JSON.stringify(payload)
      }
    )
  );

export const listScanExecutions = async (query: ScanExecutionListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/scans', query));
  return mapItems(response, mapScanExecution);
};

export const getScanExecution = async (scanId: string) =>
  mapScanExecution(
    await fetchJSON<unknown>(`/compliance/scans/${encodeURIComponent(ensureId(scanId, 'scanId'))}`)
  );

export const listComplianceFindings = async (query: ComplianceFindingListQuery = {}) => {
  const path = query.scanId
    ? `/compliance/scans/${encodeURIComponent(query.scanId)}/findings`
    : '/compliance/findings';
  const response = await fetchJSON<unknown>(withQuery(path, {
    ...query,
    scanId: undefined
  }));
  return mapItems(response, mapFinding);
};

export const getComplianceFinding = async (findingId: string) =>
  mapFindingDetail(
    await fetchJSON<unknown>(`/compliance/findings/${encodeURIComponent(ensureId(findingId, 'findingId'))}`)
  );

export const createRemediationTask = async (
  findingId: string,
  payload: CreateRemediationTaskRequest
) =>
  mapRemediationTask(
    await fetchJSON<unknown>(
      `/compliance/findings/${encodeURIComponent(ensureId(findingId, 'findingId'))}/remediation-tasks`,
      {
        method: 'POST',
        body: JSON.stringify(payload)
      }
    )
  );

export const updateRemediationTask = async (
  taskId: string,
  payload: UpdateRemediationTaskRequest
) =>
  mapRemediationTask(
    await fetchJSON<unknown>(
      `/compliance/remediation-tasks/${encodeURIComponent(ensureId(taskId, 'taskId'))}`,
      {
        method: 'PATCH',
        body: JSON.stringify(payload)
      }
    )
  );

export const listRemediationTasks = async (query: RemediationTaskListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/remediation-tasks', query));
  return mapItems(response, mapRemediationTask);
};

export const createComplianceException = async (
  findingId: string,
  payload: CreateComplianceExceptionRequest
) =>
  mapException(
    await fetchJSON<unknown>(
      `/compliance/findings/${encodeURIComponent(ensureId(findingId, 'findingId'))}/exceptions`,
      {
        method: 'POST',
        body: JSON.stringify(payload)
      }
    )
  );

export const listComplianceExceptions = async (query: ComplianceExceptionListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/exceptions', query));
  return mapItems(response, mapException);
};

export const reviewComplianceException = async (
  exceptionId: string,
  payload: ReviewComplianceExceptionRequest
) =>
  mapException(
    await fetchJSON<unknown>(
      `/compliance/exceptions/${encodeURIComponent(ensureId(exceptionId, 'exceptionId'))}/review`,
      {
        method: 'POST',
        body: JSON.stringify(payload)
      }
    )
  );

export const createRecheckTask = async (
  findingId: string,
  payload: CreateRecheckTaskRequest = {}
) =>
  mapRecheck(
    await fetchJSON<unknown>(
      `/compliance/findings/${encodeURIComponent(ensureId(findingId, 'findingId'))}/rechecks`,
      {
        method: 'POST',
        body: JSON.stringify(payload)
      }
    )
  );

export const getRecheckTask = async (recheckId: string) =>
  mapRecheck(
    await fetchJSON<unknown>(`/compliance/rechecks/${encodeURIComponent(ensureId(recheckId, 'recheckId'))}`)
  );

export const listRecheckTasks = async (query: RecheckTaskListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/rechecks', query));
  return mapItems(response, mapRecheck);
};

export const getComplianceOverview = async (query: ComplianceOverviewQuery = {}) =>
  mapOverview(await fetchJSON<unknown>(withQuery('/compliance/overview', query)));

export const getComplianceTrends = async (query: ComplianceTrendQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/trends', query));
  const record = asRecord(response);
  return {
    points: asArray(record.points).map(mapTrendPoint),
    comparisonBasis: mapTrendComparison(record.comparisonBasis)
  } as ComplianceTrendResponse;
};

export const listComplianceArchiveExports = async (query: ArchiveExportListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/compliance/archive-exports', query));
  return mapItems(response, mapArchiveExport);
};

export const createComplianceArchiveExport = async (payload: CreateArchiveExportTaskRequest) =>
  mapArchiveExport(
    await fetchJSON<unknown>('/compliance/archive-exports', {
      method: 'POST',
      body: JSON.stringify(payload)
    })
  );

export const getComplianceArchiveExport = async (exportId: string) =>
  mapArchiveExport(
    await fetchJSON<unknown>(
      `/compliance/archive-exports/${encodeURIComponent(ensureId(exportId, 'exportId'))}`
    )
  );

export const listComplianceAuditEvents = async (query: ComplianceAuditQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/audit/compliance/events', query));
  return mapItems(response, mapAuditEvent);
};
