import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AppLayout, Home } from '@/app/App';
import { ProtectedRoute } from '@/app/ProtectedRoute';

export const router = createBrowserRouter([
  {
    path: '/login',
    lazy: async () => {
      const { LoginPage } = await import('@/features/auth/pages/LoginPage');
      return { Component: LoginPage };
    }
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: '/',
        element: <AppLayout />,
        children: [
          {
            index: true,
            element: <Home />
          },
          {
            path: '/clusters',
            lazy: async () => {
              const { ClusterOverviewPage } = await import(
                '@/features/clusters/pages/ClusterOverviewPage'
              );
              return { Component: ClusterOverviewPage };
            }
          },
          {
            path: '/resources',
            lazy: async () => {
              const { ResourcesPage } = await import('@/features/resources/pages/ResourcesPage');
              return { Component: ResourcesPage };
            }
          },
          {
            path: '/workspaces',
            lazy: async () => {
              const { WorkspacePage } = await import('@/features/workspaces/pages/WorkspacePage');
              return { Component: WorkspacePage };
            }
          },
          {
            path: '/projects',
            lazy: async () => {
              const { ProjectPage } = await import('@/features/projects/pages/ProjectPage');
              return { Component: ProjectPage };
            }
          },
          {
            path: '/audit-events',
            lazy: async () => {
              const { AuditEventPage } = await import('@/features/audit/pages/AuditEventPage');
              return { Component: AuditEventPage };
            }
          },
          {
            path: '/audit-events/gitops',
            lazy: async () => {
              const { GitOpsAuditPage } = await import('@/features/audit/pages/GitOpsAuditPage');
              return { Component: GitOpsAuditPage };
            }
          },
          {
            path: '/audit-events/security-policy',
            lazy: async () => {
              const { SecurityPolicyAuditPage } = await import(
                '@/features/audit/pages/SecurityPolicyAuditPage'
              );
              return { Component: SecurityPolicyAuditPage };
            }
          },
          {
            path: '/audit-events/compliance',
            lazy: async () => {
              const { ComplianceAuditPage } = await import(
                '@/features/audit/pages/ComplianceAuditPage'
              );
              return { Component: ComplianceAuditPage };
            }
          },
          {
            path: '/audit-events/cluster-lifecycle',
            lazy: async () => {
              const { ClusterLifecycleAuditPage } = await import(
                '@/features/audit/pages/ClusterLifecycleAuditPage'
              );
              return { Component: ClusterLifecycleAuditPage };
            }
          },
          {
            path: '/audit-events/backup-restore',
            lazy: async () => {
              const { BackupRestoreAuditPage } = await import(
                '@/features/audit/pages/BackupRestoreAuditPage'
              );
              return { Component: BackupRestoreAuditPage };
            }
          },
          {
            path: '/audit-events/identity',
            lazy: async () => {
              const { IdentityGovernanceAuditPage } = await import(
                '@/features/audit/pages/IdentityGovernanceAuditPage'
              );
              return { Component: IdentityGovernanceAuditPage };
            }
          },
          {
            path: '/audit-events/platform-marketplace',
            lazy: async () => {
              const { PlatformMarketplaceAuditPage } = await import(
                '@/features/audit/pages/PlatformMarketplaceAuditPage'
              );
              return { Component: PlatformMarketplaceAuditPage };
            }
          },
          {
            path: '/audit-events/sre',
            lazy: async () => {
              const { SREAuditPage } = await import('@/features/audit/pages/SREAuditPage');
              return { Component: SREAuditPage };
            }
          },
          {
            path: '/observability',
            lazy: async () => {
              const { ObservabilityOverviewPage } = await import(
                '@/features/observability/pages/ObservabilityOverviewPage'
              );
              return { Component: ObservabilityOverviewPage };
            }
          },
          {
            path: '/observability/logs',
            lazy: async () => {
              const { LogExplorerPage } = await import(
                '@/features/observability/pages/LogExplorerPage'
              );
              return { Component: LogExplorerPage };
            }
          },
          {
            path: '/observability/events',
            lazy: async () => {
              const { EventExplorerPage } = await import(
                '@/features/observability/pages/EventExplorerPage'
              );
              return { Component: EventExplorerPage };
            }
          },
          {
            path: '/observability/metrics',
            lazy: async () => {
              const { MetricsExplorerPage } = await import(
                '@/features/observability/pages/MetricsExplorerPage'
              );
              return { Component: MetricsExplorerPage };
            }
          },
          {
            path: '/observability/context',
            lazy: async () => {
              const { ResourceContextPage } = await import(
                '@/features/observability/pages/ResourceContextPage'
              );
              return { Component: ResourceContextPage };
            }
          },
          {
            path: '/observability/alerts',
            lazy: async () => {
              const { AlertCenterPage } = await import(
                '@/features/observability/pages/AlertCenterPage'
              );
              return { Component: AlertCenterPage };
            }
          },
          {
            path: '/observability/alert-rules',
            lazy: async () => {
              const { AlertRulePage } = await import('@/features/observability/pages/AlertRulePage');
              return { Component: AlertRulePage };
            }
          },
          {
            path: '/observability/silences',
            lazy: async () => {
              const { SilenceWindowPage } = await import(
                '@/features/observability/pages/SilenceWindowPage'
              );
              return { Component: SilenceWindowPage };
            }
          },
          {
            path: '/observability/*',
            element: <Navigate to="/observability" replace />
          },
          {
            path: '/workload-ops',
            lazy: async () => {
              const { WorkloadOperationsPage } = await import(
                '@/features/workload-ops/pages/WorkloadOperationsPage'
              );
              return { Component: WorkloadOperationsPage };
            }
          },
          {
            path: '/workload-ops/batches',
            lazy: async () => {
              const { BatchOperationPage } = await import(
                '@/features/workload-ops/pages/BatchOperationPage'
              );
              return { Component: BatchOperationPage };
            }
          },
          {
            path: '/gitops',
            lazy: async () => {
              const { GitOpsOverviewPage } = await import('@/features/gitops/pages/GitOpsOverviewPage');
              return { Component: GitOpsOverviewPage };
            }
          },
          {
            path: '/gitops/delivery-units/:unitId',
            lazy: async () => {
              const { DeliveryUnitDetailPage } = await import(
                '@/features/gitops/pages/DeliveryUnitDetailPage'
              );
              return { Component: DeliveryUnitDetailPage };
            }
          },
          {
            path: '/gitops/*',
            element: <Navigate to="/gitops" replace />
          },
          {
            path: '/security-policies',
            lazy: async () => {
              const { PolicyCenterPage } = await import(
                '@/features/security-policy/pages/PolicyCenterPage'
              );
              return { Component: PolicyCenterPage };
            }
          },
          {
            path: '/security-policies/rollout',
            lazy: async () => {
              const { PolicyRolloutPage } = await import(
                '@/features/security-policy/pages/PolicyRolloutPage'
              );
              return { Component: PolicyRolloutPage };
            }
          },
          {
            path: '/security-policies/violations',
            lazy: async () => {
              const { ViolationCenterPage } = await import(
                '@/features/security-policy/pages/ViolationCenterPage'
              );
              return { Component: ViolationCenterPage };
            }
          },
          {
            path: '/cluster-lifecycle',
            lazy: async () => {
              const { ClusterLifecycleListPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterLifecycleListPage'
              );
              return { Component: ClusterLifecycleListPage };
            }
          },
          {
            path: '/cluster-lifecycle/register',
            lazy: async () => {
              const { ClusterRegistrationPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterRegistrationPage'
              );
              return { Component: ClusterRegistrationPage };
            }
          },
          {
            path: '/cluster-lifecycle/provision',
            lazy: async () => {
              const { ClusterProvisionPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterProvisionPage'
              );
              return { Component: ClusterProvisionPage };
            }
          },
          {
            path: '/cluster-lifecycle/upgrades',
            lazy: async () => {
              const { ClusterUpgradePage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterUpgradePage'
              );
              return { Component: ClusterUpgradePage };
            }
          },
          {
            path: '/cluster-lifecycle/node-pools',
            lazy: async () => {
              const { NodePoolPage } = await import(
                '@/features/cluster-lifecycle/pages/NodePoolPage'
              );
              return { Component: NodePoolPage };
            }
          },
          {
            path: '/cluster-lifecycle/retirement',
            lazy: async () => {
              const { ClusterRetirementPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterRetirementPage'
              );
              return { Component: ClusterRetirementPage };
            }
          },
          {
            path: '/cluster-lifecycle/drivers',
            lazy: async () => {
              const { ClusterDriverPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterDriverPage'
              );
              return { Component: ClusterDriverPage };
            }
          },
          {
            path: '/cluster-lifecycle/templates',
            lazy: async () => {
              const { ClusterTemplatePage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterTemplatePage'
              );
              return { Component: ClusterTemplatePage };
            }
          },
          {
            path: '/cluster-lifecycle/capabilities',
            lazy: async () => {
              const { CapabilityMatrixPage } = await import(
                '@/features/cluster-lifecycle/pages/CapabilityMatrixPage'
              );
              return { Component: CapabilityMatrixPage };
            }
          },
          {
            path: '/cluster-lifecycle/:clusterId',
            lazy: async () => {
              const { ClusterLifecycleDetailPage } = await import(
                '@/features/cluster-lifecycle/pages/ClusterLifecycleDetailPage'
              );
              return { Component: ClusterLifecycleDetailPage };
            }
          },
          {
            path: '/cluster-lifecycle/*',
            element: <Navigate to="/cluster-lifecycle" replace />
          },
          {
            path: '/backup-restore',
            lazy: async () => {
              const { BackupPolicyPage } = await import(
                '@/features/backup-restore-dr/pages/BackupPolicyPage'
              );
              return { Component: BackupPolicyPage };
            }
          },
          {
            path: '/backup-restore/policies',
            lazy: async () => {
              const { BackupPolicyPage } = await import(
                '@/features/backup-restore-dr/pages/BackupPolicyPage'
              );
              return { Component: BackupPolicyPage };
            }
          },
          {
            path: '/backup-restore/restore-points',
            lazy: async () => {
              const { RestorePointPage } = await import(
                '@/features/backup-restore-dr/pages/RestorePointPage'
              );
              return { Component: RestorePointPage };
            }
          },
          {
            path: '/backup-restore/restore-jobs',
            lazy: async () => {
              const { RestoreJobPage } = await import(
                '@/features/backup-restore-dr/pages/RestoreJobPage'
              );
              return { Component: RestoreJobPage };
            }
          },
          {
            path: '/backup-restore/migrations',
            lazy: async () => {
              const { MigrationPlanPage } = await import(
                '@/features/backup-restore-dr/pages/MigrationPlanPage'
              );
              return { Component: MigrationPlanPage };
            }
          },
          {
            path: '/backup-restore/drills',
            lazy: async () => {
              const { DRDrillPlanPage } = await import(
                '@/features/backup-restore-dr/pages/DRDrillPlanPage'
              );
              return { Component: DRDrillPlanPage };
            }
          },
          {
            path: '/backup-restore/drills/records/:recordId',
            lazy: async () => {
              const { DRDrillRecordPage } = await import(
                '@/features/backup-restore-dr/pages/DRDrillRecordPage'
              );
              return { Component: DRDrillRecordPage };
            }
          },
          {
            path: '/backup-restore/drills/reports/:recordId',
            lazy: async () => {
              const { DRDrillReportPage } = await import(
                '@/features/backup-restore-dr/pages/DRDrillReportPage'
              );
              return { Component: DRDrillReportPage };
            }
          },
          {
            path: '/backup-restore/*',
            element: <Navigate to="/backup-restore" replace />
          },
          {
            path: '/identity-tenancy',
            lazy: async () => {
              const { IdentitySourcePage } = await import(
                '@/features/identity-tenancy/pages/IdentitySourcePage'
              );
              return { Component: IdentitySourcePage };
            }
          },
          {
            path: '/identity-tenancy/sources',
            lazy: async () => {
              const { IdentitySourcePage } = await import(
                '@/features/identity-tenancy/pages/IdentitySourcePage'
              );
              return { Component: IdentitySourcePage };
            }
          },
          {
            path: '/identity-tenancy/sessions',
            lazy: async () => {
              const { SessionGovernancePage } = await import(
                '@/features/identity-tenancy/pages/SessionGovernancePage'
              );
              return { Component: SessionGovernancePage };
            }
          },
          {
            path: '/identity-tenancy/organizations',
            lazy: async () => {
              const { OrganizationModelPage } = await import(
                '@/features/identity-tenancy/pages/OrganizationModelPage'
              );
              return { Component: OrganizationModelPage };
            }
          },
          {
            path: '/identity-tenancy/mappings',
            lazy: async () => {
              const { TenantMappingPage } = await import(
                '@/features/identity-tenancy/pages/TenantMappingPage'
              );
              return { Component: TenantMappingPage };
            }
          },
          {
            path: '/identity-tenancy/roles',
            lazy: async () => {
              const { RoleCatalogPage } = await import(
                '@/features/identity-tenancy/pages/RoleCatalogPage'
              );
              return { Component: RoleCatalogPage };
            }
          },
          {
            path: '/identity-tenancy/assignments',
            lazy: async () => {
              const { RoleAssignmentPage } = await import(
                '@/features/identity-tenancy/pages/RoleAssignmentPage'
              );
              return { Component: RoleAssignmentPage };
            }
          },
          {
            path: '/identity-tenancy/delegations',
            lazy: async () => {
              const { DelegationPage } = await import(
                '@/features/identity-tenancy/pages/DelegationPage'
              );
              return { Component: DelegationPage };
            }
          },
          {
            path: '/identity-tenancy/access-risks',
            lazy: async () => {
              const { AccessRiskPage } = await import(
                '@/features/identity-tenancy/pages/AccessRiskPage'
              );
              return { Component: AccessRiskPage };
            }
          },
          {
            path: '/identity-tenancy/*',
            element: <Navigate to="/identity-tenancy" replace />
          },
          {
            path: '/platform-marketplace',
            lazy: async () => {
              const { TemplateCatalogPage } = await import(
                '@/features/platform-marketplace/pages/TemplateCatalogPage'
              );
              return { Component: TemplateCatalogPage };
            }
          },
          {
            path: '/platform-marketplace/catalog-sources',
            lazy: async () => {
              const { CatalogSourcePage } = await import(
                '@/features/platform-marketplace/pages/CatalogSourcePage'
              );
              return { Component: CatalogSourcePage };
            }
          },
          {
            path: '/platform-marketplace/templates',
            lazy: async () => {
              const { TemplateCatalogPage } = await import(
                '@/features/platform-marketplace/pages/TemplateCatalogPage'
              );
              return { Component: TemplateCatalogPage };
            }
          },
          {
            path: '/platform-marketplace/templates/:templateId',
            lazy: async () => {
              const { TemplateDetailPage } = await import(
                '@/features/platform-marketplace/pages/TemplateDetailPage'
              );
              return { Component: TemplateDetailPage };
            }
          },
          {
            path: '/platform-marketplace/distribution',
            lazy: async () => {
              const { TemplateDistributionPage } = await import(
                '@/features/platform-marketplace/pages/TemplateDistributionPage'
              );
              return { Component: TemplateDistributionPage };
            }
          },
          {
            path: '/platform-marketplace/installations',
            lazy: async () => {
              const { InstallationRecordPage } = await import(
                '@/features/platform-marketplace/pages/InstallationRecordPage'
              );
              return { Component: InstallationRecordPage };
            }
          },
          {
            path: '/platform-marketplace/extensions',
            lazy: async () => {
              const { ExtensionCenterPage } = await import(
                '@/features/platform-marketplace/pages/ExtensionCenterPage'
              );
              return { Component: ExtensionCenterPage };
            }
          },
          {
            path: '/platform-marketplace/compatibility',
            lazy: async () => {
              const { ExtensionCompatibilityPage } = await import(
                '@/features/platform-marketplace/pages/ExtensionCompatibilityPage'
              );
              return { Component: ExtensionCompatibilityPage };
            }
          },
          {
            path: '/platform-marketplace/*',
            element: <Navigate to="/platform-marketplace" replace />
          },
          {
            path: '/sre-scale',
            lazy: async () => {
              const { HealthOverviewPage } = await import('@/features/sre-scale/pages/HealthOverviewPage');
              return { Component: HealthOverviewPage };
            }
          },
          {
            path: '/sre-scale/ha',
            lazy: async () => {
              const { HAControlPage } = await import('@/features/sre-scale/pages/HAControlPage');
              return { Component: HAControlPage };
            }
          },
          {
            path: '/sre-scale/upgrades',
            lazy: async () => {
              const { UpgradeGovernancePage } = await import('@/features/sre-scale/pages/UpgradeGovernancePage');
              return { Component: UpgradeGovernancePage };
            }
          },
          {
            path: '/sre-scale/rollback',
            lazy: async () => {
              const { RollbackValidationPage } = await import('@/features/sre-scale/pages/RollbackValidationPage');
              return { Component: RollbackValidationPage };
            }
          },
          {
            path: '/sre-scale/capacity',
            lazy: async () => {
              const { CapacityGovernancePage } = await import('@/features/sre-scale/pages/CapacityGovernancePage');
              return { Component: CapacityGovernancePage };
            }
          },
          {
            path: '/sre-scale/runbooks',
            lazy: async () => {
              const { RunbookCenterPage } = await import('@/features/sre-scale/pages/RunbookCenterPage');
              return { Component: RunbookCenterPage };
            }
          },
          {
            path: '/sre-scale/*',
            element: <Navigate to="/sre-scale" replace />
          },
          {
            path: '/compliance-hardening/baselines',
            lazy: async () => {
              const { ComplianceBaselinePage } = await import(
                '@/features/compliance-hardening/pages/ComplianceBaselinePage'
              );
              return { Component: ComplianceBaselinePage };
            }
          },
          {
            path: '/compliance-hardening/scans',
            lazy: async () => {
              const { ScanCenterPage } = await import(
                '@/features/compliance-hardening/pages/ScanCenterPage'
              );
              return { Component: ScanCenterPage };
            }
          },
          {
            path: '/compliance-hardening/findings/:findingId',
            lazy: async () => {
              const { FindingDetailPage } = await import(
                '@/features/compliance-hardening/pages/FindingDetailPage'
              );
              return { Component: FindingDetailPage };
            }
          },
          {
            path: '/compliance-hardening/remediation',
            lazy: async () => {
              const { RemediationQueuePage } = await import(
                '@/features/compliance-hardening/pages/RemediationQueuePage'
              );
              return { Component: RemediationQueuePage };
            }
          },
          {
            path: '/compliance-hardening/exceptions',
            lazy: async () => {
              const { ComplianceExceptionPage } = await import(
                '@/features/compliance-hardening/pages/ComplianceExceptionPage'
              );
              return { Component: ComplianceExceptionPage };
            }
          },
          {
            path: '/compliance-hardening/rechecks',
            lazy: async () => {
              const { RecheckCenterPage } = await import(
                '@/features/compliance-hardening/pages/RecheckCenterPage'
              );
              return { Component: RecheckCenterPage };
            }
          },
          {
            path: '/compliance-hardening/overview',
            lazy: async () => {
              const { ComplianceOverviewPage } = await import(
                '@/features/compliance-hardening/pages/ComplianceOverviewPage'
              );
              return { Component: ComplianceOverviewPage };
            }
          },
          {
            path: '/compliance-hardening/trends',
            lazy: async () => {
              const { ComplianceTrendPage } = await import(
                '@/features/compliance-hardening/pages/ComplianceTrendPage'
              );
              return { Component: ComplianceTrendPage };
            }
          },
          {
            path: '/compliance-hardening/archive',
            lazy: async () => {
              const { ComplianceArchivePage } = await import(
                '@/features/compliance-hardening/pages/ComplianceArchivePage'
              );
              return { Component: ComplianceArchivePage };
            }
          },
          {
            path: '/compliance-hardening/*',
            element: <Navigate to="/compliance-hardening/overview" replace />
          }
        ]
      }
    ]
  },
  {
    path: '*',
    element: <Navigate to="/" replace />
  }
]);
