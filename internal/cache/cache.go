package cache

import (
	"context"

	"github.com/google/uuid"

	"task-management/internal/task/model"
)

type Cache interface {
	GetTaskList(ctx context.Context) ([]*model.Task, bool)
	SetTaskList(ctx context.Context, tasks []*model.Task)
	DeleteTaskList(ctx context.Context)

	GetTaskByID(ctx context.Context, id uuid.UUID) (*model.Task, bool)
	SetTaskByID(ctx context.Context, task *model.Task)
	DeleteTaskByID(ctx context.Context, id uuid.UUID)
}