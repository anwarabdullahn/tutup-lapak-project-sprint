package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"type:varchar(255)"`
	Phone     string    `gorm:"type:varchar(50)"`
	Password  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}


// Service layer request types (better practice)
type CreateUserRequest struct {
	Email    string
	Phone    string
	Password string
}

type LoginRequest struct {
	Email    string
	Phone    string
	Password string
}

// BeforeCreate ensures UUID v7 is set by the application (no DB default)
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}
