import {
  canManageEnterpriseAudit,
  canManageEnterpriseDelivery,
  canManageEnterpriseReports,
  canReadEnterprisePolish,
  useAuthStore
} from '@/features/auth/store';

export const useEnterprisePermissions = () => {
  const user = useAuthStore((state) => state.user);
  return {
    canRead: canReadEnterprisePolish(user),
    canManageAudit: canManageEnterpriseAudit(user),
    canManageReports: canManageEnterpriseReports(user),
    canManageDelivery: canManageEnterpriseDelivery(user)
  };
};
