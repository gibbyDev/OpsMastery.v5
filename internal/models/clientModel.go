package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name                 string `json:"name" gorm:"not null"`
	Email                string `json:"email" gorm:"unique;not null"`
	Title                string `json:"title"`
	PhoneNumber          string `json:"phone_number"`
	SecondaryPhoneNumber string `json:"secondary_phone_number"`
}
