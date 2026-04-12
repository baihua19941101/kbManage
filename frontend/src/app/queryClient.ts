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
