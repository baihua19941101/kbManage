package enterprise

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	entint "kbmanage/backend/internal/integration/enterprise"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
)

const ResourceTypeEnterprise = "enterprise"

var (
	ErrEnterpriseScopeDenied = errors.New("enterprise scope access denied")
	ErrEnterpriseInvalid     = errors.New("enterprise invalid request")
)

type Service struct {
	trails            *repository.PermissionChangeTrailRepository
	operations        *repository.KeyOperationTraceRepository
	crossTeams        *repository.CrossTeamAuthorizationSnapshotRepository
	risks             *repository.GovernanceRiskEventRepository
	coverage          *repository.GovernanceCoverageSnapshotRepository
	reports           *repository.GovernanceReportPackageRepository
	exports           *repository.ExportRecordRepository
	artifacts         *repository.DeliveryArtifactRepository
	bundles           *repository.DeliveryReadinessBundleRepository
	checklists        *repository.DeliveryChecklistItemRepository
	actions           *repository.GovernanceActionItemRepository
	scope             *ScopeService
	auditProvider     entint.AuditProvider
	reportBuilder     entint.ReportBuilder
	deliveryCatalog   entint.DeliveryCatalog
	reportCache       *ReportCache
	exportCoordinator *ExportCoordinator
	trendCache        *TrendCache
	auditWriter       *auditSvc.EventWriter
}

type GovernanceReportInput struct {
	WorkspaceID      uint64  `json:"workspaceId"`
	ProjectID        *uint64 `json:"projectId"`
	ReportType       string  `json:"reportType"`
	Title            string  `json:"title"`
	AudienceType     string  `json:"audienceType"`
	TimeRange        string  `json:"timeRange"`
	VisibilityPolicy string  `json:"visibilityPolicy"`
}

type ExportRecordInput struct {
	AudienceScope string `json:"audienceScope"`
	ContentLevel  string `json:"contentLevel"`
	ExportType    string `json:"exportType"`
}

func NewService(
	trails *repository.PermissionChangeTrailRepository,
	operations *repository.KeyOperationTraceRepository,
	crossTeams *repository.CrossTeamAuthorizationSnapshotRepository,
	risks *repository.GovernanceRiskEventRepository,
	coverage *repository.GovernanceCoverageSnapshotRepository,
	reports *repository.GovernanceReportPackageRepository,
	exports *repository.ExportRecordRepository,
	artifacts *repository.DeliveryArtifactRepository,
	bundles *repository.DeliveryReadinessBundleRepository,
	checklists *repository.DeliveryChecklistItemRepository,
	actions *repository.GovernanceActionItemRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	workspaceClusterRepo *repository.WorkspaceClusterRepository,
	auditProvider entint.AuditProvider,
	reportBuilder entint.ReportBuilder,
	deliveryCatalog entint.DeliveryCatalog,
	reportCache *ReportCache,
	exportCoordinator *ExportCoordinator,
	trendCache *TrendCache,
	auditWriter *auditSvc.EventWriter,
) *Service {
	if auditProvider == nil {
		auditProvider = entint.NewStaticAuditProvider()
	}
	if reportBuilder == nil {
		reportBuilder = entint.NewStaticReportBuilder()
	}
	if deliveryCatalog == nil {
		deliveryCatalog = entint.NewStaticDeliveryCatalog()
	}
	return &Service{
		trails: trails, operations: operations, crossTeams: crossTeams, risks: risks, coverage: coverage,
		reports: reports, exports: exports, artifacts: artifacts, bundles: bundles, checklists: checklists, actions: actions,
		scope:         NewScopeService(bindingRepo, projectRepo, workspaceClusterRepo),
		auditProvider: auditProvider, reportBuilder: reportBuilder, deliveryCatalog: deliveryCatalog,
		reportCache: reportCache, exportCoordinator: exportCoordinator, trendCache: trendCache, auditWriter: auditWriter,
	}
}

func (s *Service) writeAudit(ctx context.Context, actorID uint64, action, targetType, targetRef string, outcome domain.AuditOutcome, details map[string]any) {
	if s.auditWriter == nil {
		return
	}
	actor := actorID
	_ = s.auditWriter.Write(ctx, "", &actor, action, ResourceTypeEnterprise, firstNonEmptyString(targetRef, targetType), outcome, details)
}

func (s *Service) ListPermissionTrails(ctx context.Context, userID uint64) ([]domain.PermissionChangeTrail, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	items, err := s.trails.List(ctx)
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, "enterprise.permission-trail.read", "permission-trail", "", domain.AuditOutcomeSuccess, nil)
	return items, nil
}

func (s *Service) ListKeyOperations(ctx context.Context, userID uint64) ([]domain.KeyOperationTrace, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	items, err := s.operations.List(ctx)
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, "enterprise.key-operation.read", "key-operation", "", domain.AuditOutcomeSuccess, nil)
	return items, nil
}

func (s *Service) ListCoverageSnapshots(ctx context.Context, userID uint64) ([]domain.GovernanceCoverageSnapshot, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	items, err := s.coverage.List(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.trendCache.Store(ctx, "coverage", items)
	return items, nil
}

func (s *Service) ListActionItems(ctx context.Context, userID uint64) ([]domain.GovernanceActionItem, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	return s.actions.List(ctx)
}

func (s *Service) ListReports(ctx context.Context, userID uint64) ([]domain.GovernanceReportPackage, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	return s.reports.List(ctx)
}

func (s *Service) CreateReport(ctx context.Context, userID uint64, input GovernanceReportInput) (*domain.GovernanceReportPackage, error) {
	if err := s.scope.EnsureScopePermission(ctx, userID, input.WorkspaceID, input.ProjectID, "enterprise:manage-reports"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.ReportType) == "" || strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.AudienceType) == "" {
		return nil, ErrEnterpriseInvalid
	}
	built := s.reportBuilder.Build(ctx, entint.ReportBuildInput{ReportType: input.ReportType, AudienceType: input.AudienceType, TimeRange: input.TimeRange})
	item := &domain.GovernanceReportPackage{
		WorkspaceID:       input.WorkspaceID,
		ProjectID:         input.ProjectID,
		ReportType:        strings.TrimSpace(input.ReportType),
		Title:             strings.TrimSpace(input.Title),
		AudienceType:      strings.TrimSpace(input.AudienceType),
		TimeRange:         strings.TrimSpace(input.TimeRange),
		SummarySection:    built.SummarySection,
		DetailSection:     built.DetailSection,
		AttachmentCatalog: strings.Join(built.AttachmentCatalog, ","),
		VisibilityPolicy:  firstNonEmptyString(input.VisibilityPolicy, "default"),
		GeneratedAt:       time.Now(),
		GeneratedBy:       userID,
		Status:            domain.GovernanceReportStatusReady,
	}
	if err := s.reports.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.reportCache.Store(ctx, strconv.FormatUint(item.ID, 10), item)
	s.writeAudit(ctx, userID, "enterprise.report.create", "report", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"reportType": item.ReportType})
	return item, nil
}

func (s *Service) CreateExportRecord(ctx context.Context, userID, reportID uint64, input ExportRecordInput) (*domain.ExportRecord, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:manage-reports"); err != nil {
		return nil, err
	}
	if _, err := s.reports.GetByID(ctx, reportID); err != nil {
		return nil, err
	}
	item := &domain.ExportRecord{
		PackageID:      reportID,
		AudienceScope:  firstNonEmptyString(input.AudienceScope, "default"),
		ContentLevel:   firstNonEmptyString(input.ContentLevel, "summary"),
		ExportType:     firstNonEmptyString(input.ExportType, "report"),
		Result:         "succeeded",
		AuditReference: "enterprise-export-" + strconv.FormatUint(reportID, 10),
		ExportedAt:     time.Now(),
		ExportedBy:     userID,
	}
	if err := s.exports.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.exportCoordinator.Lock(ctx, strconv.FormatUint(item.ID, 10))
	s.writeAudit(ctx, userID, "enterprise.export.create", "export-record", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"reportId": reportID})
	return item, nil
}

func (s *Service) ListDeliveryArtifacts(ctx context.Context, userID uint64) ([]domain.DeliveryArtifact, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	return s.artifacts.List(ctx)
}

func (s *Service) ListDeliveryBundles(ctx context.Context, userID uint64) ([]domain.DeliveryReadinessBundle, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	return s.bundles.List(ctx)
}

func (s *Service) ListDeliveryChecklist(ctx context.Context, userID, bundleID uint64) ([]domain.DeliveryChecklistItem, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "enterprise:read"); err != nil {
		return nil, err
	}
	return s.checklists.ListByBundleID(ctx, bundleID)
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
