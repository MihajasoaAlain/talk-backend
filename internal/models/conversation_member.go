package models

import "time"

type ConversationMember struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"index;not null"`
	UserID         uint   `gorm:"index;not null"`
	Role           string `gorm:"not null;default:'member'"`

	CreatedAt time.Time
}
