package dtos

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
