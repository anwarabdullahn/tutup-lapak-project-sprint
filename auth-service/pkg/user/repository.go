package user

import (
	"auth-service/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByPhone(phone string) (*entities.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *GormRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) FindByPhone(phone string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
