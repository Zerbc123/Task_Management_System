package services

import (
	"task-management/internal/auth/dto"
	"task-management/internal/shared/utils"
	userModel "task-management/internal/user/model"
	userRepository "task-management/internal/user/repository"
	"testing"

	"github.com/google/uuid"
)

type mockUserRepository struct {
	users map[string]*userModel.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*userModel.User),
	}
}

func (m *mockUserRepository) Create(user *userModel.User) error {
	if _, exists := m.users[user.Email]; exists {
		return userRepository.ErrEmailAlreadyExists
	}

	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*userModel.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, userRepository.ErrUserNotFound
	}

	return user, nil
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*userModel.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, userRepository.ErrUserNotFound
}

func TestRegisterSuccess(t *testing.T) {
	repo := newMockUserRepository()
	service := NewAuthService(repo)

	user, err := service.Register(dto.RegisterRequest{
		Email:    "ngoc@gmail.com",
		Password: "123456",
		FullName: "Lam Ngoc",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == uuid.Nil {
		t.Fatal("expected user id to be generated")
	}

	if user.Email != "ngoc@gmail.com" {
		t.Fatalf("expected email ngoc@gmail.com, got %s", user.Email)
	}

	if user.FullName != "Lam Ngoc" {
		t.Fatalf("expected full name Lam Ngoc, got %s", user.FullName)
	}

	if user.PasswordHash == "123456" {
		t.Fatal("expected password to be hashed")
	}
}

func TestRegisterEmailAlreadyExists(t *testing.T) {
	repo := newMockUserRepository()
	service := NewAuthService(repo)

	_, err := service.Register(dto.RegisterRequest{
		Email:    "ngoc@gmail.com",
		Password: "123456",
		FullName: "Lam Ngoc",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Register(dto.RegisterRequest{
		Email:    "ngoc@gmail.com",
		Password: "abcdef",
		FullName: "Be Ngoc",
	})

	if err != userRepository.ErrEmailAlreadyExists {
		t.Fatalf("expected EmailAlreadyExists, got %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newMockUserRepository()
	service := NewAuthService(repo)

	passwordHash, err := utils.HashPassword("123456")
	if err != nil {
		t.Fatal(err)
	}

	user := &userModel.User{
		ID:           uuid.New(),
		Email:        "ngoc@gmail.com",
		PasswordHash: passwordHash,
		FullName:     "ngoc",
	}

	if err := repo.Create(user); err != nil {
		t.Fatal(err)
	}

	res, err := service.Login(dto.LoginRequest{
		Email: "ngoc@gmail.com",
		Password: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Token == "" {
		t.Fatal("expected JWT token")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repo := newMockUserRepository()
	service := NewAuthService(repo)
	
	passwordHash, err := utils.HashPassword("123456")
	if err != nil {
		t.Fatal(err)
	}

	user := &userModel.User{
		ID:           uuid.New(),
		Email:        "ngoc@gmail.com",
		PasswordHash: passwordHash,
		FullName:     "ngoc",
	}

	if err := repo.Create(user); err != nil {
		t.Fatal(err)
	}

	_, err = service.Login(dto.LoginRequest{
		Email:    "duy@gmail.com",
		Password: "wrong-password",
	})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
