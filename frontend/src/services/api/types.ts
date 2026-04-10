export type Pagination<T> = {
  items: T[];
  count?: number;
};

export type OperationStatus = 'pending' | 'running' | 'succeeded' | 'failed';

export type OperationDTO = {
  id: string;
  type: string;
  clusterId: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  riskLevel?: 'low' | 'medium' | 'high' | 'critical';
  status: OperationStatus;
  reason?: string;
  createdAt: string;
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
