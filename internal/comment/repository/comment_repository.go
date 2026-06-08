package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"task-management/internal/comment/model"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]*model.Comment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}