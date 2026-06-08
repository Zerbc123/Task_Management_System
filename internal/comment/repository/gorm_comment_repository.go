package repository

import (
	"context"
	"errors"
	"task-management/internal/comment/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormCommentRepository struct {
	db *gorm.DB
}

func NewGormCommentRepository(db *gorm.DB) CommentRepository {
	return &gormCommentRepository{
		db: db,
	}
}

func (r *gormCommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *gormCommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment

	err := r.db.WithContext(ctx).First(&comment, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCommentNotFound
	}

	return &comment, err
}

func (r *gormCommentRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&comments).Error

	return comments, err
}

func (r *gormCommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *gormCommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Comment{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return ErrCommentNotFound
	}

	return result.Error
}
