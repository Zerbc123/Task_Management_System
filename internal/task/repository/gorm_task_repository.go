package repository

import (
	"gorm.io/gorm"
	"task-management/internal/task/model"
	"github.com/google/uuid"
	"errors"
)

type gormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) TaskRepository {
	return &gormTaskRepository{
		db: db,
	}
}

func (r *gormTaskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *gormTaskRepository) GetAll() ([]*model.Task, error) {
	var tasks []*model.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *gormTaskRepository) GetByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

func (r *gormTaskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *gormTaskRepository) Delete(id uuid.UUID) error {
	result :=  r.db.Delete(&model.Task{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return result.Error
}
