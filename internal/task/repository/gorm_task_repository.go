package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"task-management/internal/task/model"
)

type gormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) TaskRepository {
	return &gormTaskRepository{
		db: db,
	}
}

func (r *gormTaskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *gormTaskRepository) GetAll(ctx context.Context) ([]*model.Task, error) {
	var tasks []*model.Task

	err := r.db.WithContext(ctx).
		Find(&tasks).Error

	return tasks, err
}

func (r *gormTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := r.db.WithContext(ctx).
		First(&task, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}

	return &task, nil
}

func (r *gormTaskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *gormTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Delete(&model.Task{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return result.Error
}
