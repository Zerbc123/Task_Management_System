package services

import (
	"task-management/internal/task/dto"
	"task-management/internal/task/model"
	"task-management/internal/task/repository"
	"github.com/google/uuid"
	"time"
)

type CategoryService interface {
	CreateCategory(req dto.CreateCategoryRequest) (*model.Category, error)
	GetCategories() ([]*model.Category, error)
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	UpdateCategory(id uuid.UUID, req dto.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(id uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req dto.CreateCategoryRequest) (*model.Category, error) {
	now := time.Now()
	category := &model.Category{
		ID: 		uuid.New(),
		Name: 		req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetCategories() ([]*model.Category, error) {
	return s.repo.List()
}

func (s *categoryService) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) UpdateCategory(id uuid.UUID, req dto.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	category.UpdatedAt = time.Now()

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	return s.repo.Delete(id)
}