package workspace

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

var errWorkspaceNameRequired = errors.New("workspace name is required")

type CreateWorkspaceRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Service struct {
	repo *repository.WorkspaceRepository
}

func NewService(repo *repository.WorkspaceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateWorkspaceRequest) (*domain.Workspace, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errWorkspaceNameRequired
	}

	item := &domain.Workspace{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Workspace, error) {
	return s.repo.List(ctx)
}
