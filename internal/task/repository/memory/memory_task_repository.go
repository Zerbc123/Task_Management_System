package memory

import (
	"task-management/internal/task/model"
	"task-management/internal/task/repository"

	"github.com/google/uuid"
)

type memoryTaskRepository struct {
	tasks map[uuid.UUID]*model.Task
}

func NewMemoryTaskRepository() repository.TaskRepository {
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
		return nil, repository.ErrTaskNotFound
	}

	return task, nil
}

func (r *memoryTaskRepository) Update(task *model.Task) error {

	_, ok := r.tasks[task.ID]

	if !ok {
		return repository.ErrTaskNotFound
	}

	r.tasks[task.ID] = task

	return nil
}

func (r *memoryTaskRepository) Delete(id uuid.UUID) error {

	_, ok := r.tasks[id]

	if !ok {
		return repository.ErrTaskNotFound
	}

	delete(r.tasks, id)

	return nil
}