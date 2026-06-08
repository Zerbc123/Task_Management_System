package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"task-management/internal/task/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
