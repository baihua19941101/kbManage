package marketplace

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	marketplaceint "kbmanage/backend/internal/integration/marketplace"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"

	"gorm.io/gorm"
)

const (
	ResourceTypePlatformMarketplace = "platformmarketplace"

	ActionCatalogSourceCreate = "platformmarketplace.catalog-source.create"
	ActionCatalogSourceSync   = "platformmarketplace.catalog-source.sync"
	ActionTemplateRelease     = "platformmarketplace.template.release"
	ActionExtensionRegister   = "platformmarketplace.extension.register"
	ActionExtensionEnable     = "platformmarketplace.extension.enable"
	ActionExtensionDisable    = "platformmarketplace.extension.disable"
)

var (
	ErrMarketplaceScopeDenied = errors.New("marketplace scope access denied")
	ErrMarketplaceConflict    = errors.New("marketplace operation conflict")
	ErrMarketplaceInvalid     = errors.New("marketplace invalid request")
	ErrMarketplaceBlocked     = errors.New("marketplace operation blocked")
)

type Service struct {
	sources       *repository.CatalogSourceRepository
	templates     *repository.ApplicationTemplateRepository
	versions      *repository.TemplateVersionRepository
	releases      *repository.TemplateReleaseScopeRepository
	installations *repository.InstallationRecordRepository
	extensions    *repository.ExtensionPackageRepository
	compatibility *repository.CompatibilityStatementRepository
	lifecycle     *repository.ExtensionLifecycleRepository
	localAudit    *repository.MarketplaceAuditRepository
	scope         *ScopeService
	catalog       marketplaceint.CatalogProvider
	registry      marketplaceint.ExtensionRegistry
	auditWriter   *auditSvc.EventWriter
}

func NewService(
	sourceRepo *repository.CatalogSourceRepository,
	templateRepo *repository.ApplicationTemplateRepository,
	versionRepo *repository.TemplateVersionRepository,
	releaseRepo *repository.TemplateReleaseScopeRepository,
	installationRepo *repository.InstallationRecordRepository,
	extensionRepo *repository.ExtensionPackageRepository,
	compatibilityRepo *repository.CompatibilityStatementRepository,
	lifecycleRepo *repository.ExtensionLifecycleRepository,
	localAuditRepo *repository.MarketplaceAuditRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	workspaceClusterRepo *repository.WorkspaceClusterRepository,
	catalogProvider marketplaceint.CatalogProvider,
	registry marketplaceint.ExtensionRegistry,
	auditWriter *auditSvc.EventWriter,
) *Service {
	if catalogProvider == nil {
		catalogProvider = marketplaceint.NewStaticCatalogProvider()
	}
	if registry == nil {
		registry = marketplaceint.NewStaticExtensionRegistry()
	}
	return &Service{
		sources:       sourceRepo,
		templates:     templateRepo,
		versions:      versionRepo,
		releases:      releaseRepo,
		installations: installationRepo,
		extensions:    extensionRepo,
		compatibility: compatibilityRepo,
		lifecycle:     lifecycleRepo,
		localAudit:    localAuditRepo,
		scope:         NewScopeService(bindingRepo, projectRepo, workspaceClusterRepo),
		catalog:       catalogProvider,
		registry:      registry,
		auditWriter:   auditWriter,
	}
}

type CatalogSourceListFilter struct {
	SourceType string
	Status     string
	Keyword    string
}

type TemplateListFilter struct {
	CatalogSourceID uint64
	Category        string
	Status          string
	Keyword         string
}

type InstallationListFilter struct {
	ScopeType string
	ScopeID   uint64
	Status    string
}

type ExtensionListFilter struct {
	Type    string
	Status  string
	Keyword string
}

type CreateCatalogSourceInput struct {
	Name            string                        `json:"name"`
	SourceType      string                        `json:"sourceType"`
	EndpointRef     string                        `json:"endpointRef"`
	Status          string                        `json:"status"`
	VisibilityScope string                        `json:"visibilityScope"`
	TemplateSeeds   []marketplaceint.TemplateSeed `json:"templateSeeds"`
}

type CreateTemplateReleaseInput struct {
	VersionID      uint64 `json:"versionId"`
	Version        string `json:"version"`
	ScopeID        uint64 `json:"scopeId"`
	ScopeType      string `json:"scopeType"`
	TargetRef      string `json:"targetRef"`
	VisibilityMode string `json:"visibilityMode"`
}

type CreateExtensionPackageInput struct {
	Name                  string                             `json:"name"`
	ExtensionType         string                             `json:"extensionType"`
	Version               string                             `json:"version"`
	VisibilityScope       string                             `json:"visibilityScope"`
	EntrySummary          string                             `json:"entrySummary"`
	PermissionDeclaration []string                           `json:"permissionDeclaration"`
	Compatibility         []marketplaceint.CompatibilitySeed `json:"compatibility"`
}

type ExtensionLifecycleInput struct {
	ScopeType string `json:"scopeType"`
	ScopeID   uint64 `json:"scopeId"`
	Reason    string `json:"reason"`
}

type TemplateDetail struct {
	Template      *domain.ApplicationTemplate     `json:"template"`
	Versions      []domain.TemplateVersion        `json:"versions"`
	Releases      []domain.TemplateReleaseScope   `json:"releases"`
	Compatibility []domain.CompatibilityStatement `json:"compatibility"`
}

type ExtensionCompatibilityView struct {
	Extension      *domain.ExtensionPackage        `json:"extension"`
	Statements     []domain.CompatibilityStatement `json:"statements"`
	BlockedReasons []string                        `json:"blockedReasons"`
}

func (s *Service) writeAudit(ctx context.Context, actorID uint64, action, targetType, targetRef string, outcome domain.AuditOutcome, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	if s.localAudit != nil {
		payload, _ := json.Marshal(details)
		_ = s.localAudit.Create(ctx, &domain.MarketplaceAuditEvent{
			Action:         action,
			ActorUserID:    actorID,
			TargetType:     targetType,
			TargetRef:      targetRef,
			Outcome:        string(outcome),
			DetailSnapshot: string(payload),
			OccurredAt:     time.Now(),
		})
	}
	if s.auditWriter != nil {
		actor := actorID
		_ = s.auditWriter.Write(
			ctx,
			"",
			&actor,
			action,
			ResourceTypePlatformMarketplace,
			targetRef,
			outcome,
			details,
		)
	}
}

func normalizeSourceStatus(value string) domain.CatalogSourceStatus {
	switch strings.TrimSpace(value) {
	case string(domain.CatalogSourceStatusDraft):
		return domain.CatalogSourceStatusDraft
	case string(domain.CatalogSourceStatusDisabled):
		return domain.CatalogSourceStatusDisabled
	case string(domain.CatalogSourceStatusDegraded):
		return domain.CatalogSourceStatusDegraded
	default:
		return domain.CatalogSourceStatusActive
	}
}

func normalizeTemplatePublishStatus(value string) domain.TemplatePublishStatus {
	switch strings.TrimSpace(value) {
	case string(domain.TemplatePublishStatusDisabled):
		return domain.TemplatePublishStatusDisabled
	case string(domain.TemplatePublishStatusRetired):
		return domain.TemplatePublishStatusRetired
	case string(domain.TemplatePublishStatusHistoryOnly):
		return domain.TemplatePublishStatusHistoryOnly
	default:
		return domain.TemplatePublishStatusActive
	}
}

func normalizeTemplateVersionStatus(value string) domain.TemplateVersionStatus {
	switch strings.TrimSpace(value) {
	case string(domain.TemplateVersionStatusDraft):
		return domain.TemplateVersionStatusDraft
	case string(domain.TemplateVersionStatusDeprecated):
		return domain.TemplateVersionStatusDeprecated
	case string(domain.TemplateVersionStatusRetired):
		return domain.TemplateVersionStatusRetired
	default:
		return domain.TemplateVersionStatusActive
	}
}

func normalizeExtensionStatus(value string) domain.ExtensionPackageStatus {
	switch strings.TrimSpace(value) {
	case string(domain.ExtensionPackageStatusDraft):
		return domain.ExtensionPackageStatusDraft
	case string(domain.ExtensionPackageStatusEnabled):
		return domain.ExtensionPackageStatusEnabled
	case string(domain.ExtensionPackageStatusDisabled):
		return domain.ExtensionPackageStatusDisabled
	case string(domain.ExtensionPackageStatusRetired):
		return domain.ExtensionPackageStatusRetired
	default:
		return domain.ExtensionPackageStatusRegistered
	}
}

func scopeRef(scopeType string, scopeID uint64) string {
	return strings.TrimSpace(scopeType) + ":" + strconv.FormatUint(scopeID, 10)
}

func mustJSON(v any) string {
	if v == nil {
		return ""
	}
	payload, _ := json.Marshal(v)
	return string(payload)
}

func parseScopeRef(ref string) (string, uint64, error) {
	parts := strings.Split(strings.TrimSpace(ref), ":")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid scope ref")
	}
	id, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return parts[0], id, nil
}

func supersededVersionID(versions []domain.TemplateVersion, supersedesVersion string) *uint64 {
	for i := range versions {
		if versions[i].Version == strings.TrimSpace(supersedesVersion) {
			id := versions[i].ID
			return &id
		}
	}
	return nil
}

func isBlockedCompatibility(items []domain.CompatibilityStatement) []string {
	reasons := make([]string, 0)
	for _, item := range items {
		if item.Result == domain.CompatibilityResultBlocked {
			reasons = append(reasons, item.Summary)
		}
	}
	return reasons
}

func ensureDBReady(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return err
}
