package project

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

var (
	errWorkspaceIDRequired = errors.New("workspace id is required")
	errProjectNameRequired = errors.New("project name is required")
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Service struct {
	projectRepo   *repository.ProjectRepository
	workspaceRepo *repository.WorkspaceRepository
}

func NewService(projectRepo *repository.ProjectRepository, workspaceRepo *repository.WorkspaceRepository) *Service {
	return &Service{projectRepo: projectRepo, workspaceRepo: workspaceRepo}
}

func (s *Service) Create(ctx context.Context, workspaceID uint64, req CreateProjectRequest) (*domain.Project, error) {
	if workspaceID == 0 {
		return nil, errWorkspaceIDRequired
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errProjectNameRequired
	}

	if _, err := s.workspaceRepo.GetByID(ctx, workspaceID); err != nil {
		return nil, err
	}

	item := &domain.Project{
		WorkspaceID: workspaceID,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	}
	if err := s.projectRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]domain.Project, error) {
	if workspaceID == 0 {
		return nil, errWorkspaceIDRequired
	}
	if _, err := s.workspaceRepo.GetByID(ctx, workspaceID); err != nil {
		return nil, err
	}
	return s.projectRepo.ListByWorkspace(ctx, workspaceID)
}
