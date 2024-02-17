package models

import (
	"time"

	"github.com/google/uuid"
)

type OTPCode struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Code      string    `gorm:"not null;size:6"`
	Email     *string   `gorm:"uniqueIndex"`
	Telephone *string   `gorm:"uniqueIndex"`
	UserID    string    `gorm:"uniqueIndex"`
	Verified  bool      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type OTPConfirmationInput struct {
	Code      string `json:"code" binding:"required,len=6"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type RequestNewOTPInput struct {
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type ChangeEmailOrTelefoneBeforeToBeVerifiedInput struct {
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	UserID    string `json:"userId"`
}
