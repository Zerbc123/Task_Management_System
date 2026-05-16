package repository

import (
	"errors"

	"github.com/google/uuid"
	"task-management/internal/task/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(task *model.Task) error
	GetAll() ([]*model.Task, error)
	GetByID(id uuid.UUID) (*model.Task, error)
	Update(task *model.Task) error
	Delete(id uuid.UUID) error
}

type memoryTaskRepository struct {
	tasks map[uuid.UUID]*model.Task
}

func NewMemoryTaskRepository() TaskRepository {
	return &memoryTaskRepository{
		tasks: make(map[uuid.UUID]*model.Task),
	}
}

func (r *memoryTaskRepository) Create(task *model.Task) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *memoryTaskRepository) GetAll() ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)

	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *memoryTaskRepository) GetByID(id uuid.UUID) (*model.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (r *memoryTaskRepository) Update(task *model.Task) error {
	if _, ok := r.tasks[task.ID]; !ok {
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *memoryTaskRepository) Delete(id uuid.UUID) error {
	if _, ok := r.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}