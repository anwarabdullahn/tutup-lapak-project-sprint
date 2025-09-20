package product

import (
	"context"
	"errors"
	"strings"

	"product-service/pkg/entities"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, in CreateInput) (*entities.Product, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	List(ctx context.Context) ([]entities.Product, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

type CreateInput struct {
	Name        string
	Description string
	Category    string
	PriceCents  int64
}

type UpdateInput struct {
	Name        *string
	Description *string
	Category    *string
	PriceCents  *int64
}

var allowedCategories = map[string]struct{}{
	"Food": {}, "Beverage": {}, "Clothes": {}, "Furniture": {}, "Tools": {},
}

func validateCategory(cat string) bool {
	_, ok := allowedCategories[cat]
	return ok
}

func (s *service) Create(ctx context.Context, in CreateInput) (*entities.Product, error) {
	in.Category = strings.TrimSpace(in.Category)
	if !validateCategory(in.Category) {
		return nil, errors.New("invalid category")
	}
	p := &entities.Product{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
		Category:    in.Category,
		PriceCents:  in.PriceCents,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	return s.repo.GetByID(ctx, id.String())
}

func (s *service) List(ctx context.Context) ([]entities.Product, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*entities.Product, error) {
	p, err := s.repo.GetByID(ctx, id.String())
	if err != nil || p == nil {
		return p, err
	}
	if in.Name != nil {
		p.Name = strings.TrimSpace(*in.Name)
	}
	if in.Description != nil {
		p.Description = *in.Description
	}
	if in.Category != nil {
		cat := strings.TrimSpace(*in.Category)
		if !validateCategory(cat) {
			return nil, errors.New("invalid category")
		}
		p.Category = cat
	}
	if in.PriceCents != nil {
		p.PriceCents = *in.PriceCents
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id.String())
}
