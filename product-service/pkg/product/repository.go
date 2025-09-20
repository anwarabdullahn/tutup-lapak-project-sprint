package product

import (
	"context"
	"errors"
	"product-service/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, p *entities.Product) error
	GetByID(ctx context.Context, id string) (*entities.Product, error)
	List(ctx context.Context) ([]entities.Product, error)
	Update(ctx context.Context, p *entities.Product) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *entities.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) GetByID(ctx context.Context, id string) (*entities.Product, error) {
	var prod entities.Product
	if err := r.db.WithContext(ctx).First(&prod, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prod, nil
}

func (r *repository) List(ctx context.Context) ([]entities.Product, error) {
	var list []entities.Product
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, p *entities.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Product{}, "id = ?", id).Error
}
