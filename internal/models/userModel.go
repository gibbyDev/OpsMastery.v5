package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `json:"id"`
	gorm.Model
	Email             string    `json:"email" gorm:"unique"`
	Password          string    `json:"password"`
	Name              string    `json:"name"`
	Role              string    `json:"role" gorm:"default:Admin"`
	Active            bool      `json:"active" gorm:"default:true"`
	VerificationToken string    `json:"-"`
	ResetToken        string    `json:"-"`
	ResetTokenExpiry  time.Time `json:"-"`

	Username    string `json:"username,omitempty"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	// Store the image as a byte slice (BLOB in DB), omit from JSON by default
	ProfilePhoto []byte `json:"-" gorm:"type:bytea"`
}
