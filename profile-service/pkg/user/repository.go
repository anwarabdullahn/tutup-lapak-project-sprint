package user

import (
	"context"
	"errors"
	"profile-service/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByPhone(ctx context.Context, phone string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	GetByID(id string) (*entities.User, error)
	UpdateProfile(user *entities.User) (*entities.User, error)
	IsFileExist(id string) (bool, error)
	UpdateEmail(user *entities.User) (*entities.User, error)
	UpdatePhone(user *entities.User) (*entities.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) FindByPhone(ctx context.Context, phone string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) UpdateProfile(user *entities.User) (*entities.User, error) {
	if r.db == nil {
		return nil, errors.New("db is nil")
	}
	if user == nil {
		return nil, errors.New("user is nil")
	}

	// update hanya field yang dikirim
	if err := r.db.Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error; err != nil {
		return nil, err
	}

	// reload data setelah update (supaya struct terisi lengkap, termasuk relasi kalau preload)
	if err := r.db.Preload("File").
		First(user, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *GormRepository) UpdateEmail(user *entities.User) (*entities.User, error) {

	if r.db == nil {
		return nil, errors.New("db is nil")
	}
	if user == nil {
		return nil, errors.New("user is nil")
	}

	if err := r.db.Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"email": user.Email,
		}).Error; err != nil {

		return nil, err
	}
	if err := r.db.Preload("File").
		First(user, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) UpdatePhone(user *entities.User) (*entities.User, error) {

	if err := r.db.Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("File").
		First(user, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) GetByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("File").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormRepository) IsFileExist(id string) (bool, error) {

	var file entities.File
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
