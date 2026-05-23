package repository

import (
	"errors"

	"github.com/google/uuid"
	"task-management/internal/project/model"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectRepository interface {
	Create(project *model.Project) error
	GetAll() ([]*model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Update(project *model.Project) error
	Delete(id uuid.UUID) error
}