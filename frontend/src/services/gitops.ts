import { fetchJSON } from '@/services/api/client';
import type {
  GitOpsActionRequestDTO,
  GitOpsDeliveryUnitDTO,
  GitOpsListQueryDTO,
  GitOpsOperationDTO,
  GitOpsSourceDTO,
  Pagination
} from '@/services/api/types';

type QueryValue = string | number | boolean | undefined | null;
export type ResourceId = string | number;

export type GitOpsSourceType = 'git' | 'package' | string;
export type GitOpsSourceStatus = 'pending' | 'ready' | 'failed' | 'disabled' | string;
export type GitOpsTargetGroupStatus = 'active' | 'stale' | 'disabled' | string;
export type GitOpsActionType = GitOpsActionRequestDTO['actionType'];

export type GitOpsSourceItem = GitOpsSourceDTO & {
  id: ResourceId;
  sourceType?: GitOpsSourceType;
  endpoint?: string;
  defaultRef?: string;
  credentialRef?: string;
  workspaceId?: string | number;
  projectId?: string | number;
  status?: GitOpsSourceStatus;
  lastVerifiedAt?: string;
  lastErrorMessage?: string;
};

export type GitOpsSourceFormData = {
  name: string;
  sourceType: Exclude<GitOpsSourceType, string> | GitOpsSourceType;
  endpoint: string;
  defaultRef?: string;
  credentialRef?: string;
  workspaceId: number;
  projectId?: number;
};

export type GitOpsTargetGroupItem = {
  id: ResourceId;
  name: string;
  workspaceId: number;
  projectId?: number;
  clusterRefs?: number[];
  selectorSummary?: string;
  status?: GitOpsTargetGroupStatus;
};

export type GitOpsTargetGroupFormData = {
  name: string;
  workspaceId: number;
  projectId?: number;
  clusterRefs?: number[];
  selectorSummary?: string;
  description?: string;
};

export type GitOpsEnvironmentStage = {
  name: string;
  orderIndex: number;
  targetGroupId: number;
  promotionMode: 'manual' | 'automatic' | string;
  paused?: boolean;
};

export type GitOpsConfigurationOverlay = {
  overlayType: 'values' | 'patch' | 'manifest-snippet' | string;
  overlayRef: string;
  precedence?: number;
  effectiveScope?: string;
};

export type GitOpsDeliveryUnitItem = GitOpsDeliveryUnitDTO & {
  id: ResourceId;
  sourceId?: ResourceId;
  deliveryStatus?: string;
  desiredRevision?: string;
  desiredAppVersion?: string;
  desiredConfigVersion?: string;
  lastSyncedAt?: string;
};

export type GitOpsDeliveryUnitDetail = GitOpsDeliveryUnitItem & {
  sourcePath?: string;
  defaultNamespace?: string;
  syncMode?: 'manual' | 'auto' | string;
  environments?: GitOpsEnvironmentStage[];
  overlays?: GitOpsConfigurationOverlay[];
};

export type GitOpsDeliveryUnitFormData = {
  name: string;
  workspaceId: number;
  projectId?: number;
  sourceId: number;
  sourcePath?: string;
  defaultNamespace?: string;
  syncMode: 'manual' | 'auto';
  desiredRevision?: string;
  desiredAppVersion?: string;
  desiredConfigVersion?: string;
  environments: GitOpsEnvironmentStage[];
  overlays?: GitOpsConfigurationOverlay[];
};

export type GitOpsDeliveryUnitStatus = {
  deliveryStatus?: string;
  driftStatus?: string;
  lastSyncedAt?: string;
  lastSyncResult?: string;
  lastErrorMessage?: string;
  environments?: Array<{
    environment?: string;
    syncStatus?: string;
    driftStatus?: string;
    targetCount?: number;
    succeededCount?: number;
    failedCount?: number;
  }>;
};

export type GitOpsDeliveryDiff = {
  summary?: {
    added?: number;
    modified?: number;
    removed?: number;
    unavailable?: number;
  };
  items?: Array<{
    objectRef?: string;
    environment?: string;
    diffType?: 'added' | 'modified' | 'removed' | 'unavailable' | string;
    desiredSummary?: string;
    liveSummary?: string;
  }>;
};

export type GitOpsReleaseRevision = {
  id: ResourceId;
  sourceRevision?: string;
  appVersion?: string;
  configVersion?: string;
  status?: 'active' | 'historical' | 'failed' | 'rolled_back' | string;
  rollbackAvailable?: boolean;
  createdBy?: number | string;
  createdAt?: string;
  releaseNotesSummary?: string;
};

export type GitOpsDeliveryDiffQuery = {
  environment?: string;
  targetGroupId?: number;
};

export type GitOpsReleaseRevisionQuery = {
  environment?: string;
};

export type GitOpsSubmitActionInput = GitOpsActionRequestDTO;

export type GitOpsSourceListQuery = GitOpsListQueryDTO & {
  sourceType?: GitOpsSourceType;
  status?: GitOpsSourceStatus;
};

export type GitOpsDeliveryUnitListQuery = GitOpsListQueryDTO & {
  environment?: string;
  status?: string;
  driftStatus?: string;
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
  return queryString ? path + '?' + queryString : path;
};

const gitOpsSourcePath = (sourceId: ResourceId) =>
  `/gitops/sources/${encodeURIComponent(String(sourceId))}`;

const gitOpsTargetGroupPath = (targetGroupId: ResourceId) =>
  `/gitops/target-groups/${encodeURIComponent(String(targetGroupId))}`;

const gitOpsUnitPath = (unitId: ResourceId) =>
  `/gitops/delivery-units/${encodeURIComponent(String(unitId))}`;

const normalizeGitOpsOperation = (operation: GitOpsOperationDTO): GitOpsOperationDTO => {
  return {
    ...operation,
    operationType: operation.operationType || operation.actionType || 'sync',
    resultMessage: operation.resultMessage || operation.resultSummary
  };
};

export const listGitOpsSources = async (query: GitOpsListQueryDTO = {}) => {
  return fetchJSON<Pagination<GitOpsSourceItem>>(withQuery('/gitops/sources', query));
};

export const createGitOpsSource = async (payload: GitOpsSourceFormData) => {
  return fetchJSON<GitOpsSourceItem>('/gitops/sources', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateGitOpsSource = async (
  sourceId: ResourceId,
  payload: Partial<GitOpsSourceFormData> & { disabled?: boolean }
) => {
  return fetchJSON<GitOpsSourceItem>(gitOpsSourcePath(sourceId), {
    method: 'PATCH',
    body: JSON.stringify(payload)
  });
};

export const verifyGitOpsSource = async (sourceId: ResourceId) => {
  return fetchJSON<GitOpsOperationDTO>(`${gitOpsSourcePath(sourceId)}/verify`, {
    method: 'POST'
  });
};

export const listGitOpsTargetGroups = async (query: GitOpsListQueryDTO = {}) => {
  return fetchJSON<Pagination<GitOpsTargetGroupItem>>(withQuery('/gitops/target-groups', query));
};

export const createGitOpsTargetGroup = async (payload: GitOpsTargetGroupFormData) => {
  return fetchJSON<GitOpsTargetGroupItem>('/gitops/target-groups', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateGitOpsTargetGroup = async (
  targetGroupId: ResourceId,
  payload: Partial<GitOpsTargetGroupFormData> & { disabled?: boolean }
) => {
  return fetchJSON<GitOpsTargetGroupItem>(gitOpsTargetGroupPath(targetGroupId), {
    method: 'PATCH',
    body: JSON.stringify(payload)
  });
};

export const listGitOpsDeliveryUnits = async (query: GitOpsDeliveryUnitListQuery = {}) => {
  return fetchJSON<Pagination<GitOpsDeliveryUnitItem>>(withQuery('/gitops/delivery-units', query));
};

export const createGitOpsDeliveryUnit = async (payload: GitOpsDeliveryUnitFormData) => {
  return fetchJSON<GitOpsDeliveryUnitDetail>('/gitops/delivery-units', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateGitOpsDeliveryUnit = async (
  unitId: ResourceId,
  payload: Partial<GitOpsDeliveryUnitFormData>
) => {
  return fetchJSON<GitOpsDeliveryUnitDetail>(gitOpsUnitPath(unitId), {
    method: 'PATCH',
    body: JSON.stringify(payload)
  });
};

export const getGitOpsDeliveryUnit = async (unitId: ResourceId) => {
  return fetchJSON<GitOpsDeliveryUnitDetail>(gitOpsUnitPath(unitId));
};

export const getGitOpsDeliveryUnitStatus = async (
  unitId: ResourceId,
  query: { environment?: string } = {}
) => {
  return fetchJSON<GitOpsDeliveryUnitStatus>(
    withQuery(`${gitOpsUnitPath(unitId)}/status`, query)
  );
};

export const getGitOpsDeliveryUnitDiff = async (
  unitId: ResourceId,
  query: GitOpsDeliveryDiffQuery = {}
) => {
  return fetchJSON<GitOpsDeliveryDiff>(
    withQuery(`${gitOpsUnitPath(unitId)}/diff`, query)
  );
};

export const listGitOpsReleaseRevisions = async (
  unitId: ResourceId,
  query: GitOpsReleaseRevisionQuery = {}
) => {
  return fetchJSON<{ items: GitOpsReleaseRevision[] }>(
    withQuery(`${gitOpsUnitPath(unitId)}/releases`, query)
  );
};

export const submitGitOpsAction = async (
  unitId: ResourceId,
  payload: GitOpsSubmitActionInput
): Promise<GitOpsOperationDTO> => {
  const operation = await fetchJSON<GitOpsOperationDTO>(`${gitOpsUnitPath(unitId)}/actions`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });

  return normalizeGitOpsOperation(operation);
};

export const submitGitOpsRollback = async (
  unitId: ResourceId,
  payload: {
    targetReleaseId: number;
    environment?: string;
    reason?: string;
  }
) => {
  return submitGitOpsAction(unitId, {
    actionType: 'rollback',
    targetReleaseId: payload.targetReleaseId,
    environment: payload.environment,
    reason: payload.reason
  });
};

export const getGitOpsOperation = async (operationId: ResourceId) => {
  const operation = await fetchJSON<GitOpsOperationDTO>(
    `/gitops/operations/${encodeURIComponent(String(operationId))}`
  );

  return normalizeGitOpsOperation(operation);
};
