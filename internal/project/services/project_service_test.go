package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"task-management/internal/project/dto"
	"task-management/internal/project/model"
)

var errProjectNotFound = errors.New("project not found")

type mockProjectRepository struct {
	projects map[uuid.UUID]*model.Project
}

func newMockProjectRepository() *mockProjectRepository {
	return &mockProjectRepository{
		projects: make(map[uuid.UUID]*model.Project),
	}
}

func (m *mockProjectRepository) Create(project *model.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) GetAll() ([]*model.Project, error) {
	projects := make([]*model.Project, 0)

	for _, project := range m.projects {
		projects = append(projects, project)
	}

	return projects, nil
}

func (m *mockProjectRepository) GetByID(id uuid.UUID) (*model.Project, error) {
	project, exists := m.projects[id]
	if !exists {
		return nil, errProjectNotFound
	}

	return project, nil
}

func (m *mockProjectRepository) Update(project *model.Project) error {
	if _, exists := m.projects[project.ID]; !exists {
		return errProjectNotFound
	}

	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) Delete(id uuid.UUID) error {
	if _, exists := m.projects[id]; !exists {
		return errProjectNotFound
	}

	delete(m.projects, id)
	return nil
}

func setupProjectService() ProjectService {
	repo := newMockProjectRepository()
	return NewProjectService(repo)
}

func TestCreateProject(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()

	project, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Task Management",
		Description: "Backend project",
	}, ownerID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if project.ID == uuid.Nil {
		t.Fatal("expected project ID to be generated")
	}

	if project.Name != "Task Management" {
		t.Fatalf("expected project name Task Management, got %s", project.Name)
	}

	if project.OwnerID != ownerID {
		t.Fatalf("expected owner ID %s, got %s", ownerID, project.OwnerID)
	}
}

func TestGetProjects(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()

	_, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Project 1",
		Description: "Description 1",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	projects, err := service.GetProjects()
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
}

func TestGetProjectByID(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()

	createdProject, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Project detail",
		Description: "Get project by ID",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	project, err := service.GetProjectByID(createdProject.ID)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID != createdProject.ID {
		t.Fatalf("expected project ID %s, got %s", createdProject.ID, project.ID)
	}
}

func TestUpdateProject(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()

	project, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Old project",
		Description: "Old description",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	updatedProject, err := service.UpdateProject(
		project.ID,
		ownerID,
		dto.UpdateProjectRequest{
			Name:        "New project",
			Description: "New description",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if updatedProject.Name != "New project" {
		t.Fatalf("expected New project, got %s", updatedProject.Name)
	}

	if updatedProject.Description != "New description" {
		t.Fatalf("expected New description, got %s", updatedProject.Description)
	}
}

func TestUpdateProjectForbidden(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()
	anotherUserID := uuid.New()

	project, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Owner project",
		Description: "Only owner can update",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.UpdateProject(
		project.ID,
		anotherUserID,
		dto.UpdateProjectRequest{
			Name:        "Hacked project",
			Description: "Should not update",
		},
	)

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteProject(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()

	project, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Project to delete",
		Description: "Delete this project",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteProject(project.ID, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.GetProjectByID(project.ID)
	if err == nil {
		t.Fatal("expected error after deleting project")
	}
}

func TestDeleteProjectForbidden(t *testing.T) {
	service := setupProjectService()

	ownerID := uuid.New()
	anotherUserID := uuid.New()

	project, err := service.CreateProject(dto.CreateProjectRequest{
		Name:        "Owner project",
		Description: "Only owner can delete",
	}, ownerID)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteProject(project.ID, anotherUserID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
