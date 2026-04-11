import { fetchJSON } from '@/services/api/client';

export type RoleBindingItem = {
  id: string;
  subjectType: 'user' | 'group';
  subjectId: string;
  scopeType: 'workspace' | 'project';
  scopeId: string;
  scopeRoleId?: string;
  roleKey: string;
  grantedBy?: string;
  createdAt?: string;
};

export type CreateRoleBindingPayload = {
  subjectType: 'user' | 'group';
  subjectId: string;
  scopeType: 'workspace' | 'project';
  scopeId: string;
  roleKey: string;
};

const toText = (value: unknown): string => {
  if (typeof value === 'string') {
    return value;
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(value);
  }
  return '';
};

const mapItem = (item: Record<string, unknown>): RoleBindingItem => ({
  id: toText(item.id),
  subjectType: (toText(item.subjectType) as RoleBindingItem['subjectType']) || 'user',
  subjectId: toText(item.subjectId),
  scopeType: (toText(item.scopeType) as RoleBindingItem['scopeType']) || 'workspace',
  scopeId: toText(item.scopeId),
  scopeRoleId: toText(item.scopeRoleId) || undefined,
  roleKey: toText(item.roleKey),
  grantedBy: toText(item.grantedBy) || undefined,
  createdAt: toText(item.createdAt) || undefined
});

export const createRoleBinding = async (
  payload: CreateRoleBindingPayload
): Promise<RoleBindingItem> => {
  const response = await fetchJSON<Record<string, unknown>>('/role-bindings', {
    method: 'POST',
    body: JSON.stringify({
      subjectType: payload.subjectType,
      subjectId: payload.subjectId,
      scopeType: payload.scopeType,
      scopeId: payload.scopeId,
      roleKey: payload.roleKey
    })
  });
  return mapItem(response);
};

export const listRoleBindings = async (query?: {
  subjectType?: 'user' | 'group';
  subjectId?: string;
  scopeType?: 'workspace' | 'project';
  scopeId?: string;
}): Promise<RoleBindingItem[]> => {
  const search = new URLSearchParams();
  if (query?.subjectType) search.set('subjectType', query.subjectType);
  if (query?.subjectId) search.set('subjectId', query.subjectId);
  if (query?.scopeType) search.set('scopeType', query.scopeType);
  if (query?.scopeId) search.set('scopeId', query.scopeId);

  const qs = search.toString();
  const response = await fetchJSON<{ items?: Record<string, unknown>[] }>(
    `/role-bindings${qs ? `?${qs}` : ''}`
  );
  return Array.isArray(response.items) ? response.items.map(mapItem) : [];
};
