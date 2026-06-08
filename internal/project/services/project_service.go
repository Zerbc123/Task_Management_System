package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"task-management/internal/project/dto"
	"task-management/internal/project/model"
	"task-management/internal/project/repository"
)

var ErrForbidden = errors.New("you do not have permission to access this project")

type ProjectService interface {
	CreateProject(req dto.CreateProjectRequest, ownerID uuid.UUID) (*model.Project, error)
	GetProjects() ([]*model.Project, error)
	GetProjectByID(id uuid.UUID) (*model.Project, error)
	UpdateProject(id uuid.UUID, ownerID uuid.UUID, req dto.UpdateProjectRequest) (*model.Project, error)
	DeleteProject(id uuid.UUID, ownerID uuid.UUID) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

func (s *projectService) CreateProject(req dto.CreateProjectRequest, ownerID uuid.UUID) (*model.Project, error) {
	now := time.Now()

	project := &model.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) GetProjects() ([]*model.Project, error) {
	return s.repo.GetAll()
}

func (s *projectService) GetProjectByID(id uuid.UUID) (*model.Project, error) {
	return s.repo.GetByID(id)
}

func (s *projectService) UpdateProject(
	id uuid.UUID,
	ownerID uuid.UUID,
	req dto.UpdateProjectRequest,
) (*model.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != ownerID {
		return nil, ErrForbidden
	}

	project.Name = req.Name
	project.Description = req.Description
	project.UpdatedAt = time.Now()

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) DeleteProject(id uuid.UUID, ownerID uuid.UUID) error {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if project.OwnerID != ownerID {
		return ErrForbidden
	}

	return s.repo.Delete(id)
}
