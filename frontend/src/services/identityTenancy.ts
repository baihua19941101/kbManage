import { buildScopeQueryKey, fetchJSON } from '@/services/api/client';

type UnknownRecord = Record<string, unknown>;
type QueryValue = string | number | boolean | undefined | null;
export type ScopeSelection = Record<string, unknown>;

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
  if (typeof value === 'string') {
    if (value === 'true') {
      return true;
    }
    if (value === 'false') {
      return false;
    }
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
) => {
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

export type IdentitySourceStatus = 'active' | 'degraded' | 'draft' | 'failed' | string;
export type LoginMode = 'local' | 'external' | 'mixed' | string;
export type SessionStatus = 'active' | 'revoked' | 'expired' | 'blocked' | string;
export type RiskLevel = 'low' | 'medium' | 'high' | 'critical' | string;

export type IdentitySource = {
  id: string;
  name: string;
  sourceType?: string;
  status?: IdentitySourceStatus;
  loginMode?: LoginMode;
  scopeMode?: string;
  syncState?: string;
  lastCheckedAt?: string;
  accountCount?: number;
};

export type CreateIdentitySourcePayload = {
  name: string;
  sourceType: string;
  loginMode: LoginMode;
  scopeMode: string;
};

export type OrganizationUnit = {
  id: string;
  unitType?: string;
  name: string;
  description?: string;
  parentUnitId?: string;
  status?: string;
  identitySourceId?: string;
  memberCount?: number;
};

export type CreateOrganizationUnitPayload = {
  unitType: string;
  name: string;
  description?: string;
  parentUnitId?: string;
};

export type TenantScopeMapping = {
  id: string;
  unitId?: string;
  scopeType?: string;
  scopeRef?: string;
  inheritanceMode?: string;
  status?: string;
  conflictSummary?: string;
};

export type CreateTenantScopeMappingPayload = {
  scopeType: string;
  scopeRef: string;
  inheritanceMode: string;
};

export type RoleDefinition = {
  id: string;
  name: string;
  roleLevel?: string;
  description?: string;
  permissionSummary?: string;
  inheritancePolicy?: string;
  delegable?: boolean;
  status?: string;
};

export type CreateRoleDefinitionPayload = {
  name: string;
  roleLevel: string;
  description?: string;
  permissionSummary: string;
  inheritancePolicy: string;
  delegable?: boolean;
};

export type RoleAssignment = {
  id: string;
  subjectType?: string;
  subjectRef?: string;
  roleDefinitionId?: string;
  roleDefinitionName?: string;
  scopeType?: string;
  scopeRef?: string;
  sourceType?: string;
  validFrom?: string;
  validUntil?: string;
  status?: string;
};

export type CreateRoleAssignmentPayload = {
  subjectType: string;
  subjectRef: string;
  roleDefinitionId: string;
  scopeType: string;
  scopeRef: string;
  validUntil?: string;
};

export type DelegationGrant = {
  id: string;
  grantorRef?: string;
  delegateRef?: string;
  allowedRoleLevels: string[];
  status?: string;
  validFrom?: string;
  validUntil?: string;
  reason?: string;
};

export type CreateDelegationGrantPayload = {
  grantorRef: string;
  delegateRef: string;
  allowedRoleLevels: string[];
  validFrom: string;
  validUntil: string;
  reason?: string;
};

export type SessionRecord = {
  id: string;
  userId?: string;
  username?: string;
  identitySourceId?: string;
  loginMethod?: string;
  status?: SessionStatus;
  riskLevel?: RiskLevel;
  riskSummary?: string;
  lastSeenAt?: string;
};

export type AccessRiskSnapshot = {
  id: string;
  subjectType?: string;
  subjectRef?: string;
  riskType?: string;
  severity?: RiskLevel;
  summary?: string;
  recommendedAction?: string;
  status?: string;
  generatedAt?: string;
};

export type IdentityAuditEvent = {
  id: string;
  action?: string;
  actorUserId?: string;
  targetType?: string;
  targetRef?: string;
  outcome?: string;
  occurredAt?: string;
};

export type IdentitySourceListQuery = {
  sourceType?: string;
  status?: string;
};

export type OrganizationUnitListQuery = {
  unitType?: string;
  parentUnitId?: string;
};

export type RoleAssignmentListQuery = {
  subjectRef?: string;
  scopeType?: string;
};

export type SessionRecordListQuery = {
  status?: string;
  riskLevel?: string;
};

export type AccessRiskListQuery = {
  subjectType?: string;
  severity?: string;
};

export type IdentityAuditListQuery = {
  action?: string;
  outcome?: string;
  targetType?: string;
};

export const identityTenancyQueryKeys = {
  all: ['identity-tenancy'] as const,
  sources: (scope?: string) => ['identity-tenancy', 'sources', scope ?? 'all'] as const,
  organizations: (scope?: string) =>
    ['identity-tenancy', 'organizations', scope ?? 'all'] as const,
  mappings: (unitId?: string) => ['identity-tenancy', 'mappings', unitId ?? 'all'] as const,
  roles: (scope?: string) => ['identity-tenancy', 'roles', scope ?? 'all'] as const,
  assignments: (scope?: string) => ['identity-tenancy', 'assignments', scope ?? 'all'] as const,
  delegations: () => ['identity-tenancy', 'delegations'] as const,
  sessions: (scope?: string) => ['identity-tenancy', 'sessions', scope ?? 'all'] as const,
  risks: (scope?: string) => ['identity-tenancy', 'risks', scope ?? 'all'] as const,
  audit: (scope?: string) => ['identity-tenancy', 'audit', scope ?? 'all'] as const
};

const mapIdentitySource = (record: UnknownRecord): IdentitySource => ({
  id: toText(pick(record, ['id', 'sourceId'])) || 'unknown-source',
  name: toText(record.name) || '未命名身份源',
  sourceType: toText(record.sourceType),
  status: toText(record.status),
  loginMode: toText(record.loginMode),
  scopeMode: toText(record.scopeMode),
  syncState: toText(record.syncState),
  lastCheckedAt: toText(record.lastCheckedAt),
  accountCount: toNumber(record.accountCount)
});

const mapOrganizationUnit = (record: UnknownRecord): OrganizationUnit => ({
  id: toText(pick(record, ['id', 'unitId'])) || 'unknown-unit',
  unitType: toText(record.unitType),
  name: toText(record.name) || '未命名单元',
  description: toText(record.description),
  parentUnitId: toText(record.parentUnitId),
  status: toText(record.status),
  identitySourceId: toText(record.identitySourceId),
  memberCount: toNumber(record.memberCount)
});

const mapTenantScopeMapping = (record: UnknownRecord): TenantScopeMapping => ({
  id: toText(pick(record, ['id', 'mappingId'])) || 'unknown-mapping',
  unitId: toText(record.unitId),
  scopeType: toText(record.scopeType),
  scopeRef: toText(record.scopeRef),
  inheritanceMode: toText(record.inheritanceMode),
  status: toText(record.status),
  conflictSummary: toText(record.conflictSummary)
});

const mapRoleDefinition = (record: UnknownRecord): RoleDefinition => ({
  id: toText(pick(record, ['id', 'roleDefinitionId'])) || 'unknown-role',
  name: toText(record.name) || '未命名角色',
  roleLevel: toText(record.roleLevel),
  description: toText(record.description),
  permissionSummary: toText(record.permissionSummary),
  inheritancePolicy: toText(record.inheritancePolicy),
  delegable: toBoolean(record.delegable),
  status: toText(record.status)
});

const mapRoleAssignment = (record: UnknownRecord): RoleAssignment => ({
  id: toText(pick(record, ['id', 'assignmentId'])) || 'unknown-assignment',
  subjectType: toText(record.subjectType),
  subjectRef: toText(record.subjectRef),
  roleDefinitionId: toText(record.roleDefinitionId),
  roleDefinitionName: toText(record.roleDefinitionName),
  scopeType: toText(record.scopeType),
  scopeRef: toText(record.scopeRef),
  sourceType: toText(record.sourceType),
  validFrom: toText(record.validFrom),
  validUntil: toText(record.validUntil),
  status: toText(record.status)
});

const mapDelegationGrant = (record: UnknownRecord): DelegationGrant => ({
  id: toText(pick(record, ['id', 'grantId'])) || 'unknown-delegation',
  grantorRef: toText(record.grantorRef),
  delegateRef: toText(record.delegateRef),
  allowedRoleLevels: toStringArray(record.allowedRoleLevels),
  status: toText(record.status),
  validFrom: toText(record.validFrom),
  validUntil: toText(record.validUntil),
  reason: toText(record.reason)
});

const mapSessionRecord = (record: UnknownRecord): SessionRecord => ({
  id: toText(pick(record, ['id', 'sessionId'])) || 'unknown-session',
  userId: toText(record.userId),
  username: toText(record.username),
  identitySourceId: toText(record.identitySourceId),
  loginMethod: toText(record.loginMethod),
  status: toText(record.status),
  riskLevel: toText(record.riskLevel),
  riskSummary: toText(record.riskSummary),
  lastSeenAt: toText(record.lastSeenAt)
});

const mapAccessRiskSnapshot = (record: UnknownRecord): AccessRiskSnapshot => ({
  id: toText(pick(record, ['id', 'riskId'])) || 'unknown-risk',
  subjectType: toText(record.subjectType),
  subjectRef: toText(record.subjectRef),
  riskType: toText(record.riskType),
  severity: toText(record.severity),
  summary: toText(record.summary),
  recommendedAction: toText(record.recommendedAction),
  status: toText(record.status),
  generatedAt: toText(record.generatedAt)
});

const mapIdentityAuditEvent = (record: UnknownRecord): IdentityAuditEvent => ({
  id: toText(pick(record, ['id', 'eventId'])) || 'unknown-event',
  action: toText(record.action),
  actorUserId: toText(record.actorUserId),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  outcome: toText(record.outcome),
  occurredAt: toText(record.occurredAt)
});

export const listIdentitySources = async (query: IdentitySourceListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/sources', query));
  return normalizeListResponse(response, mapIdentitySource);
};

export const getIdentitySourceDetail = async (sourceId: string) => {
  const response = await fetchJSON<unknown>(`/identity/sources/${sourceId}`);
  return mapIdentitySource(isRecord(response) ? response : {});
};

export const createIdentitySource = async (payload: CreateIdentitySourcePayload) => {
  const response = await fetchJSON<unknown>('/identity/sources', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapIdentitySource(isRecord(response) ? response : {});
};

export const updatePreferredLoginMode = async (loginMode: LoginMode) => {
  const response = await fetchJSON<unknown>('/identity/login-mode', {
    method: 'POST',
    body: JSON.stringify({ loginMode })
  });
  return { loginMode: toText(isRecord(response) ? response.loginMode : undefined) ?? loginMode };
};

export const listOrganizationUnits = async (query: OrganizationUnitListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/organizations', query));
  return normalizeListResponse(response, mapOrganizationUnit);
};

export const createOrganizationUnit = async (payload: CreateOrganizationUnitPayload) => {
  const response = await fetchJSON<unknown>('/identity/organizations', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapOrganizationUnit(isRecord(response) ? response : {});
};

export const listTenantScopeMappings = async (unitId: string) => {
  const response = await fetchJSON<unknown>(`/identity/organizations/${unitId}/mappings`);
  return normalizeListResponse(response, mapTenantScopeMapping);
};

export const createTenantScopeMapping = async (
  unitId: string,
  payload: CreateTenantScopeMappingPayload
) => {
  const response = await fetchJSON<unknown>(`/identity/organizations/${unitId}/mappings`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapTenantScopeMapping(isRecord(response) ? response : {});
};

export const listRoleDefinitions = async (query: { roleLevel?: string } = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/roles', query));
  return normalizeListResponse(response, mapRoleDefinition);
};

export const createRoleDefinition = async (payload: CreateRoleDefinitionPayload) => {
  const response = await fetchJSON<unknown>('/identity/roles', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapRoleDefinition(isRecord(response) ? response : {});
};

export const listRoleAssignments = async (query: RoleAssignmentListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/assignments', query));
  return normalizeListResponse(response, mapRoleAssignment);
};

export const createRoleAssignment = async (payload: CreateRoleAssignmentPayload) => {
  const response = await fetchJSON<unknown>('/identity/assignments', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapRoleAssignment(isRecord(response) ? response : {});
};

export const listDelegationGrants = async () => {
  const response = await fetchJSON<unknown>('/identity/delegations');
  return normalizeListResponse(response, mapDelegationGrant);
};

export const createDelegationGrant = async (payload: CreateDelegationGrantPayload) => {
  const response = await fetchJSON<unknown>('/identity/delegations', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return mapDelegationGrant(isRecord(response) ? response : {});
};

export const listSessionRecords = async (query: SessionRecordListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/sessions', query));
  return normalizeListResponse(response, mapSessionRecord);
};

export const revokeSessionRecord = async (sessionId: string) => {
  const response = await fetchJSON<unknown>(`/identity/sessions/${sessionId}/revoke`, {
    method: 'POST'
  });
  return mapSessionRecord(isRecord(response) ? response : {});
};

export const listAccessRisks = async (query: AccessRiskListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/identity/access-risks', query));
  return normalizeListResponse(response, mapAccessRiskSnapshot);
};

export const listIdentityGovernanceAuditEvents = async (query: IdentityAuditListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/audit/identity/events', query));
  return normalizeListResponse(response, mapIdentityAuditEvent);
};

export const summarizeScopeBoundary = (scopeType?: string, scopeRef?: string) =>
  [scopeType, scopeRef].filter(Boolean).join(' / ') || '未指定边界';

export const summarizeRoleBoundary = (assignment: RoleAssignment) =>
  `${assignment.subjectRef || '未标记主体'} · ${assignment.scopeType || '未标记范围'} · ${assignment.scopeRef || '未标记对象'}`;

export const identitySourceQueryScope = (query: IdentitySourceListQuery) =>
  buildScopeQueryKey([query.sourceType, query.status]);

export const organizationQueryScope = (query: OrganizationUnitListQuery) =>
  buildScopeQueryKey([query.unitType, query.parentUnitId]);

export const roleAssignmentQueryScope = (query: RoleAssignmentListQuery) =>
  buildScopeQueryKey([query.subjectRef, query.scopeType]);

export const sessionQueryScope = (query: SessionRecordListQuery) =>
  buildScopeQueryKey([query.status, query.riskLevel]);

export const accessRiskQueryScope = (query: AccessRiskListQuery) =>
  buildScopeQueryKey([query.subjectType, query.severity]);
