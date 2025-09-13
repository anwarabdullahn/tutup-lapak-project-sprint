package userfile

import (
	"errors"
	"profile-service/pkg/entities"

	"gorm.io/gorm"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	UploadUserFile(file *entities.File) (*entities.File, error)
	GetUserFile(ID uint) (*entities.File, error)
	IsFileExist(id string) (bool, error)
}

type repository struct {
	DB *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) UploadUserFile(file *entities.File) (*entities.File, error) {
	if err := r.DB.Save(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (r *repository) GetUserFile(userID uint) (*entities.File, error) {
	var userFile entities.File
	err := r.DB.Where("user_id = ?", userID).First(&userFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userFile, nil
}

func (r *repository) IsFileExist(id string) (bool, error) {

	var file entities.File
	err := r.DB.Where("id = ?", id).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
