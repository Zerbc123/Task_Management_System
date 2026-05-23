package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"task-management/internal/auth/dto"
	"task-management/internal/shared/utils"
	userModel "task-management/internal/user/model"
	userRepository "task-management/internal/user/repository"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService interface {
	Register(req dto.RegisterRequest) (*userModel.User, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo userRepository.UserRepository
}

func NewAuthService(userRepo userRepository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*userModel.User, error) {
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, userRepository.ErrEmailAlreadyExists
	}

	if !errors.Is(err, userRepository.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := &userModel.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}