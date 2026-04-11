package auth

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

const scopeBindingPageSize = 200

type ScopeAccessService struct {
	bindingRepo          *repository.ScopeRoleBindingRepository
	projectRepo          *repository.ProjectRepository
	workspaceClusterRepo *repository.WorkspaceClusterRepository
	scopeAuthorizer      *ScopeAuthorizer
	permissionSvc        *PermissionService
}

func NewScopeAccessService(
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	scopeAuthorizer *ScopeAuthorizer,
	permissionSvc *PermissionService,
	workspaceClusterRepo ...*repository.WorkspaceClusterRepository,
) *ScopeAccessService {
	var clusterRepo *repository.WorkspaceClusterRepository
	if len(workspaceClusterRepo) > 0 {
		clusterRepo = workspaceClusterRepo[0]
	}
	if scopeAuthorizer == nil {
		scopeAuthorizer = NewScopeAuthorizer()
	}
	if permissionSvc == nil {
		permissionSvc = NewPermissionService()
	}
	return &ScopeAccessService{
		bindingRepo:          bindingRepo,
		projectRepo:          projectRepo,
		workspaceClusterRepo: clusterRepo,
		scopeAuthorizer:      scopeAuthorizer,
		permissionSvc:        permissionSvc,
	}
}

func (s *ScopeAccessService) HasScopePermission(
	ctx context.Context,
	userID uint64,
	targetType domain.ScopeType,
	targetWorkspaceID uint64,
	targetProjectID uint64,
	permission string,
) (bool, error) {
	if userID == 0 {
		return false, nil
	}
	if s == nil || s.bindingRepo == nil {
		return false, gorm.ErrInvalidDB
	}

	workspaceID, projectID, err := s.resolveTargetScope(ctx, targetType, targetWorkspaceID, targetProjectID)
	if err != nil {
		return false, err
	}

	roleKeys, err := s.matchingRoleKeys(ctx, userID, targetType, workspaceID, projectID)
	if err != nil {
		return false, err
	}
	if len(roleKeys) == 0 {
		return false, nil
	}

	required := strings.TrimSpace(permission)
	if required == "" {
		return true, nil
	}
	return s.permissionSvc.HasPermission(roleKeys, required), nil
}

func (s *ScopeAccessService) ListWorkspaceIDsByPermission(ctx context.Context, userID uint64, permission string) ([]uint64, error) {
	if userID == 0 {
		return []uint64{}, nil
	}
	if s == nil || s.bindingRepo == nil {
		return nil, gorm.ErrInvalidDB
	}

	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.listWorkspaceIDsByPermissionFromBindings(ctx, bindings, permission)
}

func (s *ScopeAccessService) ListClusterIDsByPermission(ctx context.Context, userID uint64, permission string) ([]uint64, bool, error) {
	if userID == 0 {
		return []uint64{}, true, nil
	}
	if s == nil || s.bindingRepo == nil {
		return nil, false, gorm.ErrInvalidDB
	}
	if s.workspaceClusterRepo == nil {
		return []uint64{}, false, nil
	}

	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	if len(bindings) == 0 {
		return []uint64{}, false, nil
	}

	workspaceIDs, err := s.listWorkspaceIDsByPermissionFromBindings(ctx, bindings, permission)
	if err != nil {
		return nil, false, err
	}
	if len(workspaceIDs) == 0 {
		return []uint64{}, true, nil
	}

	workspaceClusterItems, err := s.workspaceClusterRepo.ListByWorkspaces(ctx, workspaceIDs)
	if err != nil {
		return nil, false, err
	}

	allowedClusterIDs := make(map[uint64]struct{}, len(workspaceClusterItems))
	for _, item := range workspaceClusterItems {
		if item.ClusterID == 0 {
			continue
		}
		allowedClusterIDs[item.ClusterID] = struct{}{}
	}
	return setToSortedIDs(allowedClusterIDs), true, nil
}

func (s *ScopeAccessService) CanAccessClusterByPermission(ctx context.Context, userID uint64, clusterID uint64, permission string) (bool, error) {
	if userID == 0 || clusterID == 0 {
		return false, nil
	}
	clusterIDs, constrained, err := s.ListClusterIDsByPermission(ctx, userID, permission)
	if err != nil {
		return false, err
	}
	if !constrained {
		return true, nil
	}
	for _, id := range clusterIDs {
		if id == clusterID {
			return true, nil
		}
	}
	return false, nil
}

func (s *ScopeAccessService) resolveTargetScope(
	ctx context.Context,
	targetType domain.ScopeType,
	workspaceID uint64,
	projectID uint64,
) (uint64, uint64, error) {
	if targetType != domain.ScopeTypeProject {
		return workspaceID, projectID, nil
	}
	if workspaceID != 0 || projectID == 0 || s.projectRepo == nil {
		return workspaceID, projectID, nil
	}

	item, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workspaceID, projectID, nil
		}
		return 0, 0, err
	}
	return item.WorkspaceID, projectID, nil
}

func (s *ScopeAccessService) matchingRoleKeys(
	ctx context.Context,
	userID uint64,
	targetType domain.ScopeType,
	targetWorkspaceID uint64,
	targetProjectID uint64,
) ([]string, error) {
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return nil, err
	}

	projectWorkspaceCache := make(map[uint64]uint64)
	allowed := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		grantedType := parseScopeType(binding.ScopeType)
		if grantedType == "" {
			continue
		}
		grantedWorkspaceID, grantedProjectID, err := s.resolveBindingScope(ctx, grantedType, binding.ScopeID, projectWorkspaceCache)
		if err != nil {
			return nil, err
		}

		if !s.scopeAuthorizer.CanAccess(
			grantedType,
			grantedWorkspaceID,
			grantedProjectID,
			targetType,
			targetWorkspaceID,
			targetProjectID,
		) {
			continue
		}

		roleKey := strings.TrimSpace(binding.RoleKey)
		if roleKey == "" {
			continue
		}
		allowed[roleKey] = struct{}{}
	}

	result := make([]string, 0, len(allowed))
	for roleKey := range allowed {
		result = append(result, roleKey)
	}
	return result, nil
}

func (s *ScopeAccessService) listUserBindings(ctx context.Context, userID uint64) ([]repository.ScopeRoleBindingWithRole, error) {
	items := make([]repository.ScopeRoleBindingWithRole, 0, scopeBindingPageSize)
	offset := 0
	for {
		batch, err := s.bindingRepo.List(ctx, repository.ScopeRoleBindingFilter{
			SubjectType: "user",
			SubjectID:   userID,
			Limit:       scopeBindingPageSize,
			Offset:      offset,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, batch...)
		if len(batch) < scopeBindingPageSize {
			break
		}
		offset += scopeBindingPageSize
	}
	return items, nil
}

func (s *ScopeAccessService) listWorkspaceIDsByPermissionFromBindings(
	ctx context.Context,
	bindings []repository.ScopeRoleBindingWithRole,
	permission string,
) ([]uint64, error) {
	projectWorkspaceCache := make(map[uint64]uint64)
	allowed := make(map[uint64]struct{}, len(bindings))
	required := strings.TrimSpace(permission)

	for _, binding := range bindings {
		grantedType := parseScopeType(binding.ScopeType)
		if grantedType != domain.ScopeTypeWorkspace && grantedType != domain.ScopeTypeProject {
			continue
		}
		roleKey := strings.TrimSpace(binding.RoleKey)
		if roleKey == "" {
			continue
		}
		if required != "" && !s.permissionSvc.HasPermission([]string{roleKey}, required) {
			continue
		}

		workspaceID, projectID, err := s.resolveBindingScope(ctx, grantedType, binding.ScopeID, projectWorkspaceCache)
		if err != nil {
			return nil, err
		}
		if workspaceID == 0 {
			continue
		}

		if !s.scopeAuthorizer.CanAccess(
			grantedType,
			workspaceID,
			projectID,
			domain.ScopeTypeWorkspace,
			workspaceID,
			0,
		) {
			continue
		}
		allowed[workspaceID] = struct{}{}
	}

	return setToSortedIDs(allowed), nil
}

func (s *ScopeAccessService) resolveBindingScope(
	ctx context.Context,
	grantedType domain.ScopeType,
	scopeID uint64,
	projectWorkspaceCache map[uint64]uint64,
) (uint64, uint64, error) {
	switch grantedType {
	case domain.ScopeTypeWorkspace:
		return scopeID, 0, nil
	case domain.ScopeTypeProject:
		if scopeID == 0 {
			return 0, 0, nil
		}
		if cached, ok := projectWorkspaceCache[scopeID]; ok {
			return cached, scopeID, nil
		}
		if s.projectRepo == nil {
			projectWorkspaceCache[scopeID] = 0
			return 0, scopeID, nil
		}

		item, err := s.projectRepo.GetByID(ctx, scopeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				projectWorkspaceCache[scopeID] = 0
				return 0, scopeID, nil
			}
			return 0, 0, err
		}
		projectWorkspaceCache[scopeID] = item.WorkspaceID
		return item.WorkspaceID, scopeID, nil
	default:
		return 0, 0, nil
	}
}

func setToSortedIDs(values map[uint64]struct{}) []uint64 {
	ids := make([]uint64, 0, len(values))
	for id := range values {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func ParseClusterIDFromReference(raw string) (uint64, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}

	for _, segment := range strings.Split(trimmed, "/") {
		part := strings.TrimSpace(segment)
		if !strings.HasPrefix(strings.ToLower(part), "cluster:") {
			continue
		}
		idText := strings.TrimSpace(part[len("cluster:"):])
		id, err := strconv.ParseUint(idText, 10, 64)
		if err != nil || id == 0 {
			return 0, false
		}
		return id, true
	}
	return 0, false
}

func parseScopeType(scopeType string) domain.ScopeType {
	switch strings.ToLower(strings.TrimSpace(scopeType)) {
	case string(domain.ScopeTypePlatform):
		return domain.ScopeTypePlatform
	case string(domain.ScopeTypeWorkspace):
		return domain.ScopeTypeWorkspace
	case string(domain.ScopeTypeProject):
		return domain.ScopeTypeProject
	default:
		return ""
	}
}
