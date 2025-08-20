package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	Name         string
	Role         string `gorm:"default:User"`
	Active       bool   `gorm:"default:true"`
	Username     string `gorm:"uniqueIndex"`
	Address      string
	PhoneNumber  string
	ProfilePhoto []byte `gorm:"type:bytea"`
	ClientID     *uint
	Client       *Client

	VerificationToken string    `json:"-"`
	ResetToken        string    `json:"-"`
	ResetTokenExpiry  time.Time `json:"-"`
}
