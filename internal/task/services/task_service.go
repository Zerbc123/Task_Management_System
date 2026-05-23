package services

import (
	"time"

	"github.com/google/uuid"

	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	"task-management/internal/task/repository"
)

type TaskService interface {
	CreateTask(req dto.CreateTaskRequest) (*model.Task, error)
	GetTasks() ([]*model.Task, error)
	GetTaskByID(id uuid.UUID) (*model.Task, error)
	UpdateTask(id uuid.UUID, req dto.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(id uuid.UUID) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) CreateTask(req dto.CreateTaskRequest) (*model.Task, error) {
	now := time.Now()

	status := req.Status
	if status == "" {
		status = model.StatusTodo
	}

	task := &model.Task{
		ID:          uuid.New(),
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		AssigneeID:  req.AssigneeID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTasks() ([]*model.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) GetTaskByID(id uuid.UUID) (*model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) UpdateTask(id uuid.UUID, req dto.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.ProjectID != uuid.Nil {
		task.ProjectID = req.ProjectID
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.AssigneeID = req.AssigneeID
	task.UpdatedAt = time.Now()

	if task.Status == "" {
		task.Status = model.StatusTodo
	}

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uuid.UUID) error {
	return s.repo.Delete(id)
}