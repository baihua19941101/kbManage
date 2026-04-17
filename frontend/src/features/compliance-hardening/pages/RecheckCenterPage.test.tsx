import { canReadCompliance } from '@/features/auth/store';

describe('RecheckCenterPage', () => {
  it('uses compliance read permission gate', () => {
    expect(canReadCompliance({ id: 'u1', username: 'admin', platformRoles: ['platform-admin'] })).toBe(true);
    expect(canReadCompliance({ id: 'u2', username: 'viewer', platformRoles: [] })).toBe(false);
  });
});
