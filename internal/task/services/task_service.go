package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"task-management/internal/cache"
	"task-management/internal/notification"
	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	"task-management/internal/task/repository"
	websocketHub "task-management/internal/websocket"
)

type TaskService interface {
	Create(ctx context.Context, req dto.CreateTaskRequest) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskService struct {
	repo  repository.TaskRepository
	cache cache.Cache
	queue notification.JobQueue
	hub   *websocketHub.Hub
}

func NewTaskService(
	repo repository.TaskRepository,
	cache cache.Cache,
	queue notification.JobQueue,
	hub *websocketHub.Hub,
) TaskService {
	return &taskService{
		repo:  repo,
		cache: cache,
		queue: queue,
		hub:   hub,
	}
}

func (s *taskService) Create(ctx context.Context, req dto.CreateTaskRequest) (*model.Task, error) {
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

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	s.cache.DeleteTaskList(ctx)

	if task.AssigneeID != nil {
		if err := s.queue.Push(ctx, notification.NotificationJob{
			TaskID:    task.ID,
			UserID:    task.AssigneeID,
			EventType: "task.created",
			Message:   "A new task has been assigned to you",
			MaxRetry:  3,
		}); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *taskService) GetAll(ctx context.Context) ([]*model.Task, error) {
	if tasks, ok := s.cache.GetTaskList(ctx); ok {
		return tasks, nil
	}

	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	s.cache.SetTaskList(ctx, tasks)

	return tasks, nil
}

func (s *taskService) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	if task, ok := s.cache.GetTaskByID(ctx, id); ok {
		return task, nil
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.cache.SetTaskByID(ctx, task)

	return task, nil
}

func (s *taskService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ProjectID != uuid.Nil {
		task.ProjectID = req.ProjectID
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	task.Description = req.Description

	if req.Status != "" {
		task.Status = req.Status
	}

	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	s.cache.DeleteTaskList(ctx)
	s.cache.DeleteTaskByID(ctx, id)

	if s.hub != nil {
		s.hub.BroadcastJSON(websocketHub.TaskUpdatedEvent{
			Type:   "task.updated",
			TaskID: task.ID,
			Title:  task.Title,
			Status: string(task.Status),
		})
	}

	if task.AssigneeID != nil {
		if err := s.queue.Push(ctx, notification.NotificationJob{
			TaskID:    task.ID,
			UserID:    task.AssigneeID,
			EventType: "task.updated",
			Message:   "Task has been updated",
			MaxRetry:  3,
		}); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.cache.DeleteTaskList(ctx)
	s.cache.DeleteTaskByID(ctx, id)

	return nil
}
