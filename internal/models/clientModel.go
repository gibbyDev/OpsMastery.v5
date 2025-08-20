package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"uniqueIndex;not null"`
	Phone string

	Users   []User
	Tickets []Ticket
}
