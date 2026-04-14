export type Pagination<T> = {
  items: T[];
  count?: number;
};

export type OperationStatus =
  | 'pending'
  | 'running'
  | 'partially_succeeded'
  | 'succeeded'
  | 'failed'
  | 'canceled';

export type OperationDTO = {
  id: string | number;
  requestId?: string;
  operatorId?: string | number;
  operationType?: string;
  type?: string;
  targetRef?: string;
  clusterId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  riskLevel?: 'low' | 'medium' | 'high' | 'critical';
  status?: OperationStatus | string;
  reason?: string;
  resultMessage?: string;
  createdAt?: string;
  updatedAt?: string;
  completedAt?: string;
};

export type AuditEventDTO = {
  id: string;
  actorUserId?: string;
  action: string;
  outcome: string;
  clusterId?: string;
  targetRef?: string;
  occurredAt: string;
};

export type AuditExportTaskDTO = {
  taskId: string;
  status: 'pending' | 'running' | 'succeeded' | 'failed';
  resultTotal?: number;
  downloadUrl?: string;
  errorMessage?: string;
  createdAt: string;
  updatedAt?: string;
};

export type ResourceListQueryDTO = {
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  namespace?: string;
  kind?: string;
  keyword?: string;
  health?: string;
  limit?: number;
  offset?: number;
};

export type ResourceInventoryDTO = {
  id?: string | number;
  ID?: string | number;
  cluster?: string | number;
  Cluster?: string | number;
  clusterId?: string | number;
  clusterID?: string | number;
  ClusterID?: string | number;
  clusterName?: string;
  ClusterName?: string;
  namespace?: string;
  Namespace?: string;
  kind?: string;
  Kind?: string;
  resourceType?: string;
  ResourceType?: string;
  name?: string;
  Name?: string;
  status?: string;
  Status?: string;
  health?: string;
  Health?: string;
  labels?: Record<string, unknown>;
  Labels?: Record<string, unknown>;
  updatedAt?: string;
  UpdatedAt?: string;
  createdAt?: string;
  CreatedAt?: string;
};

export type ListResourcesResponseDTO = {
  items?: ResourceInventoryDTO[];
  Items?: ResourceInventoryDTO[];
};

export type ObservabilityScopeDTO = {
  clusterIds?: string[];
  workspaceIds?: string[];
  projectIds?: string[];
  namespaces?: string[];
  resourceKinds?: string[];
  resourceNames?: string[];
};

export type ObservabilityOverviewCardDTO = {
  title: string;
  value: number;
  unit?: string;
  trend?: string;
  severity?: string;
};

export type ObservabilityOverviewDTO = {
  cards: ObservabilityOverviewCardDTO[];
};

export type ObservabilityLogEntryDTO = {
  timestamp: string;
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  namespace?: string;
  workload?: string;
  pod?: string;
  container?: string;
  message: string;
};

export type ObservabilityEventDTO = {
  clusterId?: string;
  namespace?: string;
  involvedKind?: string;
  involvedName?: string;
  eventType?: 'normal' | 'warning';
  reason?: string;
  message?: string;
  firstSeenAt?: string;
  lastSeenAt?: string;
  count?: number;
};

export type ObservabilityMetricPointDTO = {
  timestamp: string;
  value: number;
};

export type ObservabilityMetricSeriesDTO = {
  metricKey: string;
  subjectType: 'cluster' | 'node' | 'namespace' | 'workload' | 'pod';
  subjectRef: string;
  points: ObservabilityMetricPointDTO[];
  dataFreshness?: string;
};

export type ObservabilityAlertDTO = {
  id: string | number;
  ruleId?: string;
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  severity?: 'info' | 'warning' | 'critical';
  status?: 'firing' | 'acknowledged' | 'silenced' | 'resolved';
  summary?: string;
  startsAt?: string;
  acknowledgedAt?: string;
  resolvedAt?: string;
  sourceIncidentKey?: string;
};

export type ObservabilityHandlingRecordDTO = {
  id: string | number;
  incidentId: string | number;
  actionType: string;
  content?: string;
  actedBy?: string | number;
  actedAt?: string;
};

export type ObservabilityAlertRuleDTO = {
  id: string | number;
  name: string;
  description?: string;
  severity: 'info' | 'warning' | 'critical';
  scopeSnapshotJson?: string;
  conditionExpression: string;
  evaluationWindow?: string;
  notificationStrategy?: string;
  status: 'enabled' | 'disabled';
};

export type ObservabilityNotificationTargetDTO = {
  id: string | number;
  name: string;
  targetType: string;
  configRef?: string;
  scopeSnapshot?: string;
  status?: 'active' | 'disabled' | string;
};

export type ObservabilitySilenceWindowDTO = {
  id: string | number;
  name: string;
  scopeSnapshot?: string;
  reason?: string;
  startsAt?: string;
  endsAt?: string;
  status?: 'scheduled' | 'active' | 'expired' | 'canceled' | string;
};

export type WorkloadOperationsViewDTO = {
  clusterId: number;
  workspaceId?: number;
  projectId?: number;
  namespace: string;
  resourceKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  resourceName: string;
  healthStatus: string;
  rolloutStatus: string;
  latestChangeSummary?: string;
  latestActionSummary?: string;
  availableActions?: string[];
};

export type WorkloadInstanceDTO = {
  podName: string;
  containerName?: string;
  nodeName?: string;
  phase: string;
  ready: boolean;
  restartCount?: number;
  startedAt?: string;
  lastTransitionAt?: string;
  logAvailable?: boolean;
  terminalAvailable?: boolean;
};

export type ReleaseRevisionDTO = {
  revision: number;
  sourceKind: 'replicaset' | 'controllerrevision';
  sourceName: string;
  changeCause?: string;
  createdAt?: string;
  isCurrent: boolean;
  rollbackAvailable: boolean;
  summary?: string;
};

export type SubmitWorkloadActionRequestDTO = {
  clusterId: number;
  workspaceId?: number;
  projectId?: number;
  namespace: string;
  resourceKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  resourceName: string;
  targetInstanceRef?: string;
  actionType: 'scale' | 'restart' | 'redeploy' | 'replace-instance' | 'rollback';
  riskConfirmed?: boolean;
  payload?: Record<string, unknown>;
};

export type WorkloadActionDTO = {
  id: number;
  actionType: string;
  status: 'pending' | 'running' | 'succeeded' | 'failed' | 'canceled';
  riskLevel: 'low' | 'medium' | 'high';
  progressMessage?: string;
  resultMessage?: string;
  failureReason?: string;
  startedAt?: string;
  completedAt?: string;
};

export type SubmitBatchOperationRequestDTO = {
  actionType: 'scale' | 'restart' | 'redeploy' | 'replace-instance';
  riskConfirmed?: boolean;
  payload?: Record<string, unknown>;
  targets: Array<{
    clusterId: number;
    workspaceId?: number;
    projectId?: number;
    namespace: string;
    resourceKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
    resourceName: string;
  }>;
};

export type BatchOperationTaskDTO = {
  id: number;
  actionType: string;
  status: 'pending' | 'running' | 'partially_succeeded' | 'succeeded' | 'failed' | 'canceled';
  totalTargets: number;
  succeededTargets?: number;
  failedTargets?: number;
  canceledTargets?: number;
  progressPercent?: number;
  items?: Array<{
    resourceRef?: string;
    status?: string;
    resultMessage?: string;
    failureReason?: string;
  }>;
  startedAt?: string;
  completedAt?: string;
};

export type TerminalSessionDTO = {
  id: number;
  status: 'pending' | 'active' | 'closed' | 'expired' | 'denied' | 'failed';
  podName: string;
  containerName: string;
  workloadKind?: string;
  workloadName?: string;
  streamUrl?: string;
  streamToken?: string;
  startedAt?: string;
  endedAt?: string;
  durationSeconds?: number;
  closeReason?: string;
};

export type GitOpsSourceDTO = {
  id: string;
  name: string;
  sourceType?: 'git' | 'helm' | 'oci' | string;
  status?: 'enabled' | 'disabled' | 'degraded' | string;
  workspaceId?: string;
  projectId?: string;
  lastValidatedAt?: string;
};

export type GitOpsDeliveryUnitDTO = {
  id: string;
  name: string;
  sourceId?: string;
  workspaceId?: string;
  projectId?: string;
  desiredState?: string;
  actualState?: string;
  driftStatus?: 'in_sync' | 'drifted' | 'unknown' | string;
  lastSyncAt?: string;
  lastSyncResult?: string;
};

export type GitOpsOperationType =
  | 'install'
  | 'sync'
  | 'resync'
  | 'upgrade'
  | 'pause'
  | 'resume'
  | 'promote'
  | 'rollback'
  | 'uninstall';

export type GitOpsOperationDTO = {
  id: string | number;
  unitId?: string;
  actionType?: GitOpsOperationType | string;
  operationType: GitOpsOperationType | string;
  status?: OperationStatus | string;
  progressPercent?: number;
  resultSummary?: string;
  startedAt?: string;
  completedAt?: string;
  resultMessage?: string;
  failureReason?: string;
  stages?: Array<{
    environment?: string;
    status?: string;
    targetCount?: number;
    succeededCount?: number;
    failedCount?: number;
    failureReason?: string;
  }>;
};

export type GitOpsActionRequestDTO = {
  actionType: GitOpsOperationType;
  environment?: string;
  targetReleaseId?: number;
  targetAppVersion?: string;
  targetConfigVersion?: string;
  reason?: string;
  stageId?: string;
  targetRevision?: string;
  overrideValues?: Record<string, unknown>;
};

export type GitOpsListQueryDTO = {
  keyword?: string;
  workspaceId?: string;
  projectId?: string;
  limit?: number;
  offset?: number;
};
