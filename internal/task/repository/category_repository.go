package repository

import (
	"errors"
	"github.com/google/uuid"
	"task-management/internal/task/model"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uuid.UUID) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	List() ([]*model.Category, error)
}

type memoryCategoryRepository struct {
	categories map[uuid.UUID]*model.Category
}

func NewMemoryCategoryRepository() CategoryRepository {
	return &memoryCategoryRepository{
		categories: make(map[uuid.UUID]*model.Category),
	}
}

func (r *memoryCategoryRepository) Create(category *model.Category) error {
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepository) List() ([]*model.Category, error) {
	categories := make([]*model.Category, 0)

	for _, category := range r.categories {
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *memoryCategoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	category, ok := r.categories[id]
	if !ok {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (r *memoryCategoryRepository) Update(category *model.Category) error {
	if _, ok := r.categories[category.ID]; !ok {
		return ErrCategoryNotFound
	}
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepository) Delete(id uuid.UUID) error {
	if _, ok := r.categories[id]; !ok {
		return ErrCategoryNotFound
	}
	delete(r.categories, id)
	return nil
}