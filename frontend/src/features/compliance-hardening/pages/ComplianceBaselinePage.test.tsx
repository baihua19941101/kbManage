import { canManageComplianceBaseline, canReadCompliance } from '@/features/auth/store';

describe('ComplianceBaselinePage', () => {
  it('exposes compliance baseline permissions for page gate', () => {
    expect(canReadCompliance({ id: 'u1', username: 'admin', platformRoles: ['platform-admin'] })).toBe(true);
    expect(canManageComplianceBaseline({ id: 'u1', username: 'admin', platformRoles: ['platform-admin'] })).toBe(true);
  });
});
