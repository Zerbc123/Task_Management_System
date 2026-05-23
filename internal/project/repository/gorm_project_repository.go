package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"task-management/internal/project/model"
)

type gormProjectRepository struct {
	db *gorm.DB
}

func NewGormProjectRepository(db *gorm.DB) ProjectRepository {
	return &gormProjectRepository{
		db: db,
	}
}

func (r *gormProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *gormProjectRepository) GetAll() ([]*model.Project, error) {
	var projects []*model.Project

	err := r.db.Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *gormProjectRepository) GetByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project

	err := r.db.First(&project, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProjectNotFound
	}

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *gormProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *gormProjectRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&model.Project{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}

	return nil
}