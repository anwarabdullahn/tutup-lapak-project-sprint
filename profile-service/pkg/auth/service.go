package auth

import (
	"profile-service/pkg/entities"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(user *entities.User) (*entities.User, error)
	Login(user *entities.User) (*entities.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(user *entities.User) (*entities.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) Login(req *entities.User) (*entities.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, err
	}
	return user, nil
}
