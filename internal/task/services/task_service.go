package services

import (
	"time"

	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	"task-management/internal/task/repository"

	"github.com/google/uuid"
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
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := s.repo.Create(task)
	if err != nil {
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

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.Priority = req.Priority
	task.AssignedTo = req.AssignedTo
	task.UpdatedAt = time.Now()

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uuid.UUID) error {
	return s.repo.Delete(id)
}