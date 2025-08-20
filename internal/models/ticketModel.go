package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string

	ReporterID uint
	Reporter   User `gorm:"foreignKey:ReporterID"`

	AssigneeID uint
	Assignee   User `gorm:"foreignKey:AssigneeID"`

	ClientID uint
	Client   Client

	Status   string `gorm:"default:'Open'"`   // Add this line
	Priority string `gorm:"default:'Normal'"` // Add this line

	Chats []Chat
}
