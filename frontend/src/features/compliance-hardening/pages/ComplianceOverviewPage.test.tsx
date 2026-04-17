import { queryKeys } from '@/app/queryClient';

describe('ComplianceOverviewPage', () => {
  it('uses compliance overview query key', () => {
    expect(queryKeys.compliance.overview('cluster')).toEqual(['compliance', 'overview', 'cluster']);
  });
});
