package userfile

import "profile-service/pkg/entities"

// Service is an interface from which our api module can access our repository of all our models
type Service interface {
	UploadUserFile(file *entities.File) (*entities.File, error)
	GetUserFile(ID uint) (*entities.File, error)
	isFileExist(id string) (bool, error)
}

type service struct {
	repository Repository
}

// NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) UploadUserFile(file *entities.File) (*entities.File, error) {
	return s.repository.UploadUserFile(file)
}

func (s *service) GetUserFile(userID uint) (*entities.File, error) {
	return s.repository.GetUserFile(userID)
}

func (s *service) isFileExist(id string) (bool, error) {
	return s.repository.IsFileExist(id)
}
