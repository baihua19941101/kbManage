import { getAccessToken } from '@/features/auth/store';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const pickMessage = (payload: unknown): string | undefined => {
  if (typeof payload === 'string') {
    const trimmed = payload.trim();
    return trimmed.length > 0 ? trimmed : undefined;
  }

  if (!isRecord(payload)) {
    return undefined;
  }

  const candidates = ['message', 'error', 'detail'];
  for (const key of candidates) {
    const value = payload[key];
    if (typeof value === 'string' && value.trim().length > 0) {
      return value.trim();
    }
  }

  const nestedError = payload.error;
  if (isRecord(nestedError) && typeof nestedError.message === 'string') {
    const trimmed = nestedError.message.trim();
    if (trimmed.length > 0) {
      return trimmed;
    }
  }

  return undefined;
};

const pickCode = (payload: unknown): string | undefined => {
  if (!isRecord(payload)) {
    return undefined;
  }

  const value = payload.code;
  return typeof value === 'string' && value.trim().length > 0 ? value.trim() : undefined;
};

const resolvePath = (path: string): string => {
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path;
  }

  if (path.startsWith('/')) {
    return `${API_BASE_URL}${path}`;
  }

  return `${API_BASE_URL}/${path}`;
};

const parsePayload = async (response: Response): Promise<unknown> => {
  if (response.status === 204) {
    return undefined;
  }

  const text = await response.text();
  if (!text) {
    return undefined;
  }

  const contentType = response.headers.get('content-type') || '';
  if (!contentType.includes('application/json')) {
    return text;
  }

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return text;
  }
};

const buildHeaders = (
  inputHeaders: HeadersInit | undefined,
  body: BodyInit | null | undefined,
  skipAuth: boolean
): Headers => {
  const headers = new Headers(inputHeaders);

  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json');
  }

  const isFormData = typeof FormData !== 'undefined' && body instanceof FormData;
  if (body && !isFormData && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const token = getAccessToken();
  if (!skipAuth && token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  return headers;
};

export class ApiError extends Error {
  status: number;

  code?: string;

  details?: unknown;

  url: string;

  constructor(
    status: number,
    message: string,
    options: {
      url: string;
      code?: string;
      details?: unknown;
    }
  ) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = options.code;
    this.details = options.details;
    this.url = options.url;
  }
}

export const normalizeApiError = (
  error: unknown,
  fallback = '请求失败，请稍后重试'
): string => {
  if (error instanceof ApiError && error.message.trim().length > 0) {
    return error.message.trim();
  }

  if (error instanceof Error && error.message.trim().length > 0) {
    return error.message.trim();
  }

  if (typeof error === 'string' && error.trim().length > 0) {
    return error.trim();
  }

  return fallback;
};

export type FetchJSONInit = RequestInit & {
  skipAuth?: boolean;
};

export const fetchJSON = async <T>(
  path: string,
  init: FetchJSONInit = {}
): Promise<T> => {
  const { skipAuth = false, headers: inputHeaders, ...requestInit } = init;
  const url = resolvePath(path);

  const response = await fetch(url, {
    ...requestInit,
    headers: buildHeaders(inputHeaders, requestInit.body, skipAuth)
  });

  const payload = await parsePayload(response);

  if (!response.ok) {
    throw new ApiError(
      response.status,
      pickMessage(payload) || `Request failed with status ${response.status}`,
      {
        url,
        code: pickCode(payload),
        details: payload
      }
    );
  }

  if (payload === undefined || typeof payload === 'string') {
    return {} as T;
  }

  return payload as T;
};

export const buildScopeQueryKey = (parts: Array<string | number | undefined | null>) => {
  return parts
    .filter((part): part is string | number => part !== undefined && part !== null && part !== '')
    .map((part) => String(part))
    .join(':');
};

export const isApiErrorStatus = (
  error: unknown,
  statuses: readonly number[]
): error is ApiError => {
  return error instanceof ApiError && statuses.includes(error.status);
};

export const isAuthorizationError = (error: unknown): error is ApiError => {
  return isApiErrorStatus(error, [401, 403]);
};

export const normalizeAuthorizationError = (
  error: unknown,
  fallback = '当前操作未授权，请联系管理员调整权限。'
): string => {
  if (isAuthorizationError(error)) {
    return normalizeApiError(error, fallback);
  }

  return fallback;
};
