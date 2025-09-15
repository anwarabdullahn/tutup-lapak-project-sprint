package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Purchase represents a purchase order
type Purchase struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	PaymentProofIds string    `gorm:"type:text"` // JSON array of file IDs
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// PurchaseItem represents an item in a purchase
type PurchaseItem struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	PurchaseID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID    string    `gorm:"type:varchar(255);not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Category     string    `gorm:"type:varchar(255);not null"`
	Qty          int       `gorm:"not null"`
	Price        float64   `gorm:"type:decimal(10,2);not null"`
	SKU          string    `gorm:"type:varchar(255)"`
	FileID       string    `gorm:"type:varchar(255)"`
	FileURI      string    `gorm:"type:text"`
	FileThumbnailURI string `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// PurchaseSender represents sender information for a purchase
type PurchaseSender struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	PurchaseID           uuid.UUID `gorm:"type:uuid;not null"`
	SenderName           string    `gorm:"type:varchar(255);not null"`
	SenderContactType    string    `gorm:"type:varchar(50);not null"` // "email" or "phone"
	SenderContactDetail  string    `gorm:"type:varchar(255);not null"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// BeforeCreate ensures UUID v7 is set by the application
func (p *Purchase) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = id
	}
	return nil
}

// BeforeCreate ensures UUID v7 is set by the application
func (pi *PurchaseItem) BeforeCreate(tx *gorm.DB) (err error) {
	if pi.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		pi.ID = id
	}
	return nil
}

// BeforeCreate ensures UUID v7 is set by the application
func (ps *PurchaseSender) BeforeCreate(tx *gorm.DB) (err error) {
	if ps.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		ps.ID = id
	}
	return nil
}
