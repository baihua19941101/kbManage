export type Pagination<T> = {
  items: T[];
  count?: number;
};

export type OperationStatus = 'pending' | 'running' | 'succeeded' | 'failed';

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
