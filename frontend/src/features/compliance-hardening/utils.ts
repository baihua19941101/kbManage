import type {
  ArchiveExportStatus,
  ComplianceExceptionStatus,
  ComplianceRecordStatus,
  CoverageStatus,
  FindingRemediationStatus,
  FindingResult,
  RemediationTaskStatus,
  RecheckStatus,
  RiskLevel,
  ScanExecutionStatus
} from '@/services/compliance';

export const formatDateTime = (value?: string) => {
  if (!value) {
    return '—';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString('zh-CN', { hour12: false });
};

export const toIsoDateTime = (value?: string) => {
  if (!value) {
    return undefined;
  }
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
};

export const toDatetimeLocal = (value?: string) => {
  if (!value) {
    return undefined;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return undefined;
  }
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60_000);
  return local.toISOString().slice(0, 16);
};

export const parseKeyValueLines = (value?: string): Record<string, string> | undefined => {
  if (!value) {
    return undefined;
  }

  const entries = value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean)
    .map((item) => {
      const [key, ...rest] = item.split('=');
      return [key?.trim(), rest.join('=').trim()] as const;
    })
    .filter((item): item is readonly [string, string] => Boolean(item[0] && item[1]));

  return entries.length > 0 ? Object.fromEntries(entries) : undefined;
};

export const stringifyKeyValueLines = (value?: Record<string, string>) => {
  if (!value) {
    return '';
  }
  return Object.entries(value)
    .map(([key, itemValue]) => `${key}=${itemValue}`)
    .join('\n');
};

export const riskColorMap: Record<RiskLevel, string> = {
  low: 'default',
  medium: 'gold',
  high: 'orange',
  critical: 'red'
};

export const findingResultColorMap: Record<FindingResult, string> = {
  pass: 'green',
  fail: 'red',
  warn: 'gold',
  skipped: 'default',
  error: 'magenta'
};

export const recordStatusColorMap: Record<ComplianceRecordStatus, string> = {
  draft: 'default',
  active: 'green',
  disabled: 'orange',
  archived: 'purple'
};

export const scanStatusColorMap: Record<ScanExecutionStatus, string> = {
  pending: 'default',
  running: 'processing',
  partially_succeeded: 'gold',
  succeeded: 'green',
  failed: 'red',
  canceled: 'default'
};

export const coverageStatusColorMap: Record<CoverageStatus, string> = {
  full: 'green',
  partial: 'gold',
  unavailable: 'red'
};

export const remediationStatusColorMap: Record<FindingRemediationStatus, string> = {
  open: 'red',
  in_progress: 'processing',
  exception_active: 'purple',
  ready_for_recheck: 'gold',
  closed: 'green'
};

export const remediationTaskStatusColorMap: Record<RemediationTaskStatus, string> = {
  todo: 'default',
  in_progress: 'processing',
  blocked: 'red',
  done: 'green',
  canceled: 'default'
};

export const exceptionStatusColorMap: Record<ComplianceExceptionStatus, string> = {
  pending: 'gold',
  approved: 'blue',
  rejected: 'red',
  active: 'green',
  expired: 'default',
  revoked: 'purple'
};

export const recheckStatusColorMap: Record<RecheckStatus, string> = {
  pending: 'default',
  running: 'processing',
  passed: 'green',
  failed: 'red',
  canceled: 'default'
};

export const archiveStatusColorMap: Record<ArchiveExportStatus, string> = {
  pending: 'default',
  running: 'processing',
  succeeded: 'green',
  failed: 'red',
  expired: 'orange'
};
