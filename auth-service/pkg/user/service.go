package user

import (
	"auth-service/pkg/entities"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req entities.CreateUserRequest) (*entities.User, error)
	Login(ctx context.Context, req entities.LoginRequest) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByPhone(ctx context.Context, phone string) (*entities.User, error)
	FindById(ctx context.Context, id string) (*entities.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req entities.CreateUserRequest) (*entities.User, error) {
	// Create new user entity from request
	user := &entities.User{
		Email: req.Email,
		Phone: req.Phone,
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashed)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password in response
	user.Password = ""
	return user, nil
}

func (s *service) Login(ctx context.Context, req entities.LoginRequest) (*entities.User, error) {
	var user *entities.User
	var err error

	// Login with email or phone
	if req.Email != "" {
		user, err = s.repo.FindByEmail(ctx, req.Email)
	} else if req.Phone != "" {
		user, err = s.repo.FindByPhone(ctx, req.Phone)
	} else {
		return nil, fmt.Errorf("email or phone is required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Don't return password in response
	user.Password = ""
	return user, nil
}

func (s *service) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	// Don't return password in response
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (s *service) FindByPhone(ctx context.Context, phone string) (*entities.User, error) {
	user, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}
	// Don't return password in response
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (s *service) FindById(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	// Don't return password in response
	if user != nil {
		user.Password = ""
	}
	return user, nil
}
