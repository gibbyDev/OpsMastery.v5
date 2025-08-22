package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name string

	TicketID *uint
	Ticket   *Ticket `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsPrivate bool `gorm:"default:false"`

	Users    []User        `gorm:"many2many:chat_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages []ChatMessage `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatMessage struct {
	gorm.Model
	ChatID uint
	Chat   Chat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	SentAt time.Time

	SenderID uint
	Sender   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Content string `gorm:"not null"`
}
