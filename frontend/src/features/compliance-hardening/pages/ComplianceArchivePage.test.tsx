import { canExportComplianceArchive } from '@/features/auth/store';

describe('ComplianceArchivePage', () => {
  it('uses compliance archive export permission gate', () => {
    expect(canExportComplianceArchive({ id: 'u1', username: 'auditor', platformRoles: ['audit-reader'] })).toBe(true);
    expect(canExportComplianceArchive({ id: 'u2', username: 'viewer', platformRoles: [] })).toBe(false);
  });
});
