package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Name string // Optional, useful for group/general chats

	// Nullable → ticket-based chats (optional)
	TicketID *uint
	Ticket   *Ticket

	// Flag for private 1:1 chats
	IsPrivate bool `gorm:"default:false"`

	Users    []User `gorm:"many2many:chat_users;"`
	Messages []ChatMessage
}

type ChatMessage struct {
	gorm.Model
	ChatID uint
	Chat   Chat

	SentAt time.Time

	SenderID uint
	Sender   User

	Content string `gorm:"not null"`
}
