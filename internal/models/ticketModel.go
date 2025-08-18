package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
	ClientID    uint   `json:"client_id"`
}
