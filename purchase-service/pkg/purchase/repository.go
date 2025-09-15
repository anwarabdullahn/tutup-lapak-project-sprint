package purchase

import (
	"context"
	"purchase-service/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) error
	CreatePurchaseItems(ctx context.Context, items []*entities.PurchaseItem) error
	CreatePurchaseSender(ctx context.Context, sender *entities.PurchaseSender) error
	GetPurchaseByID(ctx context.Context, id string) (*entities.Purchase, error)
	GetPurchaseItemsByPurchaseID(ctx context.Context, purchaseID string) ([]*entities.PurchaseItem, error)
	GetPurchaseSenderByPurchaseID(ctx context.Context, purchaseID string) (*entities.PurchaseSender, error)
	UpdatePurchasePaymentProof(ctx context.Context, purchaseID string, paymentProofIds string) error
	GetPurchasesByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Purchase, int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreatePurchase(ctx context.Context, purchase *entities.Purchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

func (r *GormRepository) CreatePurchaseItems(ctx context.Context, items []*entities.PurchaseItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *GormRepository) CreatePurchaseSender(ctx context.Context, sender *entities.PurchaseSender) error {
	return r.db.WithContext(ctx).Create(sender).Error
}

func (r *GormRepository) GetPurchaseByID(ctx context.Context, id string) (*entities.Purchase, error) {
	var purchase entities.Purchase
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&purchase).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *GormRepository) GetPurchaseItemsByPurchaseID(ctx context.Context, purchaseID string) ([]*entities.PurchaseItem, error) {
	var items []*entities.PurchaseItem
	if err := r.db.WithContext(ctx).Where("purchase_id = ?", purchaseID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GormRepository) GetPurchaseSenderByPurchaseID(ctx context.Context, purchaseID string) (*entities.PurchaseSender, error) {
	var sender entities.PurchaseSender
	if err := r.db.WithContext(ctx).Where("purchase_id = ?", purchaseID).First(&sender).Error; err != nil {
		return nil, err
	}
	return &sender, nil
}

func (r *GormRepository) UpdatePurchasePaymentProof(ctx context.Context, purchaseID string, paymentProofIds string) error {
	return r.db.WithContext(ctx).Model(&entities.Purchase{}).
		Where("id = ?", purchaseID).
		Update("payment_proof_ids", paymentProofIds).Error
}

func (r *GormRepository) GetPurchasesByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Purchase, int64, error) {
	var purchases []*entities.Purchase
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entities.Purchase{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&purchases).Error; err != nil {
		return nil, 0, err
	}

	return purchases, total, nil
}
