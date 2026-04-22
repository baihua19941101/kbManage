import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createExportRecord, createGovernanceReport } from '@/services/enterprisePolish';
import { queryKeys } from '@/app/queryClient';

export const useReportActions = () => {
  const queryClient = useQueryClient();
  const createReport = useMutation({
    mutationFn: createGovernanceReport,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.enterprisePolish.reports() })
  });
  const createExport = useMutation({
    mutationFn: ({ reportId, audienceScope, contentLevel, exportType }: { reportId: string; audienceScope: string; contentLevel: string; exportType: string }) =>
      createExportRecord(reportId, { audienceScope, contentLevel, exportType })
  });
  return { createReport, createExport };
};
