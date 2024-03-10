package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FullName            *string   `gorm:"uniqueIndex"`
	UserName            string    `gorm:"not null"`
	Email               *string   `gorm:"uniqueIndex"`
	Password            string    `gorm:"not null"`
	Telephone           *string   `gorm:"uniqueIndex"`
	IsIdentityVerified  bool      `gorm:"not null"`
	IsEmailVerified     bool      `gorm:"not null"`
	IsTelephoneVerified bool      `gorm:"not null"`
	IsAddressVerified   bool      `gorm:"not null"`
	Photo               *string   // Pode ser nulo
	Biography           *string   // Pode ser nulo
	Draft               bool      `gorm:"not null"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

type SignUpInput struct {
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type SignInInput struct {
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required,min=8"`
}
type Otp struct {
	Telephone string    `json:"telephone"`
	Email     string    `json:"email"`
	CodeOtp   string    `json:"code_opt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	UserName  string    `json:"user_name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Telephone string    `json:"telephone,omitempty"`
	Biography string    `json:"biography,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
