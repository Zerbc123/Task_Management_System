package repository

import (
	"errors"

	"task-management/internal/task/model"

	"github.com/google/uuid"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(task *model.Task) error
	GetAll() ([]*model.Task, error)
	GetByID(id uuid.UUID) (*model.Task, error)
	Update(task *model.Task) error
	Delete(id uuid.UUID) error
}
