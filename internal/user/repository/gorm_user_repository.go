package repository

import (
	"gorm.io/gorm"
	"task-management/internal/user/model"
	"github.com/google/uuid"
	"errors"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (r *gormUserRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error

	if err != nil{
		return err
	}

	return nil
}

func (r *gormUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil{
		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil{
		return nil, err
	}
	return &user, nil
}