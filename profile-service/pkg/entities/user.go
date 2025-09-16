package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// files table
type File struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	FileUri          string    `gorm:"type:varchar(255);column:fileUri"`
	FileThumbnailUri string    `gorm:"type:varchar(255);column:fileThumbnailUri"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// users table
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	Email    string    `gorm:"type:varchar(255);column:email;unique"`
	Phone    string    `gorm:"type:varchar(50);column:phone"`
	Password string    `gorm:"type:varchar(255);column:password"`

	FileId *uuid.UUID `gorm:"type:uuid;column:fileId"`
	File   *File      `gorm:"foreignKey:FileId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	BankAccountName   string `gorm:"type:varchar(32);column:bankAccountName;default:NULL"`
	BankAccountHolder string `gorm:"type:varchar(32);column:bankAccountHolder;default:NULL"`
	BankAccountNumber string `gorm:"type:varchar(32);column:bankAccountNumber;default:NULL"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;default:NULL"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;default:NULL"`
}
type AuthRequest struct {
	Email    string `gorm:"not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required,min=8,max=32"`
}

type EmailRequest struct {
	Email string `gorm:"not null" json:"email" validate:"required,email"`
}

type PhoneRequest struct {
	Phone string `gorm:"not null" json:"phone" validate:"required,min=10,max=15"`
}

type UpdateUserRequest struct {
	FileId            string `gorm:"not null" json:"fileId" validate:"required,uuid7"`
	BankAccountName   string `gorm:"not null;column:bankAccountName" json:"bankAccountName" validate:"required,min=4,max=32"`
	BankAccountHolder string `gorm:"not null;column:bankAccountHolder" json:"bankAccountHolder" validate:"required,min=4,max=32"`
	BankAccountNumber string `gorm:"not null;column:bankAccountNumber" json:"bankAccountNumber" validate:"required,min=4,max=32"`
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

var (
	ErrInvalidFileID      = errors.New("fileID is not valid")
	ErrFileNotFound       = errors.New("fileID not found")
	ErrInvalidUserID      = errors.New("userID is not valid")
	ErrInvalidPhoneNumber = errors.New("phone number is not valid")
)

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

func (u *File) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}

func (File) TableName() string { return "files" }
func (User) TableName() string { return "users" }
