package dtos

// API Request DTOs - for HTTP layer validation
type EmailRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneRegisterRequest struct {
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type EmailLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneLoginRequest struct {
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// API Response DTOs
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
