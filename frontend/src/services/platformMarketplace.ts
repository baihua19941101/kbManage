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

export type CatalogSource = {
  id: string;
  name: string;
  sourceType?: string;
  status?: string;
  syncStatus?: string;
  templateCount?: number;
  healthSummary?: string;
  visibleScope?: string;
  lastSyncAt?: string;
};

export type CatalogSourceListQuery = {
  sourceType?: string;
  status?: string;
};

export type CreateCatalogSourcePayload = {
  name: string;
  sourceType: string;
  endpoint?: string;
  visibleScope?: string;
};

export type TemplateListQuery = {
  category?: string;
  status?: string;
  sourceId?: string;
};

export type TemplateVersion = {
  id: string;
  version: string;
  status?: string;
  dependencySummary?: string;
  releaseNotes?: string;
  formSummary?: string;
  constraintSummary?: string;
  compatibilitySummary?: string;
};

export type ApplicationTemplate = {
  id: string;
  name: string;
  category?: string;
  sourceId?: string;
  sourceName?: string;
  status?: string;
  latestVersion?: string;
  scopeSummary?: string;
  dependencySummary?: string;
  releaseNoteSummary?: string;
  templateVersions?: TemplateVersion[];
};

export type TemplateReleaseScope = {
  id: string;
  templateId?: string;
  templateName?: string;
  version?: string;
  targetType?: string;
  targetRef?: string;
  status?: string;
  visibilitySummary?: string;
  releaseNotes?: string;
  restrictions?: string[];
  createdAt?: string;
};

export type CreateTemplateReleasePayload = {
  version: string;
  targetType: string;
  targetRef: string;
  visibilitySummary?: string;
  releaseNotes?: string;
};

export type InstallationRecord = {
  id: string;
  templateId?: string;
  templateName?: string;
  targetType?: string;
  targetRef?: string;
  currentVersion?: string;
  latestVersion?: string;
  status?: string;
  changeSummary?: string;
  offlineState?: string;
  updatedAt?: string;
};

export type InstallationListQuery = {
  targetType?: string;
  status?: string;
};

export type ExtensionListQuery = {
  status?: string;
  compatibilityStatus?: string;
};

export type ExtensionPackage = {
  id: string;
  name: string;
  version?: string;
  status?: string;
  visibilityScope?: string;
  permissionSummary?: string;
  compatibilityStatus?: string;
  impactSummary?: string;
  updatedAt?: string;
};

export type CreateExtensionPackagePayload = {
  name: string;
  version: string;
  visibilityScope: string;
  permissionSummary: string;
};

export type CompatibilityStatement = {
  id: string;
  extensionId?: string;
  platformVersion?: string;
  compatibilityStatus?: string;
  summary?: string;
  blockedReasons: string[];
  suggestedActions: string[];
  permissionImpact?: string;
};

export type PlatformMarketplaceAuditEvent = {
  id: string;
  action?: string;
  actorUserId?: string;
  targetType?: string;
  targetRef?: string;
  outcome?: string;
  occurredAt?: string;
};

export const platformMarketplaceQueryKeys = {
  all: ['platform-marketplace'] as const,
  catalogSources: (scope?: string) =>
    ['platform-marketplace', 'catalog-sources', scope ?? 'all'] as const,
  templates: (scope?: string) => ['platform-marketplace', 'templates', scope ?? 'all'] as const,
  templateDetail: (templateId?: string) =>
    ['platform-marketplace', 'template-detail', templateId ?? 'unknown'] as const,
  templateReleases: (templateId?: string) =>
    ['platform-marketplace', 'template-releases', templateId ?? 'all'] as const,
  installations: (scope?: string) =>
    ['platform-marketplace', 'installations', scope ?? 'all'] as const,
  extensions: (scope?: string) =>
    ['platform-marketplace', 'extensions', scope ?? 'all'] as const,
  extensionCompatibility: (extensionId?: string) =>
    ['platform-marketplace', 'extension-compatibility', extensionId ?? 'unknown'] as const,
  audit: (scope?: string) => ['platform-marketplace', 'audit', scope ?? 'all'] as const
};

const mapCatalogSource = (record: UnknownRecord): CatalogSource => ({
  id: toText(pick(record, ['id', 'sourceId'])) || 'unknown-source',
  name: toText(record.name) || '未命名目录来源',
  sourceType: toText(record.sourceType),
  status: toText(record.status),
  syncStatus: toText(pick(record, ['syncStatus', 'syncState'])),
  templateCount: toNumber(record.templateCount),
  healthSummary: toText(record.healthSummary),
  visibleScope: toText(pick(record, ['visibleScope', 'visibilityScope'])),
  lastSyncAt: toText(record.lastSyncAt)
});

const mapTemplateVersion = (record: UnknownRecord): TemplateVersion => ({
  id: toText(pick(record, ['id', 'versionId'])) || 'unknown-version',
  version: toText(record.version) || 'unknown',
  status: toText(record.status),
  dependencySummary: toText(record.dependencySummary),
  releaseNotes: toText(record.releaseNotes),
  formSummary: toText(record.formSummary),
  constraintSummary: toText(record.constraintSummary),
  compatibilitySummary: toText(record.compatibilitySummary)
});

const mapApplicationTemplate = (record: UnknownRecord): ApplicationTemplate => ({
  id: toText(pick(record, ['id', 'templateId'])) || 'unknown-template',
  name: toText(record.name) || '未命名模板',
  category: toText(record.category),
  sourceId: toText(pick(record, ['sourceId', 'catalogSourceId'])),
  sourceName: toText(record.sourceName),
  status: toText(pick(record, ['status', 'publishStatus'])),
  latestVersion: toText(record.latestVersion),
  scopeSummary: toText(record.scopeSummary),
  dependencySummary: toText(record.dependencySummary),
  releaseNoteSummary: toText(pick(record, ['releaseNoteSummary', 'releaseNotesSummary'])),
  templateVersions: Array.isArray(pick(record, ['templateVersions', 'versions']))
    ? (pick(record, ['templateVersions', 'versions']) as unknown[])
        .filter(isRecord)
        .map(mapTemplateVersion)
    : []
});

const mapTemplateReleaseScope = (record: UnknownRecord): TemplateReleaseScope => ({
  id: toText(pick(record, ['id', 'releaseId'])) || 'unknown-release',
  templateId: toText(record.templateId),
  templateName: toText(record.templateName),
  version: toText(record.version),
  targetType: toText(pick(record, ['targetType', 'scopeType'])),
  targetRef: toText(pick(record, ['targetRef', 'scopeRef'])),
  status: toText(record.status),
  visibilitySummary: toText(pick(record, ['visibilitySummary', 'visibilityMode'])),
  releaseNotes: toText(record.releaseNotes),
  restrictions: toStringArray(record.restrictions),
  createdAt: toText(pick(record, ['createdAt', 'publishedAt']))
});

const mapInstallationRecord = (record: UnknownRecord): InstallationRecord => ({
  id: toText(pick(record, ['id', 'installationId'])) || 'unknown-installation',
  templateId: toText(record.templateId),
  templateName: toText(record.templateName),
  targetType: toText(pick(record, ['targetType', 'scopeType'])),
  targetRef: toText(pick(record, ['targetRef', 'scopeRef'])),
  currentVersion: toText(pick(record, ['currentVersion', 'currentInstalledVersion'])),
  latestVersion: toText(pick(record, ['latestVersion', 'upgradeTargetVersion'])),
  status: toText(pick(record, ['status', 'lifecycleStatus'])),
  changeSummary: toText(record.changeSummary),
  offlineState: toText(record.offlineState),
  updatedAt: toText(pick(record, ['updatedAt', 'lastChangedAt']))
});

const mapExtensionPackage = (record: UnknownRecord): ExtensionPackage => ({
  id: toText(pick(record, ['id', 'extensionId'])) || 'unknown-extension',
  name: toText(record.name) || '未命名扩展',
  version: toText(record.version),
  status: toText(record.status),
  visibilityScope: toText(record.visibilityScope),
  permissionSummary: toText(pick(record, ['permissionSummary', 'permissionDeclaration'])),
  compatibilityStatus: toText(record.compatibilityStatus),
  impactSummary: toText(pick(record, ['impactSummary', 'entrySummary'])),
  updatedAt: toText(record.updatedAt)
});

const mapCompatibilityStatement = (record: UnknownRecord): CompatibilityStatement => ({
  id: toText(pick(record, ['id', 'statementId'])) || 'unknown-compatibility',
  extensionId: toText(pick(record, ['extensionId', 'ownerRef'])),
  platformVersion: toText(pick(record, ['platformVersion', 'targetRef'])),
  compatibilityStatus: toText(pick(record, ['compatibilityStatus', 'result'])),
  summary: toText(record.summary),
  blockedReasons: toStringArray(record.blockedReasons),
  suggestedActions: toStringArray(record.suggestedActions),
  permissionImpact: toText(record.permissionImpact)
});

const mapPlatformMarketplaceAuditEvent = (record: UnknownRecord): PlatformMarketplaceAuditEvent => ({
  id: toText(pick(record, ['id', 'eventId'])) || 'unknown-event',
  action: toText(record.action),
  actorUserId: toText(record.actorUserId),
  targetType: toText(record.targetType),
  targetRef: toText(record.targetRef),
  outcome: toText(record.outcome),
  occurredAt: toText(record.occurredAt)
});

export const catalogSourceQueryScope = (query: CatalogSourceListQuery) =>
  buildScopeQueryKey([query.sourceType, query.status]);

export const templateQueryScope = (query: TemplateListQuery) =>
  buildScopeQueryKey([query.category, query.status, query.sourceId]);

export const installationQueryScope = (query: InstallationListQuery) =>
  buildScopeQueryKey([query.targetType, query.status]);

export const extensionQueryScope = (query: ExtensionListQuery) =>
  buildScopeQueryKey([query.status, query.compatibilityStatus]);

export const listCatalogSources = async (query: CatalogSourceListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/marketplace/catalog-sources', query));
  return normalizeListResponse(response, mapCatalogSource);
};

export const createCatalogSource = async (payload: CreateCatalogSourcePayload) => {
  const response = await fetchJSON<unknown>('/marketplace/catalog-sources', {
    method: 'POST',
    body: JSON.stringify({
      name: payload.name,
      sourceType: payload.sourceType,
      endpointRef: payload.endpoint,
      visibilityScope: payload.visibleScope
    })
  });
  return mapCatalogSource(isRecord(response) ? response : {});
};

export const syncCatalogSource = async (sourceId: string) => {
  const response = await fetchJSON<unknown>(`/marketplace/catalog-sources/${sourceId}/sync`, {
    method: 'POST'
  });
  return mapCatalogSource(isRecord(response) ? response : { id: sourceId, syncStatus: 'pending' });
};

export const listTemplates = async (query: TemplateListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/marketplace/templates', query));
  return normalizeListResponse(response, mapApplicationTemplate);
};

export const getTemplateDetail = async (templateId: string) => {
  const response = await fetchJSON<unknown>(`/marketplace/templates/${templateId}`);
  if (!isRecord(response)) {
    return mapApplicationTemplate({});
  }
  if (isRecord(response.template)) {
    return mapApplicationTemplate({
      ...response.template,
      versions: Array.isArray(response.versions) ? response.versions : []
    });
  }
  return mapApplicationTemplate(response);
};

export const listTemplateReleases = async (templateId: string) => {
  const response = await fetchJSON<unknown>(`/marketplace/templates/${templateId}/releases`);
  return normalizeListResponse(response, mapTemplateReleaseScope);
};

export const createTemplateRelease = async (
  templateId: string,
  payload: CreateTemplateReleasePayload
) => {
  const response = await fetchJSON<unknown>(`/marketplace/templates/${templateId}/releases`, {
    method: 'POST',
    body: JSON.stringify({
      version: payload.version,
      scopeType: payload.targetType,
      scopeId: toNumber(payload.targetRef) ?? 0,
      targetRef: payload.targetRef,
      visibilityMode: payload.visibilitySummary || 'scope'
    })
  });
  return mapTemplateReleaseScope(isRecord(response) ? response : {});
};

export const listInstallations = async (query: InstallationListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/marketplace/installations', query));
  return normalizeListResponse(response, mapInstallationRecord);
};

export const listExtensions = async (query: ExtensionListQuery = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/marketplace/extensions', query));
  return normalizeListResponse(response, mapExtensionPackage);
};

export const createExtensionPackage = async (payload: CreateExtensionPackagePayload) => {
  const response = await fetchJSON<unknown>('/marketplace/extensions', {
    method: 'POST',
    body: JSON.stringify({
      name: payload.name,
      extensionType: 'plugin',
      version: payload.version,
      visibilityScope: payload.visibilityScope,
      entrySummary: payload.permissionSummary,
      permissionDeclaration: [payload.permissionSummary]
    })
  });
  return mapExtensionPackage(isRecord(response) ? response : {});
};

export const enableExtension = async (extensionId: string) => {
  const response = await fetchJSON<unknown>(`/marketplace/extensions/${extensionId}/enable`, {
    method: 'POST'
  });
  return mapExtensionPackage(isRecord(response) ? response : { id: extensionId, status: 'pending' });
};

export const disableExtension = async (extensionId: string) => {
  const response = await fetchJSON<unknown>(`/marketplace/extensions/${extensionId}/disable`, {
    method: 'POST'
  });
  return mapExtensionPackage(isRecord(response) ? response : { id: extensionId, status: 'disabled' });
};

export const getExtensionCompatibility = async (extensionId: string) => {
  const response = await fetchJSON<unknown>(
    `/marketplace/extensions/${extensionId}/compatibility`
  );
  if (!isRecord(response)) {
    return mapCompatibilityStatement({});
  }
  if (Array.isArray(response.statements)) {
    const first = response.statements.find(isRecord);
    return mapCompatibilityStatement(first ?? {});
  }
  return mapCompatibilityStatement(response);
};

export const listPlatformMarketplaceAuditEvents = async (query: Record<string, QueryValue> = {}) => {
  const response = await fetchJSON<unknown>(withQuery('/audit/platform-marketplace/events', query));
  return normalizeListResponse(response, mapPlatformMarketplaceAuditEvent);
};

export const hasUpgradeGap = (record: InstallationRecord) =>
  Boolean(record.currentVersion && record.latestVersion && record.currentVersion !== record.latestVersion);

export const summarizeTargetScope = (targetType?: string, targetRef?: string) =>
  [targetType, targetRef].filter(Boolean).join(' / ') || '未指定范围';

export const isExtensionBlocked = (compatibility?: CompatibilityStatement) =>
  compatibility?.compatibilityStatus === 'blocked' ||
  compatibility?.compatibilityStatus === 'incompatible' ||
  compatibility?.blockedReasons.length !== 0;

export const isTemplateOffline = (template: ApplicationTemplate) =>
  template.status === 'offline' || template.status === 'retired';

export const hasPermissionDeclaration = (extensionPackage: ExtensionPackage) =>
  toBoolean(extensionPackage.permissionSummary) ?? Boolean(extensionPackage.permissionSummary);
