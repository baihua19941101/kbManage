package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

type ObservabilityAccessSeed struct {
	WorkspaceID uint64
	ProjectID   uint64
	ClusterID   uint64
}

func SeedObservabilityAccess(t *testing.T, db *gorm.DB, userID uint64, namePrefix string, roleKey string) ObservabilityAccessSeed {
	t.Helper()
	if db == nil {
		t.Fatal("seed observability access requires non-nil db")
	}
	if userID == 0 {
		t.Fatal("seed observability access requires userID")
	}

	ctx := context.Background()
	suffix := fmt.Sprintf("%s-%d", namePrefix, time.Now().UnixNano())

	workspace := &domain.Workspace{Name: "ws-" + suffix, Description: "seeded workspace"}
	if err := db.WithContext(ctx).Create(workspace).Error; err != nil {
		t.Fatalf("seed workspace failed: %v", err)
	}

	project := &domain.Project{
		WorkspaceID: workspace.ID,
		Name:        "proj-" + suffix,
		Description: "seeded project",
	}
	if err := db.WithContext(ctx).Create(project).Error; err != nil {
		t.Fatalf("seed project failed: %v", err)
	}

	cluster := &domain.Cluster{
		Name:      "cluster-" + suffix,
		APIServer: "https://cluster-" + suffix + ".example.test",
		Status:    domain.ClusterStatusHealthy,
	}
	if err := db.WithContext(ctx).Create(cluster).Error; err != nil {
		t.Fatalf("seed cluster failed: %v", err)
	}

	workspaceBinding := &repository.WorkspaceClusterBinding{
		WorkspaceID: workspace.ID,
		ClusterID:   cluster.ID,
	}
	if err := db.WithContext(ctx).Create(workspaceBinding).Error; err != nil {
		t.Fatalf("seed workspace-cluster binding failed: %v", err)
	}

	projectBinding := &repository.ProjectClusterBinding{
		ProjectID: project.ID,
		ClusterID: cluster.ID,
	}
	if err := db.WithContext(ctx).Create(projectBinding).Error; err != nil {
		t.Fatalf("seed project-cluster binding failed: %v", err)
	}

	if roleKey == "" {
		roleKey = "workspace-owner"
	}
	var role repository.ScopeRole
	if err := db.WithContext(ctx).
		Where("scope_type = ? AND role_key = ?", "workspace", roleKey).
		First(&role).Error; err != nil {
		t.Fatalf("load scope role failed: %v", err)
	}

	binding := &repository.ScopeRoleBinding{
		SubjectType: "user",
		SubjectID:   userID,
		ScopeType:   "workspace",
		ScopeID:     workspace.ID,
		ScopeRoleID: role.ID,
		GrantedBy:   userID,
	}
	if err := db.WithContext(ctx).Create(binding).Error; err != nil {
		t.Fatalf("seed scope role binding failed: %v", err)
	}

	return ObservabilityAccessSeed{
		WorkspaceID: workspace.ID,
		ProjectID:   project.ID,
		ClusterID:   cluster.ID,
	}
}
