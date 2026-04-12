package observability

import (
	"context"
	"errors"
	"sort"

	"kbmanage/backend/internal/domain"
	authSvc "kbmanage/backend/internal/service/auth"
)

var ErrInvalidObservabilityUser = errors.New("observability user id is required")

type AccessContext struct {
	UserID             uint64
	Permission         string
	ClusterConstrained bool
	WorkspaceIDs       []uint64
	ProjectIDs         []uint64
	ClusterIDs         []uint64
}

type accessContextKey struct{}

type ClusterScopeResolver interface {
	ResolveObservabilityScopeByClusterIDs(ctx context.Context, clusterIDs []uint64) ([]uint64, []uint64, error)
}

type ScopeService struct {
	authorizer      *authSvc.ScopeAuthorizer
	scopeAccess     *authSvc.ScopeAccessService
	clusterResolver ClusterScopeResolver
}

func NewScopeService(
	authorizer *authSvc.ScopeAuthorizer,
	scopeAccess *authSvc.ScopeAccessService,
	clusterResolver ClusterScopeResolver,
) *ScopeService {
	if authorizer == nil {
		authorizer = authSvc.NewScopeAuthorizer()
	}
	return &ScopeService{
		authorizer:      authorizer,
		scopeAccess:     scopeAccess,
		clusterResolver: clusterResolver,
	}
}

func WithAccessContext(ctx context.Context, access AccessContext) context.Context {
	access.ClusterIDs = dedupeUint64(access.ClusterIDs)
	access.WorkspaceIDs = dedupeUint64(access.WorkspaceIDs)
	access.ProjectIDs = dedupeUint64(access.ProjectIDs)
	return context.WithValue(ctx, accessContextKey{}, access)
}

func AccessContextFromContext(ctx context.Context) (AccessContext, bool) {
	if ctx == nil {
		return AccessContext{}, false
	}
	access, ok := ctx.Value(accessContextKey{}).(AccessContext)
	return access, ok
}

type ScopeFilter struct {
	ClusterIDs   []uint64
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Namespaces   []string
}

func (s *ScopeService) FilterByScope(ctx context.Context, userID uint64, filter ScopeFilter) (ScopeFilter, error) {
	if access, ok := AccessContextFromContext(ctx); ok && userID == 0 {
		userID = access.UserID
	}
	if userID == 0 {
		return ScopeFilter{}, ErrInvalidObservabilityUser
	}

	filter.ClusterIDs = dedupeUint64(filter.ClusterIDs)
	filter.WorkspaceIDs = dedupeUint64(filter.WorkspaceIDs)
	filter.ProjectIDs = dedupeUint64(filter.ProjectIDs)

	access, hasAccess := AccessContextFromContext(ctx)
	if hasAccess {
		if access.ClusterConstrained {
			if len(filter.ClusterIDs) == 0 {
				filter.ClusterIDs = append([]uint64(nil), access.ClusterIDs...)
			} else {
				filter.ClusterIDs = intersectUint64(filter.ClusterIDs, access.ClusterIDs)
			}
		}

		if len(filter.WorkspaceIDs) > 0 && len(access.WorkspaceIDs) > 0 {
			filter.WorkspaceIDs = intersectUint64(filter.WorkspaceIDs, access.WorkspaceIDs)
		}
		if len(filter.ProjectIDs) > 0 && len(access.ProjectIDs) > 0 {
			filter.ProjectIDs = intersectUint64(filter.ProjectIDs, access.ProjectIDs)
		}
	}

	if len(filter.ClusterIDs) > 0 && s != nil && s.clusterResolver != nil {
		mappedWorkspaceIDs, mappedProjectIDs, err := s.clusterResolver.ResolveObservabilityScopeByClusterIDs(ctx, filter.ClusterIDs)
		if err != nil {
			return ScopeFilter{}, err
		}
		filter.WorkspaceIDs = mergeUint64(filter.WorkspaceIDs, mappedWorkspaceIDs)
		filter.ProjectIDs = mergeUint64(filter.ProjectIDs, mappedProjectIDs)
	}

	if hasAccess && access.ClusterConstrained &&
		len(filter.ClusterIDs) == 0 &&
		len(filter.WorkspaceIDs) == 0 &&
		len(filter.ProjectIDs) == 0 {
		return ScopeFilter{}, ErrObservabilityScopeDenied
	}
	if len(filter.WorkspaceIDs) == 0 && len(filter.ProjectIDs) == 0 && len(filter.ClusterIDs) == 0 {
		return ScopeFilter{}, ErrObservabilityScopeDenied
	}

	return filter, nil
}

func (s *ScopeService) CanAccessScope(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	target ScopeFilter,
) bool {
	if s == nil || s.authorizer == nil {
		return false
	}
	targetWorkspaceIDs := dedupeUint64(target.WorkspaceIDs)
	targetProjectIDs := dedupeUint64(target.ProjectIDs)
	clusterScoped := len(target.ClusterIDs) > 0
	return s.authorizer.CanAccessObservabilityMapped(
		grantedType,
		grantedWorkspaceID,
		grantedProjectID,
		targetWorkspaceIDs,
		targetProjectIDs,
		clusterScoped,
	)
}

func dedupeUint64(values []uint64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}
	set := make(map[uint64]struct{}, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		set[value] = struct{}{}
	}
	out := make([]uint64, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func mergeUint64(left, right []uint64) []uint64 {
	out := append([]uint64{}, left...)
	out = append(out, right...)
	return dedupeUint64(out)
}

func intersectUint64(left, right []uint64) []uint64 {
	if len(left) == 0 || len(right) == 0 {
		return []uint64{}
	}
	set := make(map[uint64]struct{}, len(right))
	for _, item := range right {
		if item == 0 {
			continue
		}
		set[item] = struct{}{}
	}
	out := make([]uint64, 0, len(left))
	for _, item := range left {
		if item == 0 {
			continue
		}
		if _, ok := set[item]; ok {
			out = append(out, item)
		}
	}
	return dedupeUint64(out)
}
