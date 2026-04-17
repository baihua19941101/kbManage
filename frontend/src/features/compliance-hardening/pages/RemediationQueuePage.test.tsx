import { canManageComplianceRemediation, canReadCompliance } from '@/features/auth/store';

describe('RemediationQueuePage', () => {
  it('uses compliance remediation permission gate', () => {
    expect(canReadCompliance({ id: 'u1', username: 'admin', platformRoles: ['platform-admin'] })).toBe(true);
    expect(canManageComplianceRemediation({ id: 'u1', username: 'admin', platformRoles: ['platform-admin'] })).toBe(true);
  });
});
