import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { ApiError, normalizeApiError } from '@/services/api/client';

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const suppressGlobalError = (meta: unknown): boolean => {
  if (!isRecord(meta)) {
    return false;
  }

  return meta.suppressGlobalError === true;
};

export const normalizeErrorMessage = (
  error: unknown,
  fallback = '请求失败，请稍后重试'
): string => {
  if (error instanceof ApiError) {
    return error.message;
  }

  if (error instanceof Error && error.message.trim().length > 0) {
    return error.message.trim();
  }

  if (typeof error === 'string' && error.trim().length > 0) {
    return error.trim();
  }

  if (isRecord(error) && typeof error.message === 'string' && error.message.trim().length > 0) {
    return error.message.trim();
  }

  return fallback;
};

export const queryKeys = {
  observability: {
    all: ['observability'] as const,
    overview: (scope?: string) => ['observability', 'overview', scope ?? 'default'] as const,
    logs: (scope?: string) => ['observability', 'logs', scope ?? 'default'] as const,
    events: (scope?: string) => ['observability', 'events', scope ?? 'default'] as const,
    metrics: (scope?: string) => ['observability', 'metrics', scope ?? 'default'] as const,
    alerts: (scope?: string) => ['observability', 'alerts', scope ?? 'default'] as const
  },
  workloadOps: {
    all: ['workloadOps'] as const,
    context: (scope?: string) => ['workloadOps', 'context', scope ?? 'default'] as const,
    instances: (scope?: string) => ['workloadOps', 'instances', scope ?? 'default'] as const,
    revisions: (scope?: string) => ['workloadOps', 'revisions', scope ?? 'default'] as const,
    action: (id?: number | string) => ['workloadOps', 'action', id ?? 'unknown'] as const,
    batch: (id?: number | string) => ['workloadOps', 'batch', id ?? 'unknown'] as const,
    terminal: (id?: number | string) => ['workloadOps', 'terminal', id ?? 'unknown'] as const
  },
  gitops: {
    all: ['gitops'] as const,
    sources: (scope?: string) => ['gitops', 'sources', scope ?? 'default'] as const,
    deliveryUnits: (scope?: string) => ['gitops', 'deliveryUnits', scope ?? 'default'] as const,
    operation: (id?: number | string) => ['gitops', 'operation', id ?? 'unknown'] as const
  },
  securityPolicy: {
    all: ['securityPolicy'] as const,
    list: (scope?: string) => ['securityPolicy', 'list', scope ?? 'default'] as const,
    detail: (policyId?: number | string) =>
      ['securityPolicy', 'detail', policyId ?? 'unknown'] as const,
    assignments: (policyId?: number | string) =>
      ['securityPolicy', 'assignments', policyId ?? 'unknown'] as const
  },
  compliance: {
    all: ['compliance'] as const,
    baselines: (scope?: string) => ['compliance', 'baselines', scope ?? 'default'] as const,
    profiles: (scope?: string) => ['compliance', 'profiles', scope ?? 'default'] as const,
    scans: (scope?: string) => ['compliance', 'scans', scope ?? 'default'] as const,
    findings: (scope?: string) => ['compliance', 'findings', scope ?? 'default'] as const,
    remediation: (scope?: string) => ['compliance', 'remediation', scope ?? 'default'] as const,
    exceptions: (scope?: string) => ['compliance', 'exceptions', scope ?? 'default'] as const,
    rechecks: (scope?: string) => ['compliance', 'rechecks', scope ?? 'default'] as const,
    overview: (scope?: string) => ['compliance', 'overview', scope ?? 'default'] as const,
    trends: (scope?: string) => ['compliance', 'trends', scope ?? 'default'] as const,
    exports: (scope?: string) => ['compliance', 'exports', scope ?? 'default'] as const,
    audit: (scope?: string) => ['compliance', 'audit', scope ?? 'default'] as const
  },
  backupRestore: {
    all: ['backupRestore'] as const,
    policies: (scope?: string) => ['backupRestore', 'policies', scope ?? 'default'] as const,
    restorePoints: (scope?: string) =>
      ['backupRestore', 'restorePoints', scope ?? 'default'] as const,
    restoreJobs: (scope?: string) =>
      ['backupRestore', 'restoreJobs', scope ?? 'default'] as const,
    drillPlans: (scope?: string) =>
      ['backupRestore', 'drillPlans', scope ?? 'default'] as const,
    drillRecord: (recordId?: number | string) =>
      ['backupRestore', 'drillRecord', recordId ?? 'unknown'] as const,
    audit: (scope?: string) => ['backupRestore', 'audit', scope ?? 'default'] as const
  },
  identityTenancy: {
    all: ['identityTenancy'] as const,
    sources: (scope?: string) => ['identityTenancy', 'sources', scope ?? 'default'] as const,
    organizations: (scope?: string) =>
      ['identityTenancy', 'organizations', scope ?? 'default'] as const,
    roles: (scope?: string) => ['identityTenancy', 'roles', scope ?? 'default'] as const,
    assignments: (scope?: string) =>
      ['identityTenancy', 'assignments', scope ?? 'default'] as const,
    sessions: (scope?: string) => ['identityTenancy', 'sessions', scope ?? 'default'] as const,
    risks: (scope?: string) => ['identityTenancy', 'risks', scope ?? 'default'] as const,
    audit: (scope?: string) => ['identityTenancy', 'audit', scope ?? 'default'] as const
  },
  platformMarketplace: {
    all: ['platformMarketplace'] as const,
    catalogSources: (scope?: string) =>
      ['platformMarketplace', 'catalogSources', scope ?? 'default'] as const,
    templates: (scope?: string) => ['platformMarketplace', 'templates', scope ?? 'default'] as const,
    installations: (scope?: string) =>
      ['platformMarketplace', 'installations', scope ?? 'default'] as const,
    extensions: (scope?: string) =>
      ['platformMarketplace', 'extensions', scope ?? 'default'] as const,
    compatibility: (scope?: string) =>
      ['platformMarketplace', 'compatibility', scope ?? 'default'] as const,
    audit: (scope?: string) => ['platformMarketplace', 'audit', scope ?? 'default'] as const
  },
  sreScale: {
    all: ['sreScale'] as const,
    haPolicies: (scope?: string) => ['sreScale', 'haPolicies', scope ?? 'default'] as const,
    health: (scope?: string) => ['sreScale', 'health', scope ?? 'default'] as const,
    maintenanceWindows: (scope?: string) =>
      ['sreScale', 'maintenanceWindows', scope ?? 'default'] as const,
    upgrades: (scope?: string) => ['sreScale', 'upgrades', scope ?? 'default'] as const,
    capacity: (scope?: string) => ['sreScale', 'capacity', scope ?? 'default'] as const,
    scaleEvidence: (scope?: string) => ['sreScale', 'scaleEvidence', scope ?? 'default'] as const,
    runbooks: (scope?: string) => ['sreScale', 'runbooks', scope ?? 'default'] as const,
    audit: (scope?: string) => ['sreScale', 'audit', scope ?? 'default'] as const
  },
  enterprisePolish: {
    all: ['enterprisePolish'] as const,
    permissionTrails: (scope?: string) =>
      ['enterprisePolish', 'permissionTrails', scope ?? 'default'] as const,
    keyOperations: (scope?: string) =>
      ['enterprisePolish', 'keyOperations', scope ?? 'default'] as const,
    coverage: (scope?: string) => ['enterprisePolish', 'coverage', scope ?? 'default'] as const,
    actionItems: (scope?: string) => ['enterprisePolish', 'actionItems', scope ?? 'default'] as const,
    reports: (scope?: string) => ['enterprisePolish', 'reports', scope ?? 'default'] as const,
    deliveryArtifacts: (scope?: string) =>
      ['enterprisePolish', 'deliveryArtifacts', scope ?? 'default'] as const,
    deliveryBundles: (scope?: string) =>
      ['enterprisePolish', 'deliveryBundles', scope ?? 'default'] as const,
    deliveryChecklist: (bundleId?: string) =>
      ['enterprisePolish', 'deliveryChecklist', bundleId ?? 'unknown'] as const,
    audit: (scope?: string) => ['enterprisePolish', 'audit', scope ?? 'default'] as const
  }
};

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      if (suppressGlobalError(query.meta)) {
        return;
      }

      message.error(normalizeApiError(error, '查询失败，请稍后重试'));
    }
  }),
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      if (mutation.options.onError || suppressGlobalError(mutation.meta)) {
        return;
      }

      message.error(normalizeApiError(error, '提交失败，请稍后重试'));
    }
  }),
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
      refetchOnWindowFocus: false
    },
    mutations: {
      retry: 0
    }
  }
});
