package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	SenderID    uint      `json:"sender_id"`
	RecipientID uint      `json:"recipient_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	// Optional: preload sender/recipient info
	Sender    User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Recipient User `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
}
