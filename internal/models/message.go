package models

import "time"

type Message struct {
	ID             uint `gorm:"primaryKey"`
	ConversationID uint `gorm:"index;not null"`
	SenderID       uint `gorm:"index;not null"`

	Content string    `gorm:"type:text;not null"`
	SentAt  time.Time `gorm:"index;not null"`

	DeliveredAt *time.Time
	ReadAt      *time.Time
}
