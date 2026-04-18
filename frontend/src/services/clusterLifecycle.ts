import { fetchJSON } from '@/services/api/client';

type UnknownRecord = Record<string, unknown>;
type QueryValue = string | number | boolean | undefined | null;

const UNKNOWN_TEXT = '未知';
const UNGROUPED_DOMAIN = '未归类';

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

const toBoolean = (value: unknown): boolean | undefined => {
  if (typeof value === 'boolean') {
    return value;
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

const toStringArray = (value: unknown): string[] | undefined => {
  if (!Array.isArray(value)) {
    return undefined;
  }
  return value
    .map((item) => toText(item))
    .filter((item): item is string => Boolean(item));
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

export type ClusterLifecycleMode = 'imported' | 'registered' | 'provisioned' | string;
export type ClusterLifecycleStatus =
  | 'pending'
  | 'active'
  | 'degraded'
  | 'upgrading'
  | 'disabled'
  | 'retiring'
  | 'retired'
  | 'failed'
  | string;
export type RegistrationStatus =
  | 'not_required'
  | 'pending'
  | 'issued'
  | 'connected'
  | 'failed'
  | string;
export type ClusterHealthStatus = 'healthy' | 'warning' | 'critical' | 'unknown' | string;
export type LifecycleOperationStatus =
  | 'pending'
  | 'running'
  | 'partially_succeeded'
  | 'succeeded'
  | 'failed'
  | 'canceled'
  | 'blocked'
  | string;
export type LifecycleRiskLevel = 'low' | 'medium' | 'high' | 'critical' | string;
export type ClusterDriverStatus = 'draft' | 'active' | 'deprecated' | 'disabled' | string;
export type CapabilitySupportLevel = 'native' | 'extended' | 'partial' | 'unsupported' | string;
export type CapabilityCompatibilityStatus =
  | 'compatible'
  | 'conditional'
  | 'incompatible'
  | string;
export type UpgradePlanStatus =
  | 'draft'
  | 'approved'
  | 'running'
  | 'succeeded'
  | 'failed'
  | 'canceled'
  | string;
export type UpgradePrecheckStatus = 'pending' | 'passed' | 'warning' | 'failed' | string;
export type NodePoolRole = 'control-plane' | 'worker' | 'mixed' | string;
export type NodePoolStatus =
  | 'pending'
  | 'active'
  | 'scaling'
  | 'upgrading'
  | 'degraded'
  | 'failed'
  | string;
export type CapabilityDomain =
  | 'network'
  | 'storage'
  | 'identity'
  | 'observability'
  | 'security'
  | 'backup'
  | 'release'
  | string;

export type ClusterLifecycleRecord = {
  id: string;
  name: string;
  displayName?: string;
  lifecycleMode?: ClusterLifecycleMode;
  infrastructureType?: string;
  driverRef?: string;
  driverVersion?: string;
  workspaceId?: string;
  projectId?: string;
  status?: ClusterLifecycleStatus;
  registrationStatus?: RegistrationStatus;
  healthStatus?: ClusterHealthStatus;
  kubernetesVersion?: string;
  targetVersion?: string;
  nodePoolSummary?: string;
  lastValidationStatus?: string;
  lastValidationAt?: string;
  lastOperationId?: string;
  retirementReason?: string;
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type LifecycleOperation = {
  id: string;
  clusterId?: string;
  operationType?: string;
  status?: LifecycleOperationStatus;
  riskLevel?: LifecycleRiskLevel;
  resultSummary?: string;
  failureReason?: string;
  startedAt?: string;
  completedAt?: string;
};

export type RegistrationBundle = {
  clusterId: string;
  registrationToken?: string;
  commandSnippet?: string;
  expiresAt?: string;
};

export type ClusterDriverVersion = {
  id: string;
  driverKey: string;
  version: string;
  displayName?: string;
  providerType?: string;
  status?: ClusterDriverStatus;
  capabilityProfileVersion?: string;
  schemaVersion?: string;
  releaseNotes?: string;
};

export type CapabilityMatrixEntry = {
  id: string;
  ownerType?: string;
  ownerRef?: string;
  capabilityDomain?: CapabilityDomain;
  supportLevel?: CapabilitySupportLevel;
  compatibilityStatus?: CapabilityCompatibilityStatus;
  constraintsSummary?: string;
  recommendedFor?: string;
  updatedAt?: string;
};

export type ClusterTemplate = {
  id: string;
  name: string;
  description?: string;
  infrastructureType?: string;
  driverKey?: string;
  driverVersionRange?: string;
  requiredCapabilities?: string[];
  parameterSchemaRef?: string;
  defaultValuesRef?: string;
  status?: string;
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type ValidationResult = {
  overallStatus?: string;
  blockers?: string[];
  warnings?: string[];
  passedChecks?: string[];
};

export type UpgradePlan = {
  id: string;
  clusterId?: string;
  fromVersion?: string;
  toVersion?: string;
  windowStart?: string;
  windowEnd?: string;
  precheckStatus?: UpgradePrecheckStatus;
  impactSummary?: string;
  status?: UpgradePlanStatus;
  lastOperationId?: string;
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type NodePoolProfile = {
  id: string;
  clusterId?: string;
  name: string;
  role?: NodePoolRole;
  desiredCount?: number;
  currentCount?: number;
  minCount?: number;
  maxCount?: number;
  version?: string;
  zoneRefs?: string[];
  status?: NodePoolStatus;
  lastOperationId?: string;
  updatedAt?: string;
};

export type LifecycleAuditEvent = {
  id: string;
  action?: string;
  actorUserId?: string;
  clusterId?: string;
  targetType?: string;
  targetRef?: string;
  outcome?: string;
  occurredAt?: string;
  details?: Record<string, unknown>;
};

export type ClusterLifecycleDetail = ClusterLifecycleRecord & {
  cluster?: ClusterLifecycleRecord;
  nodePools?: NodePoolProfile[];
  upgradePlans?: UpgradePlan[];
  recentOperations?: LifecycleOperation[];
};

export type ClusterLifecycleListQuery = {
  status?: string;
  infrastructureType?: string;
  driverKey?: string;
  keyword?: string;
};

export type ImportClusterRequest = {
  name: string;
  infrastructureType: string;
  accessEndpoint: string;
  credentialRef?: string;
  workspaceId?: string;
  projectId?: string;
};

export type RegisterClusterRequest = {
  name: string;
  infrastructureType: string;
  workspaceId?: string;
  projectId?: string;
};

export type CreateClusterRequest = {
  name: string;
  infrastructureType: string;
  driverRef: string;
  templateId: string;
  parameterOverrides?: Record<string, unknown>;
  workspaceId?: string;
  projectId?: string;
};

export type ValidationRequest = {
  templateId?: string;
  driverRef?: string;
  targetVersion?: string;
};

export type CreateUpgradePlanRequest = {
  toVersion: string;
  windowStart?: string;
  windowEnd?: string;
  reason?: string;
};

export type ScaleNodePoolRequest = {
  desiredCount: number;
  reason?: string;
};

export type DisableClusterRequest = {
  reason: string;
  confirmation: boolean;
};

export type RetireClusterRequest = {
  reason: string;
  confirmation: boolean;
  evidenceNote?: string;
};

export type CreateDriverRequest = {
  driverKey: string;
  version: string;
  providerType: string;
  displayName?: string;
};

export type CreateTemplateRequest = {
  name: string;
  infrastructureType: string;
  driverKey: string;
  driverVersionRange?: string;
  requiredCapabilities?: string[];
};

export type TemplateValidationRequest = {
  driverRef?: string;
  environmentSnapshot?: Record<string, unknown>;
};

export type LifecycleAuditQuery = {
  actorUserId?: string;
  action?: string;
  outcome?: string;
  infrastructureType?: string;
};

const mapClusterLifecycleRecord = (
  value: UnknownRecord,
  index: number
): ClusterLifecycleRecord => {
  const name = toText(pick(value, ['name', 'displayName'])) || `cluster-${index + 1}`;
  return {
    id: toText(pick(value, ['id', 'clusterId'])) || name,
    name,
    displayName: toText(pick(value, ['displayName'])),
    lifecycleMode: toText(pick(value, ['lifecycleMode'])),
    infrastructureType: toText(pick(value, ['infrastructureType'])),
    driverRef: toText(pick(value, ['driverRef'])),
    driverVersion: toText(pick(value, ['driverVersion'])),
    workspaceId: toText(pick(value, ['workspaceId'])),
    projectId: toText(pick(value, ['projectId'])),
    status: toText(pick(value, ['status'])),
    registrationStatus: toText(pick(value, ['registrationStatus'])),
    healthStatus: toText(pick(value, ['healthStatus'])),
    kubernetesVersion: toText(pick(value, ['kubernetesVersion'])),
    targetVersion: toText(pick(value, ['targetVersion'])),
    nodePoolSummary: toText(pick(value, ['nodePoolSummary'])),
    lastValidationStatus: toText(pick(value, ['lastValidationStatus'])),
    lastValidationAt: toText(pick(value, ['lastValidationAt'])),
    lastOperationId: toText(pick(value, ['lastOperationId'])),
    retirementReason: toText(pick(value, ['retirementReason'])),
    createdBy: toText(pick(value, ['createdBy'])),
    createdAt: toText(pick(value, ['createdAt'])),
    updatedAt: toText(pick(value, ['updatedAt']))
  };
};

const mapOperation = (value: unknown): LifecycleOperation => {
  const record = isRecord(value) ? value : {};
  return {
    id: toText(pick(record, ['id'])) || 'operation-unknown',
    clusterId: toText(pick(record, ['clusterId'])),
    operationType: toText(pick(record, ['operationType'])),
    status: toText(pick(record, ['status'])),
    riskLevel: toText(pick(record, ['riskLevel'])),
    resultSummary: toText(pick(record, ['resultSummary'])),
    failureReason: toText(pick(record, ['failureReason'])),
    startedAt: toText(pick(record, ['startedAt'])),
    completedAt: toText(pick(record, ['completedAt']))
  };
};

const mapRegistrationBundle = (value: unknown): RegistrationBundle => {
  const record = isRecord(value) ? value : {};
  return {
    clusterId: toText(pick(record, ['clusterId'])) || 'cluster-unknown',
    registrationToken: toText(pick(record, ['registrationToken'])),
    commandSnippet: toText(pick(record, ['commandSnippet'])),
    expiresAt: toText(pick(record, ['expiresAt']))
  };
};

const mapDriverVersion = (value: UnknownRecord, index: number): ClusterDriverVersion => {
  const driverKey = toText(pick(value, ['driverKey'])) || `driver-${index + 1}`;
  return {
    id: toText(pick(value, ['id'])) || driverKey,
    driverKey,
    version: toText(pick(value, ['version'])) || UNKNOWN_TEXT,
    displayName: toText(pick(value, ['displayName'])),
    providerType: toText(pick(value, ['providerType'])),
    status: toText(pick(value, ['status'])),
    capabilityProfileVersion: toText(pick(value, ['capabilityProfileVersion'])),
    schemaVersion: toText(pick(value, ['schemaVersion'])),
    releaseNotes: toText(pick(value, ['releaseNotes']))
  };
};

const mapCapabilityMatrixEntry = (
  value: UnknownRecord,
  index: number
): CapabilityMatrixEntry => ({
  id: toText(pick(value, ['id'])) || `capability-${index + 1}`,
  ownerType: toText(pick(value, ['ownerType'])),
  ownerRef: toText(pick(value, ['ownerRef'])),
  capabilityDomain: toText(pick(value, ['capabilityDomain'])),
  supportLevel: toText(pick(value, ['supportLevel'])),
  compatibilityStatus: toText(pick(value, ['compatibilityStatus'])),
  constraintsSummary: toText(pick(value, ['constraintsSummary'])),
  recommendedFor: toText(pick(value, ['recommendedFor'])),
  updatedAt: toText(pick(value, ['updatedAt']))
});

const mapTemplate = (value: UnknownRecord, index: number): ClusterTemplate => {
  const name = toText(pick(value, ['name'])) || `template-${index + 1}`;
  return {
    id: toText(pick(value, ['id'])) || name,
    name,
    description: toText(pick(value, ['description'])),
    infrastructureType: toText(pick(value, ['infrastructureType'])),
    driverKey: toText(pick(value, ['driverKey'])),
    driverVersionRange: toText(pick(value, ['driverVersionRange'])),
    requiredCapabilities: toStringArray(pick(value, ['requiredCapabilities'])),
    parameterSchemaRef: toText(pick(value, ['parameterSchemaRef'])),
    defaultValuesRef: toText(pick(value, ['defaultValuesRef'])),
    status: toText(pick(value, ['status'])),
    createdBy: toText(pick(value, ['createdBy'])),
    createdAt: toText(pick(value, ['createdAt'])),
    updatedAt: toText(pick(value, ['updatedAt']))
  };
};

const mapValidationResult = (value: unknown): ValidationResult => {
  const record = isRecord(value) ? value : {};
  return {
    overallStatus: toText(pick(record, ['overallStatus'])),
    blockers: toStringArray(pick(record, ['blockers'])),
    warnings: toStringArray(pick(record, ['warnings'])),
    passedChecks: toStringArray(pick(record, ['passedChecks']))
  };
};

const mapUpgradePlan = (value: unknown): UpgradePlan => {
  const record = isRecord(value) ? value : {};
  return {
    id: toText(pick(record, ['id'])) || 'plan-unknown',
    clusterId: toText(pick(record, ['clusterId'])),
    fromVersion: toText(pick(record, ['fromVersion'])),
    toVersion: toText(pick(record, ['toVersion'])),
    windowStart: toText(pick(record, ['windowStart'])),
    windowEnd: toText(pick(record, ['windowEnd'])),
    precheckStatus: toText(pick(record, ['precheckStatus'])),
    impactSummary: toText(pick(record, ['impactSummary'])),
    status: toText(pick(record, ['status'])),
    lastOperationId: toText(pick(record, ['lastOperationId'])),
    createdBy: toText(pick(record, ['createdBy'])),
    createdAt: toText(pick(record, ['createdAt'])),
    updatedAt: toText(pick(record, ['updatedAt']))
  };
};

const mapNodePool = (value: UnknownRecord, index: number): NodePoolProfile => {
  const name = toText(pick(value, ['name'])) || `node-pool-${index + 1}`;
  return {
    id: toText(pick(value, ['id'])) || name,
    clusterId: toText(pick(value, ['clusterId'])),
    name,
    role: toText(pick(value, ['role'])),
    desiredCount: toNumber(pick(value, ['desiredCount'])),
    currentCount: toNumber(pick(value, ['currentCount'])),
    minCount: toNumber(pick(value, ['minCount'])),
    maxCount: toNumber(pick(value, ['maxCount'])),
    version: toText(pick(value, ['version'])),
    zoneRefs: toStringArray(pick(value, ['zoneRefs'])),
    status: toText(pick(value, ['status'])),
    lastOperationId: toText(pick(value, ['lastOperationId'])),
    updatedAt: toText(pick(value, ['updatedAt']))
  };
};

const mapAuditEvent = (value: UnknownRecord, index: number): LifecycleAuditEvent => {
  const details = isRecord(pick(value, ['details'])) ? (pick(value, ['details']) as UnknownRecord) : {};
  return {
    id:
      toText(pick(value, ['id'])) ||
      `${toText(pick(value, ['occurredAt'])) || 'time-unknown'}-${toText(pick(value, ['action'])) || index}`,
    action: toText(pick(value, ['action'])),
    actorUserId: toText(pick(value, ['actorUserId'])),
    clusterId: toText(pick(value, ['clusterId'])),
    targetType: toText(pick(value, ['targetType'])),
    targetRef: toText(pick(value, ['targetRef'])),
    outcome: toText(pick(value, ['outcome'])),
    occurredAt: toText(pick(value, ['occurredAt'])),
    details
  };
};

const mapClusterLifecycleDetail = (value: unknown): ClusterLifecycleDetail => {
  const record = isRecord(value) ? value : {};
  const clusterRecord = isRecord(pick(record, ['cluster']))
    ? (pick(record, ['cluster']) as UnknownRecord)
    : record;

  const cluster = mapClusterLifecycleRecord(clusterRecord, 0);
  const nodePools = Array.isArray(record.nodePools)
    ? record.nodePools.filter(isRecord).map(mapNodePool)
    : [];
  const upgradePlans = Array.isArray(record.upgradePlans)
    ? record.upgradePlans.map(mapUpgradePlan)
    : [];
  const recentOperations = Array.isArray(record.recentOperations)
    ? record.recentOperations.map(mapOperation)
    : undefined;

  return {
    ...cluster,
    cluster,
    nodePools,
    upgradePlans,
    recentOperations
  };
};

export const clusterLifecycleQueryKeys = {
  all: ['clusterLifecycle'] as const,
  clusters: (scope?: string) => ['clusterLifecycle', 'clusters', scope ?? 'default'] as const,
  clusterDetail: (clusterId?: string) =>
    ['clusterLifecycle', 'clusterDetail', clusterId ?? 'unknown'] as const,
  drivers: (scope?: string) => ['clusterLifecycle', 'drivers', scope ?? 'default'] as const,
  templates: (scope?: string) => ['clusterLifecycle', 'templates', scope ?? 'default'] as const,
  capabilityMatrix: (scope?: string) =>
    ['clusterLifecycle', 'capabilityMatrix', scope ?? 'default'] as const,
  nodePools: (clusterId?: string) =>
    ['clusterLifecycle', 'nodePools', clusterId ?? 'unknown'] as const,
  audit: (scope?: string) => ['clusterLifecycle', 'audit', scope ?? 'default'] as const
};

export const listClusterLifecycleRecords = async (
  query: ClusterLifecycleListQuery = {}
): Promise<{ items: ClusterLifecycleRecord[] }> => {
  const response = await fetchJSON<unknown>(withQuery('/cluster-lifecycle/clusters', query));
  return normalizeListResponse(response, mapClusterLifecycleRecord);
};

export const importCluster = async (
  payload: ImportClusterRequest
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/clusters/import', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapOperation(response);
};

export const registerCluster = async (
  payload: RegisterClusterRequest
): Promise<RegistrationBundle> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/clusters/register', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapRegistrationBundle(response);
};

export const getClusterLifecycleDetail = async (
  clusterId: string
): Promise<ClusterLifecycleDetail> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}`
  );
  return mapClusterLifecycleDetail(response);
};

export const createCluster = async (
  payload: CreateClusterRequest
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/clusters', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapOperation(response);
};

export const validateClusterChange = async (
  clusterId: string,
  payload: ValidationRequest
): Promise<ValidationResult> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/validate`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapValidationResult(response);
};

export const createUpgradePlan = async (
  clusterId: string,
  payload: CreateUpgradePlanRequest
): Promise<UpgradePlan> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/upgrade-plans`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapUpgradePlan(response);
};

export const executeUpgradePlan = async (
  clusterId: string,
  planId: string
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/upgrade-plans/${encodeURIComponent(planId)}/execute`,
    { method: 'POST' }
  );
  return mapOperation(response);
};

export const listNodePools = async (
  clusterId: string
): Promise<{ items: NodePoolProfile[] }> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/node-pools`
  );
  return normalizeListResponse(response, mapNodePool);
};

export const scaleNodePool = async (
  clusterId: string,
  nodePoolId: string,
  payload: ScaleNodePoolRequest
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/node-pools/${encodeURIComponent(nodePoolId)}/scale`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapOperation(response);
};

export const disableCluster = async (
  clusterId: string,
  payload: DisableClusterRequest
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/disable`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapOperation(response);
};

export const retireCluster = async (
  clusterId: string,
  payload: RetireClusterRequest
): Promise<LifecycleOperation> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/clusters/${encodeURIComponent(clusterId)}/retire`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapOperation(response);
};

export const listClusterDrivers = async (): Promise<{ items: ClusterDriverVersion[] }> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/drivers');
  return normalizeListResponse(response, mapDriverVersion);
};

export const createClusterDriver = async (
  payload: CreateDriverRequest
): Promise<ClusterDriverVersion> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/drivers', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapDriverVersion(isRecord(response) ? response : {}, 0);
};

export const listDriverCapabilities = async (
  driverId: string
): Promise<{ items: CapabilityMatrixEntry[] }> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/drivers/${encodeURIComponent(driverId)}/capabilities`
  );
  return normalizeListResponse(response, mapCapabilityMatrixEntry);
};

export const listClusterTemplates = async (): Promise<{ items: ClusterTemplate[] }> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/templates');
  return normalizeListResponse(response, mapTemplate);
};

export const createClusterTemplate = async (
  payload: CreateTemplateRequest
): Promise<ClusterTemplate> => {
  const response = await fetchJSON<unknown>('/cluster-lifecycle/templates', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapTemplate(isRecord(response) ? response : {}, 0);
};

export const validateClusterTemplate = async (
  templateId: string,
  payload: TemplateValidationRequest
): Promise<ValidationResult> => {
  const response = await fetchJSON<unknown>(
    `/cluster-lifecycle/templates/${encodeURIComponent(templateId)}/validate`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return mapValidationResult(response);
};

export const listLifecycleAuditEvents = async (
  query: LifecycleAuditQuery = {}
): Promise<{ items: LifecycleAuditEvent[] }> => {
  const response = await fetchJSON<unknown>(withQuery('/audit/cluster-lifecycle/events', query));
  return normalizeListResponse(response, mapAuditEvent);
};

export const hasBlockingValidationIssue = (value?: ValidationResult | null): boolean =>
  (value?.blockers?.length ?? 0) > 0 || value?.overallStatus === 'failed';

export const isDestructiveLifecycleAction = (action?: string): boolean =>
  action === 'disable' || action === 'retire';

export const isConfirmedActionPayload = (
  value: DisableClusterRequest | RetireClusterRequest
): boolean => toBoolean(value.confirmation) === true;

export const flattenCapabilityDomains = (items: CapabilityMatrixEntry[]): CapabilityDomain[] => {
  const domains = new Set<CapabilityDomain>();
  items.forEach((item) => {
    if (item.capabilityDomain) {
      domains.add(item.capabilityDomain);
    }
  });
  return Array.from(domains);
};

export const groupCapabilityMatrixByDomain = (
  items: CapabilityMatrixEntry[]
): Record<string, CapabilityMatrixEntry[]> => {
  return items.reduce<Record<string, CapabilityMatrixEntry[]>>((accumulator, item) => {
    const domain = item.capabilityDomain || UNGROUPED_DOMAIN;
    accumulator[domain] = accumulator[domain] || [];
    accumulator[domain].push(item);
    return accumulator;
  }, {});
};

export const extractLifecycleSummary = (clusters: ClusterLifecycleRecord[]) => ({
  total: clusters.length,
  active: clusters.filter((cluster) => cluster.status === 'active').length,
  pending: clusters.filter((cluster) => cluster.status === 'pending').length,
  degraded: clusters.filter((cluster) => cluster.status === 'degraded').length,
  retiring: clusters.filter((cluster) => cluster.status === 'retiring').length
});

export const inferClusterDisplayName = (cluster?: ClusterLifecycleRecord | null) =>
  cluster?.displayName || cluster?.name || cluster?.id || '未命名集群';

export const extractRecentOperationMessage = (operation?: LifecycleOperation | null) =>
  operation?.failureReason || operation?.resultSummary || '动作已受理，等待后端执行结果。';

export const listDraftCapabilityEntries = (items: CapabilityMatrixEntry[]) =>
  items.filter((item) => item.compatibilityStatus !== 'compatible');

export const listProvisionableTemplates = (items: ClusterTemplate[]) =>
  items.filter((item) => item.status !== 'disabled');

export const findClusterById = (
  clusters: ClusterLifecycleRecord[],
  clusterId?: string | null
): ClusterLifecycleRecord | undefined =>
  clusters.find((cluster) => cluster.id === clusterId);

export const mapAuditTargetLabel = (event: LifecycleAuditEvent) =>
  [event.targetType, event.targetRef].filter(Boolean).join(':') || '-';

export const getDefaultEnvironmentSnapshot = (cluster?: ClusterLifecycleRecord | null) => ({
  infrastructureType: cluster?.infrastructureType,
  kubernetesVersion: cluster?.kubernetesVersion,
  driverRef: cluster?.driverRef
});

export const extractClusterOperationState = (
  cluster?: ClusterLifecycleRecord | null
): 'idle' | 'running' | 'blocked' => {
  if (!cluster) {
    return 'idle';
  }
  if (cluster.status === 'upgrading' || cluster.status === 'retiring') {
    return 'running';
  }
  if (cluster.status === 'retired' || cluster.status === 'failed') {
    return 'blocked';
  }
  return 'idle';
};

export const listCapabilityOwners = (items: CapabilityMatrixEntry[]) => {
  const owners = new Set<string>();
  items.forEach((item) => {
    const value = item.ownerRef || item.ownerType;
    if (value) {
      owners.add(value);
    }
  });
  return Array.from(owners);
};

export const listCapabilityWarnings = (items: CapabilityMatrixEntry[]) =>
  items.filter((item) => item.compatibilityStatus === 'conditional');

export const listCapabilityBlockers = (items: CapabilityMatrixEntry[]) =>
  items.filter((item) => item.compatibilityStatus === 'incompatible');

export const mergeCapabilityEntries = (
  left: CapabilityMatrixEntry[],
  right: CapabilityMatrixEntry[]
) => [...left, ...right];

export const normalizeNodePoolTarget = (
  nodePool?: NodePoolProfile | null
): number | undefined => nodePool?.desiredCount ?? nodePool?.currentCount;

export const toNodePoolScalePayload = (
  desiredCount: number,
  reason?: string
): ScaleNodePoolRequest => ({
  desiredCount,
  reason
});
