import { fetchJSON } from '@/services/api/client';
import type {
  CreateExceptionRequestDTO,
  CreatePolicyAssignmentRequestDTO,
  CreateSecurityPolicyRequestDTO,
  Pagination,
  PolicyAssignmentDTO,
  PolicyDistributionTaskDTO,
  PolicyExceptionRequestDTO,
  PolicyExceptionStatus,
  PolicyHitRecordDTO,
  PolicyRemediationStatus,
  SecurityPolicyCategory,
  SecurityPolicyDTO,
  SecurityPolicyRiskLevel,
  SecurityPolicyScopeLevel,
  SecurityPolicyStatus,
  ReviewExceptionRequestDTO,
  SwitchPolicyModeRequestDTO,
  UpdatePolicyRemediationRequestDTO,
  UpdateSecurityPolicyRequestDTO
} from '@/services/api/types';

type QueryValue = string | number | boolean | undefined | null;

export type SecurityPolicyListQuery = {
  workspaceId?: string;
  projectId?: string;
  scopeLevel?: SecurityPolicyScopeLevel;
  status?: SecurityPolicyStatus;
  category?: SecurityPolicyCategory;
};

export type PolicyHitListQuery = {
  policyId?: string;
  workspaceId?: string;
  projectId?: string;
  clusterId?: string;
  namespace?: string;
  riskLevel?: SecurityPolicyRiskLevel;
  remediationStatus?: PolicyRemediationStatus;
  from?: string;
  to?: string;
};

export type PolicyExceptionListQuery = {
  workspaceId?: string;
  projectId?: string;
  status?: PolicyExceptionStatus;
  policyId?: string;
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

const policyPath = (policyId: string) => `/security-policies/${encodeURIComponent(policyId)}`;

export const listSecurityPolicies = async (query: SecurityPolicyListQuery = {}) => {
  return fetchJSON<Pagination<SecurityPolicyDTO>>(withQuery('/security-policies', query));
};

export const getSecurityPolicy = async (policyId: string) => {
  return fetchJSON<SecurityPolicyDTO>(policyPath(policyId));
};

export const createSecurityPolicy = async (payload: CreateSecurityPolicyRequestDTO) => {
  return fetchJSON<SecurityPolicyDTO>('/security-policies', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateSecurityPolicy = async (
  policyId: string,
  payload: UpdateSecurityPolicyRequestDTO
) => {
  return fetchJSON<SecurityPolicyDTO>(policyPath(policyId), {
    method: 'PATCH',
    body: JSON.stringify(payload)
  });
};

export const listPolicyAssignments = async (policyId: string) => {
  return fetchJSON<Pagination<PolicyAssignmentDTO>>(
    `${policyPath(policyId)}/assignments`
  );
};

export const createPolicyAssignment = async (
  policyId: string,
  payload: CreatePolicyAssignmentRequestDTO
) => {
  return fetchJSON<PolicyDistributionTaskDTO>(`${policyPath(policyId)}/assignments`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const switchPolicyMode = async (
  policyId: string,
  payload: SwitchPolicyModeRequestDTO
) => {
  return fetchJSON<PolicyDistributionTaskDTO>(`${policyPath(policyId)}/mode-switch`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const listPolicyHits = async (query: PolicyHitListQuery = {}) => {
  return fetchJSON<Pagination<PolicyHitRecordDTO>>(withQuery('/security-policies/hits', query));
};

export const createExceptionRequest = async (
  hitId: string,
  payload: CreateExceptionRequestDTO
) => {
  return fetchJSON<PolicyExceptionRequestDTO>(
    `/security-policies/hits/${encodeURIComponent(hitId)}/exceptions`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
};

export const listPolicyExceptions = async (query: PolicyExceptionListQuery = {}) => {
  return fetchJSON<Pagination<PolicyExceptionRequestDTO>>(
    withQuery('/security-policies/exceptions', query)
  );
};

export const reviewExceptionRequest = async (
  exceptionId: string,
  payload: ReviewExceptionRequestDTO
) => {
  return fetchJSON<PolicyExceptionRequestDTO>(
    `/security-policies/exceptions/${encodeURIComponent(exceptionId)}/review`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
};

export const updatePolicyHitRemediation = async (
  hitId: string,
  payload: UpdatePolicyRemediationRequestDTO
) => {
  return fetchJSON<PolicyHitRecordDTO>(
    `/security-policies/hits/${encodeURIComponent(hitId)}/remediation`,
    {
      method: 'PATCH',
      body: JSON.stringify(payload)
    }
  );
};
